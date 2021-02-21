package mconfig

import (
	"fmt"
	log "github.com/mhchlib/logger"
	"github.com/mitchellh/mapstructure"
)

type FilterMode string

const (
	FilterMode_lua    FilterMode = "lua"
	FilterMode_simple FilterMode = "simple"
	FilterMode_mep    FilterMode = "mep"
)

type FilterEntity struct {
	Env    ConfigEnv
	Weight int
	Code   FilterVal
	Mode   FilterMode
}

type FilterStoreVal struct {
	Env    ConfigEnv  `json:"env"`
	Weight int        `json:"weight"`
	Code   FilterVal  `json:"code"`
	Mode   FilterMode `json:"mode"`
}

type FilterVal string

func BuildFilterStoreVal(val *FilterStoreVal) (*StoreVal, error) {
	return buildStoreVal(val)
}

func TransformMap2FilterStoreVal(val interface{}) (*FilterStoreVal, error) {
	storeVal := &FilterStoreVal{}
	err := mapstructure.Decode(val, &storeVal)
	if err != nil {
		log.Error("filter store value transform fail:", fmt.Sprintf("%+v", val), "err:", err)
	}
	return storeVal, nil
}
