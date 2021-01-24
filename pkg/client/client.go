package client

import (
	"errors"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/pkg/mconfig"
	"sync/atomic"
)

var n int32 = 1000

type ClientId int32

var management *ClientConfigRelationManagement

func InitClientManagement() {
	management = NewClientConfigRelationManagement()
}

type Client struct {
	Id                          ClientId
	metadata                    MetaData
	msgBus                      *ClientMsgBus
	appKey                      mconfig.Appkey
	configKeys                  []mconfig.ConfigKey
	configEnv                   mconfig.ConfigEnv
	isbuildClientConfigRelation bool
	close                       chan interface{}
}

type MetaData map[string]string

func NewClient(metadata MetaData, send ClientSendFunc, recv ClientRecvFunc) (*Client, error) {
	if management == nil {
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
	}
	err = c.msgBus.RecvFunc(c)
	if err != nil {
		return nil, err
	}
	log.Info("remove client", c.Id, "success")
	return c, nil
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
	clientId := client.Id
	if client.isbuildClientConfigRelation {
		err := management.removeClientConfigRelation(*client)
		if err != nil {
			return err
		}
	}
	client.msgBus.Close()
	client = nil
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
