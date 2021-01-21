package mysql

import (
	"github.com/mhchlib/mconfig/pkg"
	"github.com/mhchlib/mconfig/pkg/store"
)

type MysqlStore struct {
}

func init() {
	store.RegisterStorePlugin("mysql", Init)
}

func Init() (pkg.AppConfigStore, error) {
	return nil, nil
}
