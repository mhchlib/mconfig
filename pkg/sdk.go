package pkg

import (
	"crypto/md5"
	"encoding/json"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig-api/api/v1/sdk"
)

// MConfigSDK ...
type MConfigSDK struct {
}

// NewMConfigSDK ...
func NewMConfigSDK() *MConfigSDK {
	return &MConfigSDK{}
}

// GetVStream ...
func (m *MConfigSDK) GetVStream(stream sdk.MConfig_GetVStreamServer) error {
	request := &sdk.GetVRequest{}
	err := stream.RecvMsg(request)
	if err != nil {
		log.Error(err)
		return err
	}
	localConfiCacheMd5 := ""
	appKey := Appkey(request.AppKey)
	configsCache, err := GetConfigFromCache(appKey, request.Filters)
	if err != nil {
		log.Error(err)
		return err
	}
	if configsCache == nil {
		//no cache
		// pull pkg from store
		configsCache, err = GetConfigFromStore(appKey, request.Filters)
		if err != nil {
			log.Error(appKey, request.Filters, err)
			return err
		}
	}
	err = sendConfig(stream, configsCache)
	if err != nil {
		return err
	}
	client, err := NewClient()
	clientChanMap.AddClient(client.Id, appKey, client.MsgChan)
	defer func() {
		clientChanMap.RemoveClient(client.Id, appKey)
	}()
	clietnStreamMsg := make(chan interface{})
	go func() {
		msg := &struct{}{}
		err := stream.RecvMsg(&msg)
		log.Error(err)
		if err != nil {
			log.Error("client id：", client.Id, err)
		}
		clietnStreamMsg <- msg
	}()

	for {
		select {
		case <-client.MsgChan:
			log.Info("client: ", client.Id, " get msg event, appId: ", appKey)
			configsCache, err = GetConfigFromCache(appKey, request.Filters)
			if err != nil {
				log.Error(err)
				return err
			}
			if ok, md5 := checkNeedNotifyClient(localConfiCacheMd5, configsCache); ok {
				err := sendConfig(stream, configsCache)
				if err != nil {
					log.Error(err)
					return err
				}
				localConfiCacheMd5 = md5
			}
		case <-clietnStreamMsg:
			return nil
		}
	}
}

func checkNeedNotifyClient(localConfiCacheMd5 string, cache []*sdk.Config) (bool, string) {
	hash := md5.New()
	bs, _ := json.Marshal(cache)
	hash.Write(bs)
	sum := hash.Sum(nil)
	if localConfiCacheMd5 == string(sum) {
		return false, ""
	}
	return true, string(sum)
}

func sendConfig(stream sdk.MConfig_GetVStreamServer, configs []*sdk.Config) error {
	err := stream.Send(&sdk.GetVResponse{
		Configs: configs,
	})
	if err != nil {
		return err
	}
	return nil
}
