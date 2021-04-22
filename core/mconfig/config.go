package mconfig

import (
	"fmt"
	log "github.com/mhchlib/logger"
	"github.com/mitchellh/mapstructure"
)

// ConfigKey ...
type ConfigKey string

// ConfigVal ...
type ConfigVal string

// ConfigEntity ...
type ConfigEntity struct {
	Key ConfigKey `json:"key"`
	Val ConfigVal `json:"val"`
}

// ConfigStoreVal ...
type ConfigStoreVal struct {
	Key ConfigKey `json:"key"`
	Val ConfigVal `json:"val"`
}

// ConfigChangeNotifyMsg ...
type ConfigChangeNotifyMsg struct {
	Key ConfigKey `json:"key"`
	Val ConfigVal `json:"val"`
}

// ConfigKeys ...
func ConfigKeys(keys []string) []ConfigKey {
	configkeys := make([]ConfigKey, 0)
	for _, key := range keys {
		configkeys = append(configkeys, ConfigKey(key))
	}
	return configkeys
}

// BuildConfigStoreVal ...
func BuildConfigStoreVal(val *ConfigStoreVal) (*StoreVal, error) {
	return buildStoreVal(val)
}

// TransformMap2ConfigStoreVal ...
func TransformMap2ConfigStoreVal(val interface{}) (*ConfigStoreVal, error) {
	storeVal := &ConfigStoreVal{}
	err := mapstructure.Decode(val, &storeVal)
	if err != nil {
		log.Error("config store value transform fail:", fmt.Sprintf("%+v", val), "err:", err)
	}
	return storeVal, nil
}
