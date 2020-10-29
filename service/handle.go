package service

import (
	"context"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig-api/api/sdk"
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
	configId := ConfigId(request.Configid)
	configCache, err := GetConfigFromCache(configId)
	if err != nil {
		log.Error(err)
		return err
	}
	if configCache == "" {
		//no cache
		// pull mconfig from store
		config, err := GetConfigFromStore(configId)
		if err != nil {
			log.Error(err)
			return err
		}
		err = sendConfig(stream, config)
		if err != nil {
			return err
		}
	}
	client, err := NewClient()
	clientChanMap.AddClient(client.Id, configId, client.MsgChan)
	for {
		select {
		case <-client.MsgChan:
			log.Println("client: ", client.Id, " get msg event, configid: ", configId)
			configCache, err = GetConfigFromCache(configId)
			if err != nil {
				log.Error(err)
				return err
			}
			err := sendConfig(stream, configCache)
			if err != nil {
				return err
			}
		}
	}
}

func sendConfig(stream sdk.MConfig_GetVStreamStream, config ConfigJSONStr) error {
	err := stream.Send(&sdk.GetVResponse{
		Config: string(config),
	})
	if err != nil {
		return err
	}
	return nil
}
