package mysql

import (
	"github.com/mhchlib/mconfig/pkg"
)

type MysqlStore struct {
}

func init() {
	pkg.RegisterStorePlugin("mysql", Init)
}

func Init() (pkg.AppConfigStore, error) {
	return nil, nil
}
