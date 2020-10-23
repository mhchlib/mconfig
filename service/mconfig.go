package service

import (
	"context"
	"encoding/json"
	"errors"
	log "github.com/mhchlib/logger"
	sch "github.com/xeipuuv/gojsonschema"
)

type ConfigStatus int

const (
	Unpublished ConfigStatus = iota
	Published
	GrayPublished
)

type Mconfig struct {
	Status  ConfigStatus `json:"status"`
	Key     string       `json:"key"`
	Value   string       `json:"value"`
	Schema  string       `json:"schema"`
	Objects string       `json:"objects"`
}

type ConfigId string
type ConfigJSONStr string

type ConfigEvent struct {
	Key   ConfigId
	Value ConfigJSONStr
}

var (
	ConfigStore Store
	cancel      context.CancelFunc
)

func InitMconfig() func() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	cancel = cancelFunc
	configChan, _ := ConfigStore.WatchConfigWithPrefix(ctx)
	go handleEventMsg(configChan, ctx)
	go dispatchMsgToClient(ctx)
	return EndMconfig()
}

func dispatchMsgToClient(ctx context.Context) {
	for {
		select {
		case ConfigId, ok := <-configChangeChan:
			if !ok {
				return
			}
			log.Println("start notify change event to client ", ConfigId)
			notifyClients(ConfigId)
		case <-ctx.Done():
			log.Println("dispatchMsgToClient done ...")
			return
		}
	}
}

func notifyClients(id ConfigId) {
	clientsChans := clientChanMap.GetClientsChan(id)
	if clientsChans != nil {
		for _, v := range clientsChans {
			v <- &struct{}{}
		}
	}
}

func GetConfigFromStore(key ConfigId) (ConfigJSONStr, error) {
	config, _, err := ConfigStore.GetConfig(string(key))
	mconfigCache.putConfigCache(key, config)
	return config, err
}

func GetConfigFromCache(key ConfigId) (ConfigJSONStr, error) {
	return mconfigCache.getConfigCache(key)
}

func EndMconfig() func() {
	return func() {
		cancel()
	}
}

func ParseConfigJSONStr(value ConfigJSONStr) ([]Mconfig, error) {
	//parse ConfigJSONStr
	var mconfigs = []Mconfig{}
	err := json.Unmarshal([]byte(value), &mconfigs)
	if err != nil {
		return nil, err
	}
	log.Println(mconfigs)
	return mconfigs, nil
}

func CheckConfigsSchema(configs []Mconfig) error {
	for _, config := range configs {
		status := config.Status
		if status > Unpublished {
			ok, err := CheckConfigSchema(&config)
			if err != nil {
				return err
			}
			if ok == false {
				log.Println("CheckConfigsSchema failer...  ", config)
				return errors.New("CheckConfigsSchema failer ")
			}
		}
	}
	return nil
}

func CheckConfigSchema(config *Mconfig) (bool, error) {
	value := config.Value
	schema := config.Schema
	schemaLoader := sch.NewStringLoader(schema)
	documentLoader := sch.NewStringLoader(value)
	result, err := sch.Validate(schemaLoader, documentLoader)
	if err != nil {
		return false, err
	}
	return result.Valid(), nil
}

func handleEventMsg(configChan chan *ConfigEvent, ctx context.Context) {
	for {
		select {
		case v, ok := <-configChan:
			if !ok {
				return
			}
			log.Println("get change event ", v.Key)
			//config 2 cache
			err := mconfigCache.putConfigCache(v.Key, v.Value)
			if err != nil {
				log.Error(err)
				break
			}
			log.Println("start push change event to client ", v.Key)
			//notify client
			configChangeChan <- v.Key
		case <-ctx.Done():
			log.Println("handleEventMsg done ... ")
			return

		}
	}
}
