package mconfig

import (
	"fmt"
	log "github.com/mhchlib/logger"
	"github.com/mitchellh/mapstructure"
)

type ConfigKey string

type ConfigVal string

type ConfigEntity struct {
	Key ConfigKey `json:"key"`
	Val ConfigVal `json:"val"`
}

type ConfigStoreVal struct {
	Key ConfigKey `json:"key"`
	Val ConfigVal `json:"val"`
}

type ConfigChangeNotifyMsg struct {
	Key ConfigKey `json:"key"`
	Val ConfigVal `json:"val"`
}

func ConfigKeys(keys []string) []ConfigKey {
	configkeys := make([]ConfigKey, 0)
	for _, key := range keys {
		configkeys = append(configkeys, ConfigKey(key))
	}
	return configkeys
}

func BuildConfigStoreVal(val *ConfigStoreVal) (*StoreVal, error) {
	return buildStoreVal(val)
}

func TransformMap2ConfigStoreVal(val interface{}) (*ConfigStoreVal, error) {
	storeVal := &ConfigStoreVal{}
	err := mapstructure.Decode(val, &storeVal)
	if err != nil {
		log.Error("config store value transform fail:", fmt.Sprintf("%+v", val), "err:", err)
	}
	return storeVal, nil
}
