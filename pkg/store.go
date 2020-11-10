package pkg

import "context"

type AppConfigStore interface {
	GetAppConfigs(key string) (AppConfigsJSONStr, int64, error)
	PutAppConfigs(key string, value AppConfigsJSONStr) error
	WatchAppConfigs(key string, rev int64, ctx context.Context) (chan *ConfigEvent, error)
	WatchAppConfigsWithPrefix(ctx context.Context) (chan *ConfigEvent, error)
	//...
}
