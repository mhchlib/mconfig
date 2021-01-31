package mconfig

type Appkey string

type ConfigKey string

type ConfigVal string

type ConfigEntity struct {
	Key ConfigKey
	Val ConfigVal
}

type FilterVal map[string]string

type ConfigEnv string

const DefaultConfigEnv = "default"

func ConfigKeys(keys []string) []ConfigKey {
	configkeys := make([]ConfigKey, 0)
	for _, key := range keys {
		configkeys = append(configkeys, ConfigKey(key))
	}
	return configkeys
}

type AppData struct {
	AppKey Appkey
	Data   map[ConfigEnv]*EnvData
}

type EnvData struct {
	Filter  FilterVal
	Configs []*ConfigEntity
}
