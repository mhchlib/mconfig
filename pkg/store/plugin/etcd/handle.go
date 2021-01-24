package etcd

import (
	"errors"
	"github.com/mhchlib/mconfig/pkg/mconfig"
	"strings"
)

func parseEventKey(key string) (*KeyEntity, error) {
	//key such
	//{{namespace prefix ( custom such com.github.hchlib. )}} + {{ mode (find/watch)}} + {{content type (config/version/filter/meta)}}
	if SEPARATOR == string(key[0]) {
		key = key[1:len(key)]
	}
	splits := strings.Split(key, SEPARATOR)
	count := len(splits)
	switch count {
	case 6:
		return &KeyEntity{
			namespace: KeyNamespce(splits[0]),
			mode:      KeyMode(splits[1]),
			class:     KeyClass(splits[2]),
			appKey:    mconfig.Appkey(splits[3]),
			configKey: mconfig.ConfigKey(splits[4]),
			env:       mconfig.ConfigEnv(splits[5]),
		}, nil
	default:
		return nil, errors.New("parse event key <" + key + "> fail")
	}
}

func Prefix(prefix string, v string) string {
	return prefix + v
}
