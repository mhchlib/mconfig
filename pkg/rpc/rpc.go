package rpc

import (
	"context"
	"encoding/json"
	"errors"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig-api/api/v1/server"
	"github.com/mhchlib/mconfig/pkg/client"
	"github.com/mhchlib/mconfig/pkg/config"
	"github.com/mhchlib/mconfig/pkg/mconfig"
	"github.com/mhchlib/mconfig/pkg/store"
)

type MConfigServer struct {
}

func NewMConfigServer() *MConfigServer {
	return &MConfigServer{}
}

func (m *MConfigServer) WatchConfigStream(stream server.MConfig_WatchConfigStreamServer) error {
	request := &server.WatchConfigStreamRequest{}
	err := stream.RecvMsg(request)
	if err != nil {
		log.Error(err)
		return err
	}
	appKey := request.AppKey
	configKeys := request.ConfigKeys
	metadata := request.Metadata
	//TODO calculate
	env := "dev"
	//get data from cache or store
	configEntitys, err := config.GetConfig(mconfig.Appkey(appKey), mconfig.ConfigKeys(configKeys), mconfig.ConfigEnv(env))
	configs := make([]*server.ConfigVal, 0)
	for _, entity := range configEntitys {
		configs = append(configs, &server.ConfigVal{
			ConfigKey: string(entity.Key),
			Val:       string(entity.Val),
		})
	}
	err = stream.Send(&server.WatchConfigStreamResponse{
		Configs: configs,
	})
	if err != nil {
		return err
	}
	c, err := client.NewClient(metadata, send(stream), recv(stream))
	if err != nil {
		return err
	}
	err = config.WatchConfig(c, mconfig.Appkey(appKey), mconfig.ConfigKeys(configKeys), mconfig.ConfigEnv(env))
	if err != nil {
		return err
	}
	c.Hold()
	return nil
}

func recv(stream server.MConfig_WatchConfigStreamServer) client.ClientRecvFunc {
	return func(c *client.Client) error {
		go func() {
			for {
				data := &server.WatchConfigStreamRequest{}
				err := stream.RecvMsg(data)
				if err != nil {
					err := c.RemoveClient()
					if err != nil {
						log.Error("remove clent fail")
					}
					return
				}
			}
		}()
		return nil
	}
}

func send(stream server.MConfig_WatchConfigStreamServer) client.ClientSendFunc {
	tmp := 1
	return func(data interface{}) error {
		tmp = tmp + 1
		log.Info(tmp)

		entity, ok := data.(*mconfig.ConfigEntity)
		if !ok {
			return errors.New("translate fail")
		}
		val := &server.ConfigVal{
			ConfigKey: string(entity.Key),
			Val:       string(entity.Val),
		}
		response := &server.WatchConfigStreamResponse{
			Configs: []*server.ConfigVal{val},
		}
		log.Debug(response)
		return stream.Send(response)
	}
}

func (m *MConfigServer) GetNodeStoreData(ctx context.Context, request *server.GetNodeStoreDataRequest) (*server.GetNodeStoreDataResponse, error) {
	data, err := store.GetCurrentMConfigStore().GetSyncData()
	if err != nil {
		return nil, err
	}
	syncData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return &server.GetNodeStoreDataResponse{
		Data: syncData,
	}, nil
}
