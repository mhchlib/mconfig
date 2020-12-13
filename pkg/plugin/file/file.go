package file

import (
	"context"
	"encoding/json"
	"github.com/mhchlib/mconfig/pkg"
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
	pkg.RegisterStorePlugin("file", Init)
}

func Init(address string) (pkg.AppConfigStore, error) {
	exists, err := file.Exists(address)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = os.Mkdir(address, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	db, err := leveldb.OpenFile(address+string(os.PathSeparator)+DB_NAME, nil)
	if err != nil {
		return nil, err
	}
	fileStore := FileStore{
		DB:             db,
		ChangeEventBus: make(chan *pkg.ConfigEvent, SIZE_ChangeEventBus),
	}
	return fileStore, nil
}

func (f FileStore) GetAppConfigs(key pkg.Appkey) (*pkg.AppConfigs, error) {
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

func (f FileStore) PutAppConfigs(key pkg.Appkey, value *pkg.AppConfigs) error {
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
