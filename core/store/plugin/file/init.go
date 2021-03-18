package file

import (
	"github.com/mhchlib/mconfig/core/store"
	"github.com/micro/go-micro/v2/util/file"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
)

func init() {
	store.RegisterStorePlugin(PLUGIN_NAME, store.MODE_LOCAL, Init)
}

func Init(addressStr string) (store.MConfigStore, error) {
	exists, err := file.Exists(addressStr)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = os.Mkdir(addressStr, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	config, err := leveldb.OpenFile(addressStr+"/config", nil)
	filter, err := leveldb.OpenFile(addressStr+"/filter", nil)
	if err != nil {
		return nil, err
	}
	fileStore := FileStore{
		watchChan: make(chan *ConfigEvent, SIZE_CHANGEEVENTBUS),
		config:    config,
		filter:    filter,
	}
	return &fileStore, nil
}
