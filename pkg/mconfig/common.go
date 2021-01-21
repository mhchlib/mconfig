package mconfig

import (
	"github.com/golang/protobuf/ptypes/timestamp"
)

// Appkey ...
type Appkey string

// ConfigKey ...
type ConfigKey string

// ConfigVal ...
type ConfigVal string

type ConfigEnv string

type AppMetaData struct {
	key         Appkey
	description string
	tags        []string
	createTime  timestamp.Timestamp
	updateTime  timestamp.Timestamp
}

type ConfigMetaData struct {
	appKey      Appkey
	configKey   ConfigKey
	val         ConfigVal
	description string
	createTime  timestamp.Timestamp
	updateTime  timestamp.Timestamp
}
