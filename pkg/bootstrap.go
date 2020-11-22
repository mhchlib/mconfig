package pkg

import (
	"context"
	"errors"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig-api/api/v1/sdk"
)

func InitMconfig(store_type, store_address string) func() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	Cancel = cancelFunc
	InitStore(store_type, store_address)
	configChan, _ := appConfigStore.WatchAppConfigsWithPrefix(ctx)
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
			log.Info("start notify change event to client ", ConfigId)
			notifyClients(ConfigId)
		case <-ctx.Done():
			log.Info("dispatchMsgToClient done ...")
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
	log.Info("send app config info to ", len(clientsChans), " clients")
}

func GetConfigFromStore(key AppId, filters *sdk.ConfigFilters) ([]*sdk.Config, error) {
	appConfigsStr, _, err := appConfigStore.GetAppConfigs(string(key))
	if appConfigsStr == "" {
		return nil, Error_AppConfigNotFound
	}
	//paser config str to ob
	//log.Info(appConfigsStr)
	appConfigs, err := parseAppConfigsJSONStr(appConfigsStr)
	//log.Info(appConfigs.AppConfigs, err)
	if err != nil {
		return nil, err
	}
	go func() {
		err = mconfigCache.putConfigCache(key, appConfigs)
		if err != nil {
			log.Error(err)
		}
	}()
	configsForClient, err := filterConfigsForClient(appConfigs, filters)
	if err != nil {
		return nil, err
	}
	return configsForClient, nil
}

func GetConfigFromCache(key AppId, filters *sdk.ConfigFilters) ([]*sdk.Config, error) {
	cache, err := mconfigCache.getConfigCache(key)
	if err != nil {
		if errors.Is(err, Error_AppConfigNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	configsForClient, err := filterConfigsForClient(cache, filters)
	if err != nil {
		return nil, err
	}
	return configsForClient, nil
}

func EndMconfig() func() {
	return func() {
		Cancel()
	}
}

func handleEventMsg(configChan chan *ConfigEvent, ctx context.Context) {
	log.Info("receive config change event started ")
	defer func() {
		log.Error("receive config change event closed ")
	}()
	for {
		select {
		case v, ok := <-configChan:
			if !ok {
				return
			}
			log.Info("get change event ", v.Key)
			//config 2 cache
			appConfigs, err := parseAppConfigsJSONStr(v.Value)
			if err != nil {
				log.Error(err)
				break
			}
			err = mconfigCache.putConfigCache(v.Key, appConfigs)
			if err != nil {
				log.Error(err)
				break
			}
			//notify client
			configChangeChan <- v.Key
			log.Info("push config change event to cache ", v.Key)
		case <-ctx.Done():
			return
		}
	}
}
