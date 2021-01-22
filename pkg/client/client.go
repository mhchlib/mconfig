package client

import (
	"errors"
	"github.com/mhchlib/mconfig/pkg/mconfig"
	"sync"
	"sync/atomic"
)

var n int32 = 1000

type ClientId int32

var management *ClientConfigRelationManagement

func InitClientManagement() {
	management = NewClientConfigRelationManagement()
}

type Client struct {
	sync.RWMutex
	Id                          ClientId
	metadata                    *MetaData
	msgBus                      ClientMsgBus
	appKey                      mconfig.Appkey
	configKeys                  []mconfig.ConfigKey
	configEnv                   mconfig.ConfigEnv
	isbuildClientConfigRelation bool
	willBeRemoved               bool
}

type MetaData struct {
	// loading
}

func NewClient(metadata *MetaData) (*Client, error) {
	if management == nil {
		return nil, errors.New("client config relation management does not init...")
	}
	id, err := getClientId()
	if err != nil {
		return nil, err
	}
	return &Client{
		Id:       id,
		metadata: metadata,
		msgBus:   newClientMsgBus(),
	}, nil
}

func GetOnlineClientSet(appKey mconfig.Appkey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv) *ClientSet {
	return management.getClientSet(appKey, configKey, env)
}

func getClientId() (ClientId, error) {
	id := atomic.AddInt32(&n, 1)
	return ClientId(id), nil
}

func (client *Client) BuildClientConfigRelation(appKey mconfig.Appkey, configKeys []mconfig.ConfigKey, env mconfig.ConfigEnv) error {
	client.appKey = appKey
	client.configKeys = configKeys
	client.configEnv = env

	err := management.addClientConfigRelation(*client)
	if err != nil {
		return err
	}
	client.isbuildClientConfigRelation = true
	return nil
}

func (client *Client) RemoveClient() error {
	if client.isbuildClientConfigRelation {
		err := management.removeClientConfigRelation(*client)
		if err != nil {
			return err
		}
	}
	client.Lock()
	close(client.msgBus)
	client.willBeRemoved = true
	client.Unlock()
	client = nil
	return nil
}

func (client *Client) AddMsgBus(data interface{}) error {
	client.Lock()
	defer client.Unlock()
	if client.willBeRemoved {
		return nil
	}
	client.msgBus <- data
	return nil
}
