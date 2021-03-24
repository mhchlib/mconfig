package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig-api/api/v1/server"
	"github.com/mhchlib/mconfig/core/mconfig"
	"github.com/mhchlib/mconfig/core/syncx"
	"github.com/mhchlib/register"
	"google.golang.org/grpc"
	"strings"
	"time"
)

// MConfigStore ...
type MConfigStore interface {
	GetConfigVal(appKey mconfig.AppKey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv) (*mconfig.StoreVal, error)
	GetFilterVal(appKey mconfig.AppKey, env mconfig.ConfigEnv) (*mconfig.StoreVal, error)
	WatchDynamicVal(customer *Consumer) error
	PutConfigVal(appKey mconfig.AppKey, env mconfig.ConfigEnv, configKey mconfig.ConfigKey, content mconfig.StoreVal) error
	PutFilterVal(appKey mconfig.AppKey, env mconfig.ConfigEnv, content mconfig.StoreVal) error
	DeleteConfig(appKey mconfig.AppKey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv) error
	DeleteFilter(appKey mconfig.AppKey, env mconfig.ConfigEnv) error
	GetAppFilters(appKey mconfig.AppKey) ([]*mconfig.StoreVal, error)
	GetAppConfigs(appKey mconfig.AppKey, env mconfig.ConfigEnv) ([]*mconfig.StoreVal, error)
	GetAppConfigKeys(appKey mconfig.AppKey, env mconfig.ConfigEnv) ([]mconfig.ConfigKey, error)
	GetSyncData() (mconfig.AppData, error)
	PutSyncData(data *mconfig.AppData) error
	Close() error
}

var shareCalls syncx.SharedCalls

// share calls
func GetConfigVal(appKey mconfig.AppKey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv) (*mconfig.StoreVal, error) {
	key := fmt.Sprintf("%v-%v-%v-%v", "GetConfigVal", appKey, configKey, env)
	v, err := shareCalls.Do(key, func() (interface{}, error) {
		val, err := currentMConfigStore.GetConfigVal(appKey, configKey, env)
		return val, err
	})
	return v.(*mconfig.StoreVal), err
}

// GetFilterVal ...
func GetFilterVal(appKey mconfig.AppKey, env mconfig.ConfigEnv) (*mconfig.StoreVal, error) {
	key := fmt.Sprintf("%v-%v-%v", "GetFilterVal", appKey, env)
	v, err := shareCalls.Do(key, func() (interface{}, error) {
		val, err := currentMConfigStore.GetFilterVal(appKey, env)
		return val, err
	})
	return v.(*mconfig.StoreVal), err
}

// GetAppConfigKeys ...
func GetAppConfigKeys(appKey mconfig.AppKey, env mconfig.ConfigEnv) ([]mconfig.ConfigKey, error) {
	key := fmt.Sprintf("%v-%v-%v", "GetAppConfigKeys", appKey, env)
	v, err := shareCalls.Do(key, func() (interface{}, error) {
		val, err := currentMConfigStore.GetAppConfigKeys(appKey, env)
		return val, err
	})
	return v.([]mconfig.ConfigKey), err
}

// GetAppConfigs ...
func GetAppConfigs(appKey mconfig.AppKey, env mconfig.ConfigEnv) ([]*mconfig.StoreVal, error) {
	key := fmt.Sprintf("%v-%v-%v", "GetAppConfigs", appKey, env)
	v, err := shareCalls.Do(key, func() (interface{}, error) {
		val, err := currentMConfigStore.GetAppConfigs(appKey, env)
		return val, err
	})
	return v.([]*mconfig.StoreVal), err
}

// GetAppFilters ...
func GetAppFilters(appKey mconfig.AppKey) ([]*mconfig.StoreVal, error) {
	key := fmt.Sprintf("%v-%v", "GetAppFilters", appKey)
	v, err := shareCalls.Do(key, func() (interface{}, error) {
		val, err := currentMConfigStore.GetAppFilters(appKey)
		return val, err
	})
	return v.([]*mconfig.StoreVal), err
}

// GetSyncData ...
func GetSyncData() (mconfig.AppData, error) {
	key := fmt.Sprintf("%v", "GetSyncData")
	v, err := shareCalls.Do(key, func() (interface{}, error) {
		val, err := currentMConfigStore.GetSyncData()
		return val, err
	})
	return v.(mconfig.AppData), err
}

// DeleteConfig ...
func DeleteConfig(appKey mconfig.AppKey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv) error {
	key := fmt.Sprintf("%v-%v-%v-%v", "DeleteConfig", appKey, configKey, env)
	_, err := shareCalls.Do(key, func() (interface{}, error) {
		err := currentMConfigStore.DeleteConfig(appKey, configKey, env)
		return nil, err
	})
	return err
}

// DeleteFilter ...
func DeleteFilter(appKey mconfig.AppKey, env mconfig.ConfigEnv) error {
	key := fmt.Sprintf("%v-%v-%v", "DeleteFilter", appKey, env)
	_, err := shareCalls.Do(key, func() (interface{}, error) {
		//delete when no have config in this env
		configs, err := GetAppConfigs(appKey, env)
		if err != nil {
			return nil, err
		}
		if len(configs) == 0 {
			err = currentMConfigStore.DeleteFilter(appKey, env)
		} else {
			return nil, errors.New("this environment has some active configs, so cannot be deleted")
		}
		return nil, err
	})
	return err
}

// --------
func WatchDynamicVal(customer *Consumer) error {
	return currentMConfigStore.WatchDynamicVal(customer)
}

// PutConfigVal ...
func PutConfigVal(appKey mconfig.AppKey, env mconfig.ConfigEnv, configKey mconfig.ConfigKey, content mconfig.StoreVal) error {
	return currentMConfigStore.PutConfigVal(appKey, env, configKey, content)
}

// PutFilterVal ...
func PutFilterVal(appKey mconfig.AppKey, env mconfig.ConfigEnv, content mconfig.StoreVal) error {
	return currentMConfigStore.PutFilterVal(appKey, env, content)
}

// PutSyncData ...
func PutSyncData(data *mconfig.AppData) error {
	return currentMConfigStore.PutSyncData(data)
}

// Close ...
func Close() error {
	return currentMConfigStore.Close()
}

//CurrentMConfigStore
var currentMConfigStore MConfigStore
var currentStorePlugin *StorePlugin

// InitStore ...
func InitStore(storeAddressStr string) (string, error) {
	storeType, storeAddress, err := parseStoreAddressStr(storeAddressStr)
	if err != nil {
		return "", errors.New("parse store address str fail:" + err.Error())
	}
	plugin, ok := storePluginMap[storeType]
	if !ok {
		return "", errors.New(fmt.Sprintf("store type: %s does not be supported, you can choose: %s", storeType, storePluginNames))
	}
	store, err := plugin.Init(storeAddress)
	if err != nil {
		return "", err
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
	initShareCalls()
	return storeType, nil
}

func parseStoreAddressStr(str string) (string, string, error) {
	splits := strings.Split(str, "://")
	if len(splits) != 2 {
		return "", "", errors.New(str + " is invalid Address")
	}
	return splits[0], splits[1], nil
}

func initShareCalls() {
	shareCalls = syncx.NewSharedCalls()
}

// GetStorePlugin ...
func GetStorePlugin() *StorePlugin {
	return currentStorePlugin
}

// GetStorePluginModel ...
func GetStorePluginModel() StoreMode {
	return currentStorePlugin.Mode
}

// GetStorePluginName ...
func GetStorePluginName() string {
	return currentStorePlugin.Name
}

var syncRegClient register.Register
var syncServiceName string

// SyncOtherMconfigData ...
func SyncOtherMconfigData(regClient register.Register, serviceName string) error {
	syncRegClient = regClient
	syncServiceName = serviceName

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
