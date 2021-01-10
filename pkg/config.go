package pkg

import (
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig-api/api/v1/common"
	"github.com/mhchlib/mconfig-api/api/v1/sdk"
	sch "github.com/xeipuuv/gojsonschema"
	"strconv"
	"sync"
)

// Config ...
type Config struct {
	Schema     string `json:"schema"`
	Config     string `json:"config"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
}

// Configs ...
type Configs struct {
	Configs    ConfigsMap        `json:"configs"`
	Desc       string            `json:"desc"`
	CreateTime int64             `json:"create_time"`
	UpdateTime int64             `json:"update_time"`
	ABFilters  map[string]string `json:"ABFilters"`
}

// AppConfigsMap ...
type AppConfigsMap struct {
	mutex      sync.RWMutex
	AppConfigs *AppConfigs
}

// ConfigsMap ...
type ConfigsMap struct {
	mutex sync.RWMutex
	Entry map[string]*Config `json:"entry"`
}

// AppConfigs ...
type AppConfigs map[string]*Configs

//
//func CheckConfigsSchema(configs []ConfigEntity) error {
//	for _, config := range configs {
//		status := config.Status
//		if status > common.ConfigStatus_Unpublished {
//			ok, err := CheckConfigSchema(config.Config, config.Schema)
//			if err != nil {
//				return err
//			}
//			if ok == false {
//				log.Info("CheckConfigsSchema failer...  ", config)
//				return errors.New("CheckConfigsSchema failer ")
//			}
//		}
//	}
//	return nil
//}

func CheckConfigSchema(config *Config) (bool, error) {
	schemaLoader := sch.NewStringLoader(config.Schema)
	documentLoader := sch.NewStringLoader(config.Config)
	result, err := sch.Validate(schemaLoader, documentLoader)
	if err != nil {
		return false, err
	}
	return result.Valid(), nil
}

func filterConfigsForClient(appConfigs *AppConfigsMap, filters *sdk.ConfigFilters, appkey Appkey) ([]*sdk.Config, error) {
	configIdLen := len(filters.ConfigIds)
	configsForClient := make([]*sdk.Config, 0)
	defaultChoose := common.ConfigStatus_Published
	for i := 0; i < configIdLen; i++ {
		needConfigId := filters.ConfigIds[i]
		appConfigs.mutex.RLock()
		appConfig, ok := (*appConfigs.AppConfigs)[needConfigId]
		appConfigs.mutex.RUnlock()
		if !ok {
			continue
		}
		//match ab filter
		abFilters := appConfig.ABFilters
		matchABFilters := true
		if len(abFilters) == 0 {
			matchABFilters = false
		} else {
			for k, v := range abFilters {
				//judge the extra data include abfilter map
				if data, ok := filters.ExtraData[k]; ok {
					if data != v {
						matchABFilters = false
						break
					}
				} else {
					matchABFilters = false
					break
				}
			}
		}
		if matchABFilters {
			defaultChoose = common.ConfigStatus_ABPublished
		}
		appConfig.Configs.mutex.RLock()
		config, ok := appConfig.Configs.Entry[strconv.Itoa(int(defaultChoose))]
		appConfig.Configs.mutex.RUnlock()
		if ok {
			configsForClient = append(configsForClient, &sdk.Config{
				Key:        needConfigId,
				Schema:     config.Schema,
				Config:     config.Config,
				Status:     defaultChoose,
				Desc:       appConfig.Desc,
				CreateTime: config.CreateTime,
				UpdateTime: config.UpdateTime,
			})
		} else {
			log.Error("not found config id ", needConfigId, " status ", strconv.Itoa(int(defaultChoose)), " in app ", appkey)
			continue
		}
	}

	return configsForClient, nil
}
