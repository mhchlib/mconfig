package file

import (
	"context"
	"encoding/json"
	"github.com/mhchlib/mconfig/pkg"
	"github.com/mhchlib/mconfig/pkg/config"
	"github.com/mhchlib/mconfig/pkg/mconfig"
	"github.com/mhchlib/mconfig/pkg/store"
	"github.com/micro/go-micro/v2/util/file"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
)

const DB_NAME = "mconfig"
const SIZE_ChangeEventBus = 10

type FileStore struct {
	DB             *leveldb.DB
	ChangeEventBus chan *pkg.ConfigEvent
}

func init() {
	store.RegisterStorePlugin("file", Init)
}

func Init(filePath string) (pkg.AppConfigStore, error) {
	exists, err := file.Exists(filePath)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = os.Mkdir(filePath, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	db, err := leveldb.OpenFile(filePath+string(os.PathSeparator)+DB_NAME, nil)
	if err != nil {
		return nil, err
	}
	fileStore := FileStore{
		DB:             db,
		ChangeEventBus: make(chan *pkg.ConfigEvent, SIZE_ChangeEventBus),
	}
	return fileStore, nil
}

func (f FileStore) GetAppConfigs(key mconfig.Appkey) (*config.AppConfigs, error) {
	db := f.DB
	value, err := db.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}
	appConfigs, err := parseAppConfigsJSONStr(AppConfigsJSONStr(value))
	if err != nil {
		return nil, err
	}
	//log.Info("get app config ",key,string(value))
	return appConfigs, nil
}

func (f FileStore) PutAppConfigs(key mconfig.Appkey, value *config.AppConfigs) error {
	configJsonByte, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = f.DB.Put([]byte(key), configJsonByte, nil)
	if err != nil {
		return err
	}
	f.ChangeEventBus <- &pkg.ConfigEvent{
		Key:        key,
		AppConfigs: value,
		EventType:  pkg.Event_Update,
	}
	//log.Info("put app config ",key,string(configJsonByte))
	return nil
}

func (f FileStore) WatchAppConfigs(ctx context.Context) (chan *pkg.ConfigEvent, error) {
	configChan := make(chan *pkg.ConfigEvent)
	go func() {
		for {
			select {
			case v := <-f.ChangeEventBus:
				configChan <- v
			}
		}
	}()
	return configChan, nil
}
