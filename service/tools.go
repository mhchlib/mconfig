package service

import (
	"encoding/json"
	"errors"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig-api/api/v1/common"
	sch "github.com/xeipuuv/gojsonschema"
)

func ParseConfigJSONStr(value ConfigJSONStr) ([]ConfigEntity, error) {
	//parse ConfigJSONStr
	var mconfigs = []ConfigEntity{}
	err := json.Unmarshal([]byte(value), &mconfigs)
	if err != nil {
		return nil, err
	}
	log.Println(mconfigs)
	return mconfigs, nil
}

func CheckConfigsSchema(configs []ConfigEntity) error {
	for _, config := range configs {
		status := config.Status
		if status > common.ConfigStatus_Unpublished {
			ok, err := CheckConfigSchema(config.Config, config.Schema)
			if err != nil {
				return err
			}
			if ok == false {
				log.Println("CheckConfigsSchema failer...  ", config)
				return errors.New("CheckConfigsSchema failer ")
			}
		}
	}
	return nil
}

func CheckConfigSchema(config string, schema string) (bool, error) {
	schemaLoader := sch.NewStringLoader(schema)
	documentLoader := sch.NewStringLoader(config)
	result, err := sch.Validate(schemaLoader, documentLoader)
	if err != nil {
		return false, err
	}
	return result.Valid(), nil
}
