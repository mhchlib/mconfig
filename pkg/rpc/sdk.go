package rpc

import (
	"errors"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig-api/api/v1/sdk"
	"github.com/mhchlib/mconfig/pkg/client"
	"github.com/mhchlib/mconfig/pkg/config"
	"github.com/mhchlib/mconfig/pkg/mconfig"
)

type MConfigSDK struct {
}

func NewMConfigSDK() *MConfigSDK {
	return &MConfigSDK{}
}

func (m *MConfigSDK) WatchConfigStream(stream sdk.MConfig_WatchConfigStreamServer) error {
	request := &sdk.WatchConfigStreamRequest{}
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
	configs := make([]*sdk.ConfigVal, 0)
	for _, entity := range configEntitys {
		configs = append(configs, &sdk.ConfigVal{
			ConfigKey: string(entity.Key),
			Val:       string(entity.Val),
		})
	}
	err = stream.Send(&sdk.WatchConfigStreamResponse{
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

func recv(stream sdk.MConfig_WatchConfigStreamServer) client.ClientRecvFunc {
	return func(c *client.Client) error {
		go func() {
			for {
				data := &sdk.WatchConfigStreamRequest{}
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

func send(stream sdk.MConfig_WatchConfigStreamServer) client.ClientSendFunc {
	tmp := 1
	return func(data interface{}) error {
		tmp = tmp + 1
		log.Info(tmp)

		entity, ok := data.(*mconfig.ConfigEntity)
		if !ok {
			return errors.New("translate fail")
		}
		val := &sdk.ConfigVal{
			ConfigKey: string(entity.Key),
			Val:       string(entity.Val),
		}
		response := &sdk.WatchConfigStreamResponse{
			Configs: []*sdk.ConfigVal{val},
		}
		log.Debug(response)
		return stream.Send(response)
	}
}
