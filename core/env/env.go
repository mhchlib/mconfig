package env

import (
	"github.com/mhchlib/mconfig/core/filter"
	"github.com/mhchlib/mconfig/core/mconfig"
)

func GetEffectEnvKey(appKey mconfig.AppKey, metadata map[string]string) (mconfig.ConfigEnv, error) {
	return filter.GetEffectEnvKey(appKey, metadata)
}
