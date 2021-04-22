package mconfig

// AppKey ...
type AppKey string

// DataVersion ...
type DataVersion struct {
	Md5     string `json:"md5"`
	Version int64  `json:"version"`
}

// ConfigEnv ...
type ConfigEnv string

// DefaultConfigEnv ...
const DefaultConfigEnv = "default"

// AppData ...
type AppData map[AppKey]map[ConfigEnv]*EnvData

// EnvData ...
type EnvData struct {
	Filter  StoreVal               `json:"filter"`
	Configs map[ConfigKey]StoreVal `json:"configs"`
}

// NodeDetail ...
type NodeDetail struct {
	Apps        *AppData `json:"apps"`
	ClientCount int32    `json:"client_count"`
}
