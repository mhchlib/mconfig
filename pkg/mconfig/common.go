package mconfig

type AppKey string

type ConfigKey string

type ConfigVal string

type ConfigEntity struct {
	Key ConfigKey `json:"key"`
	Val ConfigVal `json:"val"`
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
