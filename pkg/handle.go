package pkg

import (
	"context"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig-api/api/v1/common"
	"github.com/mhchlib/mconfig-api/api/v1/sdk"
)

type MConfig struct {
}

func NewMConfig() *MConfig {
	return &MConfig{}
}
func (M MConfig) GetVStream(ctx context.Context, request *sdk.GetVRequest, stream sdk.MConfig_GetVStreamStream) error {
	defer func() {
		_ = stream.Close()
	}()
	appId := AppId(request.AppId)
	configCache, err := GetConfigFromCache(appId)
	if err != nil {
		log.Error(err)
		return err
	}
	if configCache == nil {
		//no cache
		// pull mconfig from store
		configCache, err = GetConfigFromStore(appId)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	err = sendConfig(stream, configCache)
	if err != nil {
		return err
	}
	client, err := NewClient()
	clientChanMap.AddClient(client.Id, appId, client.MsgChan)
	defer func() {
		clientChanMap.RemoveClient(client.Id, appId)
	}()
	clietnStreamMsg := make(chan interface{})
	go func() {
		msg := &struct{}{}
		err := stream.RecvMsg(&msg)
		if err != nil {
			log.Error("client id：", client.Id, err)
		}
		clietnStreamMsg <- msg
	}()

	for {
		select {
		case <-client.MsgChan:
			log.Info("client: ", client.Id, " get msg event, appId: ", appId)
			configCache, err = GetConfigFromCache(appId)
			if err != nil {
				log.Error(err)
				return err
			}
			err := sendConfig(stream, configCache)
			if err != nil {
				return err
			}
		case <-clietnStreamMsg:
			return nil
		}
	}
}

func sendConfig(stream sdk.MConfig_GetVStreamStream, configs []ConfigEntity) error {
	err := stream.Send(&sdk.GetVResponse{
		Configs: convConfigs(configs),
	})
	if err != nil {
		return err
	}
	return nil
}

func convConfigs(configEntitys []ConfigEntity) []*common.ConfigEntityForClient {
	configs := make([]*common.ConfigEntityForClient, len(configChangeChan))
	for _, v := range configEntitys {
		configs = append(configs,
			&common.ConfigEntityForClient{
				Schema:     v.Schema,
				Config:     v.Config,
				Status:     v.Status,
				Desc:       v.Desc,
				CreateTime: v.CreateTime,
				UpdateTime: v.UpdateTime,
			},
		)
	}
	return configs
}
