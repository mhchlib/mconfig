package etcd

import (
	"errors"
	"github.com/mhchlib/mconfig/pkg/mconfig"
	"strings"
)

func parseStoreKey(key string) (*KeyEntity, error) {
	//key such
	//{{namespace prefix ( custom such com.github.hchlib. )}} + {{content type (config/version/filter/meta)}}
	var realKey string
	if SEPARATOR == string(key[0]) {
		realKey = key[1:len(key)]
	}
	if strings.HasPrefix(key, prefix_config) {
		splits := strings.Split(realKey, SEPARATOR)
		if len(splits) == 5 {
			return &KeyEntity{
				namespace: KeyNamespce(splits[0]),
				class:     CLASS_CONFIG,
				appKey:    mconfig.AppKey(splits[2]),
				env:       mconfig.ConfigEnv(splits[3]),
				configKey: mconfig.ConfigKey(splits[4]),
			}, nil
		}

	}
	if strings.HasPrefix(key, prefix_filter) {
		splits := strings.Split(realKey, SEPARATOR)
		if len(splits) == 4 {
			return &KeyEntity{
				namespace: KeyNamespce(splits[0]),
				class:     CLASS_FILTER,
				appKey:    mconfig.AppKey(splits[2]),
				env:       mconfig.ConfigEnv(splits[3]),
			}, nil
		}
	}
	return nil, errors.New("parse event key <" + key + "> fail")
}

func getStoreKey(entity *KeyEntity) (string, error) {
	key := SEPARATOR
	if entity.namespace == "" {
		return "", errors.New("namespce can not be null")
	}
	key = key + string(entity.namespace) + SEPARATOR
	if entity.class == "" {
		return "", errors.New("class can not be null")
	}
	key = key + string(entity.class) + SEPARATOR
	if entity.appKey == "" {
		return "", errors.New("appkey can not be null")
	}
	key = key + string(entity.appKey) + SEPARATOR
	if entity.env == "" {
		return "", errors.New("envkey can not be null")
	}
	key = key + string(entity.env) + SEPARATOR

	if entity.configKey != "" {
		key = key + string(entity.configKey) + SEPARATOR
	}
	return key[0 : len(key)-len(SEPARATOR)], nil
}

func Prefix(prefix string, v string) string {
	return prefix + v
}
