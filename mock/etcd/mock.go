package main

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	log "github.com/mhchlib/logger"
	"io/ioutil"
	"time"
)

const PREFIX_CONFIG = "/mconfig/"
const MOCK_DATA_PATH = "/Users/huchenhao/Documents/goproject/github.com/mhchlib/mconfig/mock/data.json"

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"etcd.u.hcyang.top:31770"},
		// Endpoints: []string{"localhost:2379", "localhost:22379", "localhost:32379"}
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	kv := clientv3.NewKV(cli)
	//_ = kv
	data, err := ioutil.ReadFile(MOCK_DATA_PATH)
	if err != nil {
		log.Fatal(err)
	}

	var mockConfigs = make(map[string]interface{}, 0)

	err = json.Unmarshal(data, &mockConfigs)
	if err != nil {
		log.Fatal(err)
	}
	log.Info(mockConfigs)
	for k, v := range mockConfigs {
		//log.Info(k,v)
		bytes, _ := json.Marshal(v)
		kv.Put(context.Background(), PREFIX_CONFIG+k, string(bytes))
	}
}
