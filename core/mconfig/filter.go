package mconfig

import (
	"fmt"
	log "github.com/mhchlib/logger"
	"github.com/mitchellh/mapstructure"
)

// FilterMode ...
type FilterMode string

const (
	// FilterMode_lua ...
	FilterMode_lua FilterMode = "lua"
	// FilterMode_simple ...
	FilterMode_simple FilterMode = "simple"
	// FilterMode_mep ...
	FilterMode_mep FilterMode = "mep"
)

// FilterEntity ...
type FilterEntity struct {
	Env    ConfigEnv
	Weight int
	Code   FilterVal
	Mode   FilterMode
}

// FilterStoreVal ...
type FilterStoreVal struct {
	Env    ConfigEnv  `json:"env"`
	Weight int        `json:"weight"`
	Code   FilterVal  `json:"code"`
	Mode   FilterMode `json:"mode"`
}

// FilterVal ...
type FilterVal string

// BuildFilterStoreVal ...
func BuildFilterStoreVal(val *FilterStoreVal) (*StoreVal, error) {
	return buildStoreVal(val)
}

// TransformMap2FilterStoreVal ...
func TransformMap2FilterStoreVal(val interface{}) (*FilterStoreVal, error) {
	storeVal := &FilterStoreVal{}
	err := mapstructure.Decode(val, &storeVal)
	if err != nil {
		log.Error("filter store value transform fail:", fmt.Sprintf("%+v", val), "err:", err)
	}
	return storeVal, nil
}
