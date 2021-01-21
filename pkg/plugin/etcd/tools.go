package etcd

import (
	"encoding/json"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/pkg/config"
)

func parseAppConfigsJSONStr(value AppConfigsJSONStr) (*config.AppConfigs, error) {
	var appConfigs = make(config.AppConfigs)
	err := json.Unmarshal([]byte(value), &appConfigs)
	if err != nil {
		log.Error(Error_ParserAppConfigFail, err)
		return nil, Error_ParserAppConfigFail
	}
	return &appConfigs, nil
}
