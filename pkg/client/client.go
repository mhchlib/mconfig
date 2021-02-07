package client

import (
	"errors"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/pkg/mconfig"
	"sync/atomic"
)

var n int32 = 1000

var count int32 = 0

type ClientId int32

var relationManagement *ClientConfigRelationManagement

func InitClientManagement() {
	relationManagement = NewClientConfigRelationManagement()
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
	if relationManagement == nil {
		return nil, errors.New("client config relation management does not init...")
	}
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

func (client *Client) BuildClientConfigRelation(appKey mconfig.AppKey, configKeys []mconfig.ConfigKey, env mconfig.ConfigEnv) error {
	client.appKey = appKey
	client.configKeys = configKeys
	client.configEnv = env

	err := relationManagement.addClientConfigRelation(*client)
	if err != nil {
		return err
	}
	client.isbuildClientConfigRelation = true
	return nil
}

func (client *Client) RemoveClient() error {
	clientId := client.Id
	if client.isbuildClientConfigRelation {
		err := relationManagement.removeClientConfigRelation(*client)
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

func (client *Client) SendMsg(data interface{}) error {
	err := client.msgBus.sendMsg(data)
	return err
}

func (client *Client) Hold() {
	<-client.close
}
