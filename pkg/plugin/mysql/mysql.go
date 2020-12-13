package mysql

import (
	"context"
	"github.com/mhchlib/mconfig/pkg"
)

type MysqlStore struct {
}

func init() {
	pkg.RegisterStorePlugin("mysql", Init)
}

func Init(address string) (pkg.AppConfigStore, error) {

	return nil, nil
}

func (m MysqlStore) GetAppConfigs(key string) (pkg.AppConfigsJSONStr, int64, error) {
	panic("implement me")
}

func (m MysqlStore) PutAppConfigs(key string, value pkg.AppConfigsJSONStr) error {
	panic("implement me")
}

func (m MysqlStore) WatchAppConfigs(key string, rev int64, ctx context.Context) (chan *pkg.ConfigEvent, error) {
	panic("implement me")
}

func (m MysqlStore) WatchAppConfigsWithPrefix(ctx context.Context) (chan *pkg.ConfigEvent, error) {
	panic("implement me")
}
