package config

func InitConfigCenter() {
	initEvent()
}

//// Config ...
//type Config struct {
//	Schema     string `json:"schema"`
//	Config     string `json:"config"`
//	CreateTime int64  `json:"create_time"`
//	UpdateTime int64  `json:"update_time"`
//}
//
//// Configs ...
//type Configs struct {
//	Configs    ConfigsMap        `json:"configs"`
//	Desc       string            `json:"desc"`
//	CreateTime int64             `json:"create_time"`
//	UpdateTime int64             `json:"update_time"`
//	ABFilters  map[string]string `json:"ABFilters"`
//}
//
//// AppConfigsMap ...
//type AppConfigsMap struct {
//	mutex      sync.RWMutex
//	AppConfigs *AppConfigs
//}
//
//// ConfigsMap ...
//type ConfigsMap struct {
//	mutex sync.RWMutex
//	Entry map[string]*Config `json:"entry"`
//}
//
//// AppConfigs ...
//type AppConfigs map[string]*Configs

// GetConfigFromStore ...
//func GetConfigFromStore(key mconfig.Appkey, filters *sdk.ConfigFilters) ([]*sdk.Config, error) {
//	appConfigs, err := store.CurrentMConfigStore.GetAppConfigs(key)
//	//paser config str to ob
//	if err != nil {
//		return nil, err
//	}
//	go func() {
//		err = cache.mconfigCache.putConfigCache(key, appConfigs)
//		if err != nil {
//			log.Error(err)
//		}
//	}()
//	configsForClient, err := filterConfigsForClient(&AppConfigsMap{AppConfigs: appConfigs}, filters, key)
//	if err != nil {
//		return nil, err
//	}
//	return configsForClient, nil
//}

//// GetConfigFromCache ...
//func GetConfigFromCache(key mconfig.Appkey, filters *sdk.ConfigFilters) ([]*sdk.Config, error) {
//	cache, err := cache.mconfigCache.getConfigCache(key)
//	if err != nil {
//		if errors.Is(err, mconfig.Error_AppConfigNotFound) {
//			return nil, nil
//		} else {
//			return nil, err
//		}
//	}
//	configsForClient, err := filterConfigsForClient(cache, filters, key)
//	if err != nil {
//		return nil, err
//	}
//	return configsForClient, nil
//}

//func CheckConfigSchema(config *Config) (bool, error) {
//	schemaLoader := sch.NewStringLoader(config.Schema)
//	documentLoader := sch.NewStringLoader(config.Config)
//	result, err := sch.Validate(schemaLoader, documentLoader)
//	if err != nil {
//		return false, err
//	}
//	return result.Valid(), nil
//}
//
//func filterConfigsForClient(appConfigs *AppConfigsMap, filters *sdk.ConfigFilters, appkey mconfig.Appkey) ([]*sdk.Config, error) {
//	configIdLen := len(filters.ConfigIds)
//	configsForClient := make([]*sdk.Config, 0)
//	defaultChoose := common.ConfigStatus_Published
//	for i := 0; i < configIdLen; i++ {
//		needConfigId := filters.ConfigIds[i]
//		appConfigs.mutex.RLock()
//		appConfig, ok := (*appConfigs.AppConfigs)[needConfigId]
//		appConfigs.mutex.RUnlock()
//		if !ok {
//			continue
//		}
//		//match ab filter
//		abFilters := appConfig.ABFilters
//		matchABFilters := true
//		if len(abFilters) == 0 {
//			matchABFilters = false
//		} else {
//			for k, v := range abFilters {
//				//judge the extra data include abfilter map
//				if data, ok := filters.ExtraData[k]; ok {
//					if data != v {
//						matchABFilters = false
//						break
//					}
//				} else {
//					matchABFilters = false
//					break
//				}
//			}
//		}
//		if matchABFilters {
//			defaultChoose = common.ConfigStatus_ABPublished
//		}
//		appConfig.Configs.mutex.RLock()
//		config, ok := appConfig.Configs.Entry[strconv.Itoa(int(defaultChoose))]
//		appConfig.Configs.mutex.RUnlock()
//		if ok {
//			configsForClient = append(configsForClient, &sdk.Config{
//				Key:        needConfigId,
//				Schema:     config.Schema,
//				Config:     config.Config,
//				Status:     defaultChoose,
//				Desc:       appConfig.Desc,
//				CreateTime: config.CreateTime,
//				UpdateTime: config.UpdateTime,
//			})
//		} else {
//			log.Error("not found config id ", needConfigId, " status ", strconv.Itoa(int(defaultChoose)), " in app ", appkey)
//			continue
//		}
//	}
//
//	return configsForClient, nil
//}
