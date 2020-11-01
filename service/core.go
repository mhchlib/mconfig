package service

import (
	"context"
	"errors"
	log "github.com/mhchlib/logger"
)

func InitMconfig() func() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	Cancel = cancelFunc
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

func notifyClients(id AppId) {
	clientsChans := clientChanMap.GetClientsChan(id)
	if clientsChans != nil {
		for _, v := range clientsChans {
			v <- &struct{}{}
		}
	}
}

func GetConfigFromStore(key AppId) ([]ConfigEntity, error) {
	configStr, _, err := ConfigStore.GetConfig(string(key))
	if configStr == "" {
		return nil, errors.New("not found")
	}
	configs, err := mconfigCache.putConfigCache(key, configStr)
	return configs, err
}

func GetConfigFromCache(key AppId) ([]ConfigEntity, error) {
	return mconfigCache.getConfigCache(key)
}

func EndMconfig() func() {
	return func() {
		Cancel()
	}
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
			_, err := mconfigCache.putConfigCache(v.Key, v.Value)
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
