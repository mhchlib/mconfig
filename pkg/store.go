package pkg

import "context"

type Store interface {
	GetConfig(key string) (ConfigJSONStr, int64, error)
	PutConfig(key string, value ConfigJSONStr) error
	WatchConfig(key string, rev int64, ctx context.Context) (chan *ConfigEvent, error)
	WatchConfigWithPrefix(ctx context.Context) (chan *ConfigEvent, error)
	//...
}
