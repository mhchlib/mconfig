package client

import (
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/pkg/mconfig"
	"sync"
)

type ClientConfigKey struct {
	appKey    mconfig.AppKey
	configKey mconfig.ConfigKey
	configEnv mconfig.ConfigEnv
}

type ClientConfigRelationManagement struct {
	sync.RWMutex
	m map[ClientConfigKey]*ClientSet
}

func NewClientConfigRelationManagement() *ClientConfigRelationManagement {
	management := &ClientConfigRelationManagement{}
	management.m = make(map[ClientConfigKey]*ClientSet)
	return management
}

func GetOnlineClientSet(appKey mconfig.AppKey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv) *ClientSet {
	return relationManagement.getClientSet(appKey, configKey, env)
}

func (management *ClientConfigRelationManagement) addClientConfigRelation(client Client) error {
	appKey := client.appKey
	configKeys := client.configKeys
	env := client.configEnv
	if env == "" {
		env = mconfig.DefaultConfigEnv
	}
	for _, configKey := range configKeys {
		clientConfigKey := buildClientConfigKey(appKey, configKey, env)
		management.RLock()
		set, ok := management.m[clientConfigKey]
		management.RUnlock()
		if !ok {
			set = newClientSet()
			management.Lock()
			management.m[clientConfigKey] = set
			management.Unlock()
		}
		err := set.add(client)
		if err != nil {
			return err
		}
		log.Info("add client config relation with client id: ", client.Id, " with app: ", appKey, " config key: ", configKey, " env: ", env)
	}
	return nil
}

func buildClientConfigKey(appKey mconfig.AppKey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv) ClientConfigKey {
	return ClientConfigKey{
		appKey:    appKey,
		configKey: configKey,
		configEnv: env,
	}
}

func newClientSet() *ClientSet {
	clientSet := &ClientSet{
		m: make(map[ClientId]*Client),
	}
	return clientSet
}

func (management *ClientConfigRelationManagement) removeClientConfigRelation(client Client) error {
	appKey := client.appKey
	configKeys := client.configKeys
	env := client.configEnv
	if env == "" {
		env = mconfig.DefaultConfigEnv
	}
	for _, configKey := range configKeys {
		clientConfigKey := buildClientConfigKey(appKey, configKey, env)
		management.RLock()
		set, ok := management.m[clientConfigKey]
		management.RUnlock()
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

func (management *ClientConfigRelationManagement) getClientSet(appKey mconfig.AppKey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv) *ClientSet {
	if env == "" {
		env = mconfig.DefaultConfigEnv
	}
	clientConfigKey := buildClientConfigKey(appKey, configKey, env)
	management.RLock()
	defer management.RUnlock()
	set, ok := management.m[clientConfigKey]
	if !ok {
		return nil
	}
	return set
}
