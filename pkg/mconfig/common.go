package mconfig

type AppKey string

type ConfigKey string

type ConfigVal string

type StoreVal struct {
	Md5     string      `json:"md5"`
	Version int64       `json:"version"`
	Data    interface{} `json:"data"`
}

type ConfigEntity struct {
	Key ConfigKey `json:"key"`
	Val ConfigVal `json:"val"`
}

type ConfigStoreVal struct {
	Key ConfigKey `json:"key"`
	Val ConfigVal `json:"val"`
}

type ConfigChangeNotifyMsg struct {
	Key ConfigKey `json:"key"`
	Val ConfigVal `json:"val"`
}

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
	Weight int        `json:"weight"`
	Code   FilterVal  `json:"code"`
	Mode   FilterMode `json:"mode"`
}

type FilterVal string

type ConfigEnv string

const DefaultConfigEnv = "default"

func ConfigKeys(keys []string) []ConfigKey {
	configkeys := make([]ConfigKey, 0)
	for _, key := range keys {
		configkeys = append(configkeys, ConfigKey(key))
	}
	return configkeys
}

type AppData map[AppKey]map[ConfigEnv]*EnvData

type EnvData struct {
	Filter  FilterVal               `json:"filter"`
	Configs map[ConfigKey]ConfigVal `json:"configs"`
}

type NodeDetail struct {
	Apps        *AppData `json:"apps"`
	ClientCount int32    `json:"client_count"`
}
