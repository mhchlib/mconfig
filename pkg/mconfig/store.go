package mconfig

type StoreVal struct {
	DataVersion
	Data interface{} `json:"data"`
}

func buildStoreVal(val interface{}) (*StoreVal, error) {
	storeVal := &StoreVal{
		DataVersion: DataVersion{
			Md5:     GetInterfaceMd5(val),
			Version: createDataVersion(),
		},
		Data: val,
	}

	return storeVal, nil
}
