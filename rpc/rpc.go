package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig-api/api/v1/server"
	"github.com/mhchlib/mconfig/core/client"
	"github.com/mhchlib/mconfig/core/config"
	"github.com/mhchlib/mconfig/core/env"
	"github.com/mhchlib/mconfig/core/mconfig"
	"github.com/mhchlib/mconfig/core/store"
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
	//get config env
	//env := "env_tPssBH6pAH0"
	configEnv, err := env.GetEffectEnvKey(mconfig.AppKey(appKey), metadata)
	if err != nil {
		return err
	}
	//get data from cache or store
	configEntitys, err := config.GetConfig(mconfig.AppKey(appKey), mconfig.ConfigKeys(configKeys), configEnv)
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
	err = c.WatchConfig(mconfig.AppKey(appKey), mconfig.ConfigKeys(configKeys), configEnv)
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
		log.Info("闭包测试", tmp)
		entity, ok := data.(*mconfig.ConfigChangeNotifyMsg)
		if !ok {
			return errors.New("translate fail")
		}
		val := &server.ConfigVal{
			ConfigKey: string(entity.Key),
			Val:       string(entity.Val),
		}
		Response := &server.WatchConfigStreamResponse{
			Configs: []*server.ConfigVal{val},
		}
		log.Debug(Response)
		return stream.Send(Response)
	}
}

func (m *MConfigServer) GetNodeStoreData(ctx context.Context, request *server.GetNodeStoreDataRequest) (*server.GetNodeStoreDataResponse, error) {
	data, err := store.GetSyncData()
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

func (m *MConfigServer) UpdateConfig(ctx context.Context, request *server.UpdateConfigRequest) (*server.UpdateConfiResponse, error) {
	filterVal := &mconfig.StoreVal{}
	err := json.Unmarshal([]byte(request.Filter), filterVal)
	if err != nil {
		log.Error(err, "data:", request.Filter)
		return nil, err
	}
	err = store.PutFilterVal(mconfig.AppKey(request.App), mconfig.ConfigEnv(request.Env), *filterVal)
	if err != nil {
		return nil, err
	}
	if request.Config != "" {
		configVal := &mconfig.StoreVal{}
		err := json.Unmarshal([]byte(request.Val), configVal)
		if err != nil {
			log.Error(err, "data:", request.Config)
			return nil, err
		}
		err = store.PutConfigVal(mconfig.AppKey(request.App), mconfig.ConfigEnv(request.Env), mconfig.ConfigKey(request.Config), *configVal)
		if err != nil {
			return nil, err
		}
	}
	return &server.UpdateConfiResponse{}, nil
}

func (m *MConfigServer) GetNodeDetail(ctx context.Context, e *empty.Empty) (*server.GetNodeDetailResponse, error) {
	storeData, err := store.GetSyncData()
	if err != nil {
		return nil, err
	}
	d, err := json.Marshal(&mconfig.NodeDetail{
		Apps:        &storeData,
		ClientCount: client.GetOnLineClientCount(),
	})
	if err != nil {
		return nil, err
	}
	return &server.GetNodeDetailResponse{
		Data: d,
	}, nil
}

func (m *MConfigServer) DeletConfig(ctx context.Context, request *server.DeletConfigRequest) (*empty.Empty, error) {
	err := store.DeleteConfig(mconfig.AppKey(request.App), mconfig.ConfigKey(request.Config), mconfig.ConfigEnv(request.Env))
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
func (m *MConfigServer) DeletFilter(ctx context.Context, request *server.DeletFilterRequest) (*empty.Empty, error) {
	err := store.DeleteFilter(mconfig.AppKey(request.App), mconfig.ConfigEnv(request.Env))
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (m *MConfigServer) UpdateFilter(ctx context.Context, request *server.UpdateFilterRequest) (*empty.Empty, error) {
	filterVal := &mconfig.StoreVal{}
	err := json.Unmarshal([]byte(request.Filter), filterVal)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	err = store.PutFilterVal(mconfig.AppKey(request.App), mconfig.ConfigEnv(request.Env), *filterVal)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
