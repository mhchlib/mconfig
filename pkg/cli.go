package pkg

import (
	"context"
	"encoding/json"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig-api/api/v1/cli"
	"github.com/mhchlib/mconfig-api/api/v1/common"
	"strconv"
	"time"
)

type MConfigCLI struct {
}

func (M *MConfigCLI) PutMconfigConfig(ctx context.Context, request *cli.PutMconfigRequest, response *cli.PutMconfigResponse) error {
	configsData, _, err := appConfigStore.GetAppConfigs(request.AppKey)
	if err != nil {
		response.Code = 500
		response.Msg = err.Error()
		return nil
	}
	appConfigs, err := parseAppConfigsJSONStr(configsData)
	if err != nil {
		response.Code = 500
		response.Msg = err.Error()
		return nil
	}
	configs, ok := appConfigs.AppConfigs[request.ConfigKey]
	if !ok {
		configs = &Configs{
			Configs: ConfigsMap{
				Entry: map[string]*Config{},
			},
			Desc:       "",
			CreateTime: time.Now().Unix(),
			UpdateTime: time.Now().Unix(),
			ABFilters:  make(map[string]string),
		}
		appConfigs.AppConfigs[request.ConfigKey] = configs
	}
	if request.Desc != "" {
		configs.Desc = request.Desc
	}
	configs.UpdateTime = time.Now().Unix()
	if request.Status == common.ConfigStatus_ABPublished {
		abFilter := make(map[string]string)
		err := json.Unmarshal([]byte(request.AbFilter), &abFilter)
		if err != nil {
			response.Code = 500
			response.Msg = "ab filter reload error"
		}
		configs.ABFilters = abFilter
	}
	configsMap := configs.Configs.Entry
	config, ok := configsMap[strconv.Itoa(int(request.Status))]
	if !ok {
		config = &Config{
			CreateTime: time.Now().Unix(),
		}
		configsMap[strconv.Itoa(int(request.Status))] = config
	}
	config.UpdateTime = time.Now().Unix()
	config.Schema = request.Schema
	config.Config = request.Config
	configsNewData, _ := json.Marshal(appConfigs.AppConfigs)
	err = appConfigStore.PutAppConfigs(request.AppKey, AppConfigsJSONStr(configsNewData))
	if err != nil {
		log.Fatal(err)
	}
	response.Msg = "success"
	response.Code = 200
	return nil
}

func (M *MConfigCLI) DeleteMconfigConfig(ctx context.Context, request *cli.DeleteMconfigConfigRequest, response *cli.DeleteMconfigConfigResponse) error {
	panic("implement me")
}

func (M *MConfigCLI) InitMconfigApp(ctx context.Context, request *cli.InitMconfigAppRequest, response *cli.InitMconfigAppResponse) error {
	configsData, _, _ := appConfigStore.GetAppConfigs(request.AppKey)
	if configsData != "" {
		response.Code = 500
		response.Msg = "This app already exists"
	}
	err := appConfigStore.PutAppConfigs(request.AppKey, "{}")
	if err != nil {
		return err
	}
	response.Code = 200
	response.Msg = "Init app " + request.AppKey + " success"
	return nil
}

func (M *MConfigCLI) UpdateMconfigApp(ctx context.Context, request *cli.UpdateMconfigAppRequest, response *cli.UpdateMconfigAppResponse) error {
	panic("implement me")
}

func (M *MConfigCLI) DeleteMconfigApp(ctx context.Context, request *cli.DeleteMconfigAppRequest, response *cli.DeleteMconfigAppResponse) error {
	panic("implement me")
}

func NewMConfigCLI() *MConfigCLI {
	return &MConfigCLI{}
}