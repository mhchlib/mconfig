package mconfig

type AppKey string

type DataVersion struct {
	Md5     string `json:"md5"`
	Version int64  `json:"version"`
}

type ConfigEnv string

const DefaultConfigEnv = "default"

type AppData map[AppKey]map[ConfigEnv]*EnvData

type EnvData struct {
	Filter  StoreVal               `json:"filter"`
	Configs map[ConfigKey]StoreVal `json:"configs"`
}

type NodeDetail struct {
	Apps        *AppData `json:"apps"`
	ClientCount int32    `json:"client_count"`
}
