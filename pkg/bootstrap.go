package pkg

import (
	"context"
	"errors"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig-api/api/v1/sdk"
)

func InitMconfig(mconfig *MConfig) func() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	Cancel = cancelFunc
	InitStore(*mconfig.StoreType, *mconfig.StoreAddress)
	configChan, _ := appConfigStore.WatchAppConfigsWithPrefix(ctx)
	go handleEventMsg(configChan, ctx)
	go dispatchMsgToClient(ctx)
	return EndMconfig()
}

func dispatchMsgToClient(ctx context.Context) {
	for {
		select {
		case AppId, ok := <-configChangeChan:
			if !ok {
				return
			}
			log.Info("app: ", AppId, "is changed, notify event to clients")
			notifyClients(AppId)
		case <-ctx.Done():
			log.Info("the function dispatch msg to client is done")
			return
		}
	}
}

func notifyClients(id Appkey) {
	clientsChans := clientChanMap.GetClientsChan(id)
	if clientsChans != nil {
		for _, v := range clientsChans {
			v <- &struct{}{}
		}
	}
	log.Info("notify app config change info to ", len(clientsChans), " clients")
}

func GetConfigFromStore(key Appkey, filters *sdk.ConfigFilters) ([]*sdk.Config, error) {
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
	configsForClient, err := filterConfigsForClient(appConfigs, filters, key)
	if err != nil {
		return nil, err
	}
	return configsForClient, nil
}

func GetConfigFromCache(key Appkey, filters *sdk.ConfigFilters) ([]*sdk.Config, error) {
	cache, err := mconfigCache.getConfigCache(key)
	if err != nil {
		if errors.Is(err, Error_AppConfigNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	configsForClient, err := filterConfigsForClient(cache, filters, key)
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
	log.Info("receive app config change event is started ")
	defer func() {
		log.Error("receive app config change event is closed ")
	}()
	for {
		select {
		case v, ok := <-configChan:
			if !ok {
				return
			}
			log.Info("receive app ", v.Key, " config change event ")
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
			log.Info("push app ", v.Key, " config change event to cache ")
		case <-ctx.Done():
			return
		}
	}
}
