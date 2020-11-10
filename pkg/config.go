package pkg

import (
	"encoding/json"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig-api/api/v1/common"
	"github.com/mhchlib/mconfig-api/api/v1/sdk"
	sch "github.com/xeipuuv/gojsonschema"
	"strconv"
)

func parseAppConfigsJSONStr(value AppConfigsJSONStr) (AppConfigs, error) {
	//parse AppConfigsJSONStr
	var appConfigs = make(AppConfigs)
	err := json.Unmarshal([]byte(value), &appConfigs)
	if err != nil {
		log.Error(Error_ParserAppConfigFail, err)
		return nil, Error_ParserAppConfigFail
	}
	return appConfigs, nil
}

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

func CheckConfigSchema(config string, schema string) (bool, error) {
	schemaLoader := sch.NewStringLoader(schema)
	documentLoader := sch.NewStringLoader(config)
	result, err := sch.Validate(schemaLoader, documentLoader)
	if err != nil {
		return false, err
	}
	return result.Valid(), nil
}

func filterConfigsForClient(appConfigs AppConfigs, filters *sdk.ConfigFilters) ([]*sdk.Config, error) {
	configIdLen := len(filters.ConfigIds)
	configsForClient := make([]*sdk.Config, 0)
	defaultChoose := common.ConfigStatus_Published
	for i := 0; i < configIdLen; i++ {
		needConfigId := filters.ConfigIds[i]
		appConfig, ok := appConfigs[needConfigId]
		log.Info(appConfigs, needConfigId)
		if !ok {
			continue
		}
		//match ab filter
		abFilters := appConfig.ABFilters
		matchABFilters := true
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
		if matchABFilters {
			defaultChoose = common.ConfigStatus_ABPublished
		}
		if config, ok := appConfig.Configs[strconv.Itoa(int(defaultChoose))]; ok {
			configsForClient = append(configsForClient, &sdk.Config{
				Id:         needConfigId,
				Schema:     config.Config,
				Config:     config.Schema,
				Status:     defaultChoose,
				Desc:       appConfig.Desc,
				CreateTime: config.CreateTime,
				UpdateTime: config.UpdateTime,
			})
		} else {
			log.Error("no found config id ", needConfigId, " status ", strconv.Itoa(int(defaultChoose)), " in app ")
			continue
		}
	}

	return configsForClient, nil
}
