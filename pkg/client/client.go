package client

import (
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/pkg/config"
	"github.com/mhchlib/mconfig/pkg/filter"
	"github.com/mhchlib/mconfig/pkg/mconfig"
	"sync/atomic"
)

var n int32 = 1000

var count int32 = 0

type ClientId int32

func InitClientManagement() {
	initRelationMap()
	initEvent()
}

type Client struct {
	Id                          ClientId
	metadata                    MetaData
	msgBus                      *ClientMsgBus
	appKey                      mconfig.AppKey
	configKeys                  []mconfig.ConfigKey
	configEnv                   mconfig.ConfigEnv
	isbuildClientConfigRelation bool
	close                       chan interface{}
}

type MetaData map[string]string

func NewClient(metadata MetaData, send ClientSendFunc, recv ClientRecvFunc) (*Client, error) {
	id, err := getClientId()
	if err != nil {
		return nil, err
	}
	c := &Client{
		Id:       id,
		metadata: metadata,
		msgBus:   newClientMsgBus(send, recv),
		close:    make(chan interface{}),
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

func GetOnLineClientCount() int32 {
	return count
}

func getClientId() (ClientId, error) {
	id := atomic.AddInt32(&n, 1)
	return ClientId(id), nil
}

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

func (client *Client) SendConfigChangeNotifyMsg(data *mconfig.ConfigChangeNotifyMsg) error {
	err := client.msgBus.sendMsg(data)
	return err
}

func (client *Client) Hold() {
	<-client.close
}

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

func (client *Client) WatchConfig(appKey mconfig.AppKey, configKeys []mconfig.ConfigKey, env mconfig.ConfigEnv) error {
	err := client.BuildClientRelation(appKey, configKeys, env)
	if err != nil {
		return err
	}
	return nil
}
