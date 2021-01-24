package rpc

//
//// MConfigCLI ...
//type MConfigCLI struct {
//}
//
//// NewMConfigCLI ...
//func NewMConfigCLI() *MConfigCLI {
//	return &MConfigCLI{}
//}
//
//// PutMconfigConfig ...
//func (M *MConfigCLI) PutMconfigConfig(ctx context.Context, request *cli.PutMconfigRequest) (*cli.PutMconfigResponse, error) {
//	appConfigs, err := pkg.ConfigStore.GetAppConfigs(mconfig-server.Appkey(request.AppKey))
//	response := &cli.PutMconfigResponse{}
//	if err != nil {
//		response.Code = 500
//		response.Msg = err.Error()
//		return response, nil
//	}
//	configs, ok := (*appConfigs)[request.ConfigKey]
//	if !ok {
//		configs = &config.Configs{
//			Configs: config.ConfigsMap{
//				Entry: map[string]*config.Config{},
//			},
//			Desc:       "",
//			CreateTime: time.Now().Unix(),
//			UpdateTime: time.Now().Unix(),
//			ABFilters:  make(map[string]string),
//		}
//		(*appConfigs)[request.ConfigKey] = configs
//	}
//	if request.Desc != "" {
//		configs.Desc = request.Desc
//	}
//	configs.UpdateTime = time.Now().Unix()
//	if request.Status == common.ConfigStatus_ABPublished {
//		abFilter := make(map[string]string)
//		err := json.Unmarshal([]byte(request.AbFilter), &abFilter)
//		if err != nil {
//			response.Code = 500
//			response.Msg = "ab filter reload error"
//		}
//		configs.ABFilters = abFilter
//	}
//	configsMap := configs.Configs.Entry
//	config, ok := configsMap[strconv.Itoa(int(request.Status))]
//	if !ok {
//		config = &config.Config{
//			CreateTime: time.Now().Unix(),
//		}
//		configsMap[strconv.Itoa(int(request.Status))] = config
//	}
//	config.UpdateTime = time.Now().Unix()
//	config.Schema = request.Schema
//	config.Config = request.Config
//	err = pkg.ConfigStore.PutAppConfigs(mconfig-server.Appkey(request.AppKey), appConfigs)
//	if err != nil {
//		log.Fatal(err)
//	}
//	response.Msg = "success"
//	response.Code = 200
//	return response, nil
//}
//
//// DeleteMconfigConfig ...
//func (M *MConfigCLI) DeleteMconfigConfig(ctx context.Context, request *cli.DeleteMconfigConfigRequest) (*cli.DeleteMconfigConfigResponse, error) {
//	panic("implement me")
//}
//
//// InitMconfigApp ...
//func (M *MConfigCLI) InitMconfigApp(ctx context.Context, request *cli.InitMconfigAppRequest) (*cli.InitMconfigAppResponse, error) {
//	response := &cli.InitMconfigAppResponse{}
//	appConfigs, _ := pkg.ConfigStore.GetAppConfigs(mconfig-server.Appkey(request.AppKey))
//	if appConfigs != nil {
//		response.Code = 500
//		response.Msg = "the app already exists"
//	}
//	err := pkg.ConfigStore.PutAppConfigs(mconfig-server.Appkey(request.AppKey), &config.AppConfigs{})
//	if err != nil {
//		return response, err
//	}
//	response.Code = 200
//	response.Msg = "Init app " + request.AppKey + " success"
//	return response, nil
//}
//
//// UpdateMconfigApp ...
//func (M *MConfigCLI) UpdateMconfigApp(ctx context.Context, request *cli.UpdateMconfigAppRequest) (*cli.UpdateMconfigAppResponse, error) {
//	panic("implement me")
//}
//
//// DeleteMconfigApp ...
//func (M *MConfigCLI) DeleteMconfigApp(ctx context.Context, request *cli.DeleteMconfigAppRequest) (*cli.DeleteMconfigAppResponse, error) {
//	panic("implement me")
//}
