package env

import (
	"github.com/mhchlib/mconfig/pkg/filter"
	"github.com/mhchlib/mconfig/pkg/mconfig"
)

func GetEffectEnvKey(appKey mconfig.AppKey, metadata map[string]string) (mconfig.ConfigEnv, error) {
	return filter.GetEffectEnvKey(appKey, metadata)
}
