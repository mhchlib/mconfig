package client

import (
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/pkg/mconfig"
	"sync"
)

type RelationDetail struct {
	appKey     mconfig.Appkey
	configKeys []mconfig.ConfigKey
}

type RelationMap struct {
	sync.RWMutex
	m map[ClientId]*RelationDetail
}

type ClientConfigRelationManagement struct {
	sync.RWMutex
	m map[mconfig.Appkey]*ClientConfigRelations
}

type ClientConfigRelations struct {
	sync.RWMutex
	m map[mconfig.ConfigKey]*ClientSet
}

var (
	dict RelationMap
)

func NewClientConfigRelations() *ClientConfigRelations {
	return &ClientConfigRelations{
		m: make(map[mconfig.ConfigKey]*ClientSet),
	}
}

func NewClientConfigRelationManagement() *ClientConfigRelationManagement {
	management := &ClientConfigRelationManagement{}
	management.m = make(map[mconfig.Appkey]*ClientConfigRelations)
	return management
}

func (management *ClientConfigRelationManagement) addClientConfigRelation(client Client, appKey mconfig.Appkey, configKeys []mconfig.ConfigKey) error {
	management.RLock()
	clientConfigRelations, ok := management.m[appKey]
	management.RUnlock()
	if ok == false {
		clientConfigRelations = NewClientConfigRelations()
		management.Lock()
		management.m[appKey] = clientConfigRelations
		management.Unlock()
	}
	for _, configKey := range configKeys {
		clientConfigRelations.RLock()
		clientSet, ok := clientConfigRelations.m[configKey]
		clientConfigRelations.RUnlock()
		if ok == false {
			clientSet = newClientSet()
			clientConfigRelations.Lock()
			clientConfigRelations.m[configKey] = clientSet
			clientConfigRelations.Unlock()
		}
		clientSet.add(client)
	}
	dict.Lock()
	if dict.m == nil {
		dict.m = make(map[ClientId]*RelationDetail)
	}
	dict.m[client.Id] = &RelationDetail{
		appKey:     appKey,
		configKeys: configKeys,
	}
	dict.Unlock()
	log.Info("add client config relation with client id: ", client.Id, " with app: ", appKey, " config keys: ", configKeys)
	return nil
}

func newClientSet() *ClientSet {
	clientSet := &ClientSet{
		m: make(map[ClientId]*Client),
	}
	return clientSet
}

func (management *ClientConfigRelationManagement) removeClientConfigRelation(client Client) error {
	dict.RLock()
	detail, ok := dict.m[client.Id]
	dict.RUnlock()
	if !ok {
		return nil
	}
	appKey := detail.appKey
	configKeys := detail.configKeys
	management.RLock()
	clientConfigRelations, ok := management.m[appKey]
	management.RUnlock()
	if ok == false {
		return nil
	}
	clientConfigRelations.RLock()
	for _, configKey := range configKeys {
		clientSet, ok := clientConfigRelations.m[configKey]
		if !ok {
			continue
		}
		clientSet.remove(client)
	}
	clientConfigRelations.RUnlock()
	dict.Lock()
	delete(dict.m, client.Id)
	dict.Unlock()
	return nil
}

func (management *ClientConfigRelationManagement) getClientSet(appKey mconfig.Appkey, configKey mconfig.ConfigKey) *ClientSet {
	management.RLock()
	clientConfigRelations, ok := management.m[appKey]
	management.RUnlock()
	if ok == false {
		return nil
	}
	clientConfigRelations.RLock()
	clientSet, ok := clientConfigRelations.m[configKey]
	clientConfigRelations.RUnlock()
	if !ok {
		return nil
	}
	return clientSet
}
