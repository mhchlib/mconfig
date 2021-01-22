package mconfig

import (
	"github.com/golang/protobuf/ptypes/timestamp"
)

type Appkey string

type ConfigKey string

type ConfigVal string

type FilterVal string

type ConfigEnv string

const DefaultConfigEnv = "default"

type AppMetaData struct {
	key         Appkey
	name        string
	description string
	tags        []string
	createTime  timestamp.Timestamp
	updateTime  timestamp.Timestamp
}

type ConfigMetaData struct {
	appKey      Appkey
	configKey   ConfigKey
	val         ConfigVal
	name        string
	description string
	createTime  timestamp.Timestamp
	updateTime  timestamp.Timestamp
}
