package etcd

import (
	"encoding/json"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/pkg"
)

func parseAppConfigsJSONStr(value AppConfigsJSONStr) (*pkg.AppConfigs, error) {
	var appConfigs = make(pkg.AppConfigs)
	err := json.Unmarshal([]byte(value), &appConfigs)
	if err != nil {
		log.Error(Error_ParserAppConfigFail, err)
		return nil, Error_ParserAppConfigFail
	}
	return &appConfigs, nil
}
