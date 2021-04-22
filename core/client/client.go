package client

import (
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/core/cache"
	"github.com/mhchlib/mconfig/core/config"
	"github.com/mhchlib/mconfig/core/filter"
	"github.com/mhchlib/mconfig/core/mconfig"
	"sync/atomic"
)

var n int32 = 1000

var count int32 = 0

// ClientId ...
type ClientId int32

// InitClientManagement ...
func InitClientManagement() {
	initRelationMap()
	initEvent()
}

// Client ...
type Client struct {
	Id                          ClientId
	metadata                    MetaData
	msgBus                      *ClientMsgBus
	appKey                      mconfig.AppKey
	configKeys                  []mconfig.ConfigKey
	configEnv                   mconfig.ConfigEnv
	isbuildClientConfigRelation bool
	close                       chan interface{}
	configUpdateMsgCache        cache.Cache
}

// MetaData ...
type MetaData map[string]string

// NewClient ...
func NewClient(metadata MetaData, send ClientSendFunc, recv ClientRecvFunc) (*Client, error) {
	id, err := getClientId()
	if err != nil {
		return nil, err
	}
	c := &Client{
		Id:                   id,
		metadata:             metadata,
		msgBus:               newClientMsgBus(send, recv),
		close:                make(chan interface{}),
		configUpdateMsgCache: cache.NewCache(),
	}
	err = c.msgBus.RecvFunc(c)
	if err != nil {
		return nil, err
	}
	increaseClientCount()
	log.Info("add client", c.Id, "success")
	return c, nil
}

func increaseClientCount() {
	atomic.AddInt32(&count, 1)
}

func reduceClientCount() {
	atomic.AddInt32(&count, -1)
}

// GetOnLineClientCount ...
func GetOnLineClientCount() int32 {
	return count
}

func getClientId() (ClientId, error) {
	id := atomic.AddInt32(&n, 1)
	return ClientId(id), nil
}

// BuildClientRelation ...
func (client *Client) BuildClientRelation(appKey mconfig.AppKey, configKeys []mconfig.ConfigKey, env mconfig.ConfigEnv) error {
	client.appKey = appKey
	client.configKeys = configKeys
	client.configEnv = env
	err := buildClientRelation(client)
	if err != nil {
		return err
	}
	client.isbuildClientConfigRelation = true
	return nil
}

// RemoveClient ...
func (client *Client) RemoveClient() error {
	clientId := client.Id
	if client.isbuildClientConfigRelation {
		err := removeClientRelation(client)
		if err != nil {
			return err
		}
	}
	client.msgBus.Close()
	client.close <- struct{}{}
	client = nil
	reduceClientCount()
	log.Info("remove client", clientId, "success")
	return nil
}

// SendConfigChangeNotifyMsg ...
func (client *Client) SendConfigChangeNotifyMsg(data *mconfig.ConfigChangeNotifyMsg) error {
	//check cache exist
	exist := client.checkConfigUpdateMsgCacheExist(data.Key, data)
	if exist {
		return nil
	}
	err := client.msgBus.sendMsg(data)
	if err == nil {
		//put cache
		err = client.putConfigUpdateMsgCache(data.Key, data)
		if err != nil {
			log.Error("client msg bus put cache error")
		}
	}
	return err
}

// Hold ...
func (client *Client) Hold() {
	<-client.close
}

// ReCalEffectEnv ...
func (client *Client) ReCalEffectEnv() error {
	env, err := filter.GetEffectEnvKey(client.appKey, client.metadata)
	if err != nil {
		return err
	}
	client.configEnv = env
	err = client.ReloadNewConfig()
	if err != nil {
		return err
	}
	return nil
}

// ReloadNewConfig ...
func (client *Client) ReloadNewConfig() error {
	configs, err := config.GetConfig(client.appKey, client.configKeys, client.configEnv)
	if err != nil {
		log.Error(err)
		return err
	}
	for _, c := range configs {
		err = client.SendConfigChangeNotifyMsg(&mconfig.ConfigChangeNotifyMsg{
			Key: c.Key,
			Val: c.Val,
		})
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

// WatchConfig ...
func (client *Client) WatchConfig(appKey mconfig.AppKey, configKeys []mconfig.ConfigKey, env mconfig.ConfigEnv) error {
	err := client.BuildClientRelation(appKey, configKeys, env)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) checkConfigUpdateMsgCacheExist(key interface{}, data interface{}) bool {
	md5 := mconfig.GetInterfaceMd5(data)
	value, err := client.configUpdateMsgCache.GetCache(key)
	if err != nil {
		return false
	}
	if value.(string) == md5 {
		return true
	}
	return false
}

func (client *Client) putConfigUpdateMsgCache(key interface{}, data interface{}) error {
	return client.configUpdateMsgCache.PutCache(key, mconfig.GetInterfaceMd5(data))
}
