package client

import (
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/core/config"
	"github.com/mhchlib/mconfig/core/filter"
	"github.com/mhchlib/mconfig/core/mconfig"
	"sync"
)

var configRelationMap *ClientConfigRelationMap
var appRelationMap *ClientAppRelationMap

type ClientConfigKey struct {
	appKey    mconfig.AppKey
	configKey mconfig.ConfigKey
	configEnv mconfig.ConfigEnv
}

type ClientConfigRelationMap struct {
	sync.RWMutex
	m map[ClientConfigKey]*ClientSet
}

type ClientAppRelationMap struct {
	sync.RWMutex
	m map[mconfig.AppKey]*ClientSet
}

func initRelationMap() {
	configRelationMap = &ClientConfigRelationMap{}
	configRelationMap.m = make(map[ClientConfigKey]*ClientSet)
	appRelationMap = &ClientAppRelationMap{}
	appRelationMap.m = make(map[mconfig.AppKey]*ClientSet)
	return
}

func buildClientRelation(client *Client) error {
	err := configRelationMap.addClientConfigRelation(client)
	if err != nil {
		return err
	}
	err = appRelationMap.addClientConfigRelation(client)
	if err != nil {
		return err
	}
	//register config notify
	err = config.RegisterAppNotify(client.appKey)
	if err != nil {
		return err
	}
	return nil
}

func removeClientRelation(client *Client) error {
	err := configRelationMap.removeClientConfigRelation(client)
	if err != nil {
		return err
	}
	err = appRelationMap.removeClientConfigRelation(client)
	if err != nil {
		return err
	}
	//unregister app notify
	count := getOnlineClientSetCountByAppRealtion(client.appKey)
	if count == 0 {
		err := config.UnRegisterAppNotify(client.appKey)
		if err != nil {
			log.Error(err)
		}
		//删除缓存数据
		err = config.DeleteConfigFromCacheByApp(client.appKey)
		if err != nil {
			log.Error(err)
		}
		err = filter.DeleteFilterFromCacheByApp(client.appKey)
		if err != nil {
			log.Error(err)
		}
	}
	return nil
}

func getOnlineClientSetByConfigRealtion(appKey mconfig.AppKey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv) *ClientSet {
	return configRelationMap.getClientSet(appKey, configKey, env)
}

func getOnlineClientSetByAppRealtion(appKey mconfig.AppKey) *ClientSet {
	return appRelationMap.getClientSet(appKey)
}

func getOnlineClientSetCountByAppRealtion(appKey mconfig.AppKey) int {
	return appRelationMap.getClientSetCount(appKey)
}

func (m *ClientAppRelationMap) addClientConfigRelation(client *Client) error {
	appKey := client.appKey
	m.RLock()
	s, ok := m.m[appKey]
	m.RUnlock()
	if !ok {
		s = NewClientSet()
		m.Lock()
		m.m[appKey] = s
		m.Unlock()
	}
	err := s.add(client)
	if err != nil {
		return err
	}
	log.Info("add client config relation with client id: ", client.Id, " with app: ", appKey)
	return nil
}

func (m *ClientAppRelationMap) removeClientConfigRelation(client *Client) error {
	appKey := client.appKey
	m.RLock()
	set, ok := m.m[appKey]
	m.RUnlock()
	if !ok {
		return nil
	}
	err := set.remove(client)
	if err != nil {
		return err
	}
	log.Info("remove client app relation with client id: ", client.Id, " with app: ", appKey)
	return nil
}

func (m *ClientAppRelationMap) getClientSet(appKey mconfig.AppKey) *ClientSet {
	m.RLock()
	defer m.RUnlock()
	s, ok := m.m[appKey]
	if !ok {
		return nil
	}
	return s
}

func (m *ClientAppRelationMap) getClientSetCount(appKey mconfig.AppKey) int {
	m.RLock()
	defer m.RUnlock()
	set := m.m[appKey]
	if set == nil {
		return 0
	}
	return set.count()
}

func (m *ClientConfigRelationMap) addClientConfigRelation(client *Client) error {
	appKey := client.appKey
	configKeys := client.configKeys
	env := client.configEnv
	if env == "" {
		env = mconfig.DefaultConfigEnv
	}
	for _, configKey := range configKeys {
		clientConfigKey := buildClientConfigKey(appKey, configKey, env)
		m.RLock()
		s, ok := m.m[clientConfigKey]
		m.RUnlock()
		if !ok {
			s = NewClientSet()
			m.Lock()
			m.m[clientConfigKey] = s
			m.Unlock()
		}
		err := s.add(client)
		if err != nil {
			return err
		}
		log.Info("add client config relation with client id: ", client.Id, " with app: ", appKey, " config key: ", configKey, " env: ", env)
	}
	return nil
}

func (m *ClientConfigRelationMap) removeClientConfigRelation(client *Client) error {
	appKey := client.appKey
	configKeys := client.configKeys
	env := client.configEnv
	if env == "" {
		env = mconfig.DefaultConfigEnv
	}
	for _, configKey := range configKeys {
		clientConfigKey := buildClientConfigKey(appKey, configKey, env)
		m.RLock()
		set, ok := m.m[clientConfigKey]
		m.RUnlock()
		if !ok {
			return nil
		}
		err := set.remove(client)
		if err != nil {
			return err
		}
		log.Info("remove client config relation with client id: ", client.Id, " with app: ", appKey, " config key: ", configKey, " env: ", env)
	}
	return nil
}

func (m *ClientConfigRelationMap) getClientSet(appKey mconfig.AppKey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv) *ClientSet {
	if env == "" {
		env = mconfig.DefaultConfigEnv
	}
	clientConfigKey := buildClientConfigKey(appKey, configKey, env)
	m.RLock()
	defer m.RUnlock()
	s, ok := m.m[clientConfigKey]
	if !ok {
		return nil
	}
	return s
}

func buildClientConfigKey(appKey mconfig.AppKey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv) ClientConfigKey {
	return ClientConfigKey{
		appKey:    appKey,
		configKey: configKey,
		configEnv: env,
	}
}
