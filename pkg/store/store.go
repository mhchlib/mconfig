package store

import (
	"context"
	"encoding/json"
	"errors"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig-api/api/v1/server"
	"github.com/mhchlib/mconfig/pkg/mconfig"
	"github.com/mhchlib/register/reg"
	"google.golang.org/grpc"
	"time"
)

// MConfigStore ...
type MConfigStore interface {
	GetConfigVal(appKey mconfig.AppKey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv) (mconfig.ConfigVal, error)
	WatchDynamicVal(customer *Consumer) error

	PutConfigVal(appKey mconfig.AppKey, env mconfig.ConfigEnv, configKey mconfig.ConfigKey, content mconfig.ConfigVal) error
	PutFilterVal(appKey mconfig.AppKey, env mconfig.ConfigEnv, content mconfig.FilterVal) error

	DeleteConfig(appKey mconfig.AppKey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv) error
	DeleteFilter(appKey mconfig.AppKey, env mconfig.ConfigEnv) error

	GetSyncData() (mconfig.AppData, error)
	PutSyncData(data *mconfig.AppData) error
	Close() error
}

//CurrentMConfigStore
var currentMConfigStore MConfigStore
var currentStorePlugin *StorePlugin

// InitStore ...
func InitStore(storeType string, storeAddress string) {
	plugin, ok := storePluginMap[storeType]
	if !ok {
		log.Fatal("store type:", storeType, "does not be supported, you can choose:", storePluginNames)
	}
	store, err := plugin.Init(storeAddress)
	if err != nil {
		log.Fatal(err)
	}
	currentMConfigStore = store
	log.Info("store init success with", storeType, storeAddress)
	go func() {
		err = currentMConfigStore.WatchDynamicVal(newConsumer())
		if err != nil {
			log.Error(err)
		}
	}()
	currentStorePlugin = plugin
}

func GetStorePlugin() *StorePlugin {
	return currentStorePlugin
}

func GetCurrentMConfigStore() MConfigStore {
	return currentMConfigStore
}

func CheckSyncData() bool {
	if currentStorePlugin.Mode == MODE_LOCAL {
		return true
	}
	return false
}

func SyncOtherMconfigData(regClient reg.Register, serviceName string) error {
	allServices, err := regClient.ListAllServices(serviceName)
	if err != nil {
		return err
	}
	for _, service := range allServices {
		metadata := service.Metadata
		mode := StoreMode(metadata["mode"].(string))
		if MODE_SHARE == mode {
			withTimeout, _ := context.WithTimeout(context.Background(), time.Second*5)
			dial, err := grpc.DialContext(withTimeout, service.Address, grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Info(err, " addr: ", service)
				continue
			}
			mconfigService := server.NewMConfigClient(dial)
			withTimeout, _ = context.WithTimeout(context.Background(), time.Second*20)
			syncResponse, err := mconfigService.GetNodeStoreData(withTimeout, &server.GetNodeStoreDataRequest{})
			if err != nil {
				log.Error(err)
				return err
			}
			//sync data to store
			syncData := &mconfig.AppData{}
			err = json.Unmarshal(syncResponse.Data, &syncData)
			if err != nil {
				log.Error(err)
				return err
			}
			log.Info("sync node data:", string(syncResponse.Data))
			err = currentMConfigStore.PutSyncData(syncData)
			if err != nil {
				log.Error(err)
				return err
			}
			return nil
		}
	}
	return errors.New("not found sync node")
}
