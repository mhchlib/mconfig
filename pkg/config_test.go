package pkg

import (
	"github.com/go-acme/lego/v3/log"
	"testing"
)

//func TestParseAppConfigsJSONStr(t *testing.T) {
//	//v := "{\"1000-100\":{\"configs\":{\"0\":{\"config\":\"{'name':'demo1','age':12}\",\"create_time\":1604249335,\"schema\":\"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}\",\"update_time\":1604249335},\"1\":{\"config\":\"{'name':'demo1','age':12}\",\"create_time\":1604249335,\"schema\":\"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}\",\"update_time\":1604249335},\"2\":{\"config\":\"{'name':'demo1','age':12}\",\"create_time\":1604249335,\"schema\":\"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}\",\"update_time\":1604249335}},\"create_time\":1604249335,\"desc\":\"test\",\"update_time\":1604249335},\"1000-101\":{\"configs\":{\"0\":{\"config\":\"{'name':'demo1','age':12}\",\"create_time\":1604249335,\"schema\":\"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}\",\"update_time\":1604249335},\"1\":{\"config\":\"{'name':'demo1','age':12}\",\"create_time\":1604249335,\"schema\":\"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}\",\"update_time\":1604249335},\"2\":{\"config\":\"{'name':'demo1','age':12}\",\"create_time\":1604249335,\"schema\":\"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}\",\"update_time\":1604249335}},\"create_time\":1604249335,\"desc\":\"test\",\"update_time\":1604249335},\"1000-102\":{\"configs\":{\"0\":{\"config\":\"{'name':'demo1','age':12}\",\"create_time\":1604249335,\"schema\":\"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}\",\"update_time\":1604249335},\"1\":{\"config\":\"{'name':'demo1','age':12}\",\"create_time\":1604249335,\"schema\":\"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}\",\"update_time\":1604249335},\"2\":{\"config\":\"{'name':'demo1','age':12}\",\"create_time\":1604249335,\"schema\":\"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}\",\"update_time\":1604249335}},\"create_time\":1604249335,\"desc\":\"test\",\"update_time\":1604249335},\"1000-103\":{\"configs\":{\"0\":{\"config\":\"{'name':'demo1','age':12}\",\"create_time\":1604249335,\"schema\":\"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}\",\"update_time\":1604249335},\"1\":{\"config\":\"{'name':'demo1','age':12}\",\"create_time\":1604249335,\"schema\":\"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}\",\"update_time\":1604249335},\"2\":{\"config\":\"{'name':'demo1','age':12}\",\"create_time\":1604249335,\"schema\":\"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}\",\"update_time\":1604249335}},\"create_time\":1604249335,\"desc\":\"test\",\"update_time\":1604249335},\"1000-104\":{\"configs\":{\"0\":{\"config\":\"{'name':'demo1','age':12}\",\"create_time\":1604249335,\"schema\":\"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}\",\"update_time\":1604249335},\"1\":{\"config\":\"{'name':'demo1','age':12}\",\"create_time\":1604249335,\"schema\":\"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}\",\"update_time\":1604249335},\"2\":{\"config\":\"{'name':'demo1','age':12}\",\"create_time\":1604249335,\"schema\":\"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}\",\"update_time\":1604249335}},\"create_time\":1604249335,\"desc\":\"test\",\"update_time\":1604249335}}"
//	v := "{\"1000-100\":{\"ABFilters\":{\"ip\":\"192.0.0.1\"},\"configs\":{\"entry\":{\"0\":{\"config\":\"{'name':'demo1','age':12}\",\"create_time\":1604249335,\"schema\":\"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}\",\"update_time\":1604249335},\"1\":{\"config\":\"{'name':'demo1','age':12}\",\"create_time\":1604249335,\"schema\":\"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}\",\"update_time\":1604249335},\"2\":{\"config\":\"{'name':'demo1','age':12}\",\"create_time\":1604249335,\"schema\":\"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}\",\"update_time\":1604249335}}},\"create_time\":1604249335,\"desc\":\"test\",\"update_time\":1604249335},\"1000-101\":{\"ABFilters\":{\"ip\":\"192.0.0.1\"},\"configs\":{\"entry\":{\"0\":{\"config\":\"{'name':'demo1','age':12}\",\"create_time\":1604249335,\"schema\":\"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}\",\"update_time\":1604249335},\"1\":{\"config\":\"{'name':'demo1','age':12}\",\"create_time\":1604249335,\"schema\":\"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}\",\"update_time\":1604249335},\"2\":{\"config\":\"{'name':'demo1','age':12}\",\"create_time\":1604249335,\"schema\":\"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}\",\"update_time\":1604249335}}},\"create_time\":1604249335,\"desc\":\"test\",\"update_time\":1604249335},\"1000-102\":{\"ABFilters\":{\"ip\":\"192.0.0.1\"},\"configs\":{\"entry\":{\"0\":{\"config\":\"{'name':'demo1','age':12}\",\"create_time\":1604249335,\"schema\":\"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}\",\"update_time\":1604249335},\"1\":{\"config\":\"{'name':'demo1','age':12}\",\"create_time\":1604249335,\"schema\":\"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}\",\"update_time\":1604249335},\"2\":{\"config\":\"{'name':'demo1','age':12}\",\"create_time\":1604249335,\"schema\":\"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}\",\"update_time\":1604249335}}},\"create_time\":1604249335,\"desc\":\"test\",\"update_time\":1604249335},\"1000-103\":{\"ABFilters\":{\"ip\":\"192.0.0.1\"},\"configs\":{\"entry\":{\"0\":{\"config\":\"{'name':'demo1','age':12}\",\"create_time\":1604249335,\"schema\":\"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}\",\"update_time\":1604249335},\"1\":{\"config\":\"{'name':'demo1','age':12}\",\"create_time\":1604249335,\"schema\":\"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}\",\"update_time\":1604249335},\"2\":{\"config\":\"{'name':'demo1','age':12}\",\"create_time\":1604249335,\"schema\":\"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}\",\"update_time\":1604249335}}},\"create_time\":1604249335,\"desc\":\"test\",\"update_time\":1604249335},\"1000-104\":{\"ABFilters\":{\"ip\":\"192.0.0.1\"},\"configs\":{\"entry\":{\"0\":{\"config\":\"{'name':'demo1','age':12}\",\"create_time\":1604249335,\"schema\":\"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}\",\"update_time\":1604249335},\"1\":{\"config\":\"{'name':'demo1','age':12}\",\"create_time\":1604249335,\"schema\":\"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}\",\"update_time\":1604249335},\"2\":{\"config\":\"{'name':'demo1','age':12}\",\"create_time\":1604249335,\"schema\":\"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}\",\"update_time\":1604249335}}},\"create_time\":1604249335,\"desc\":\"test\",\"update_time\":1604249335}}"
//	vs, err := parseAppConfigsJSONStr(AppConfigsJSONStr(v))
//	if err != nil {
//		log.Fatal(err)
//	}
//	log.Printf("%+v", vs)
//}

func TestCheckConfigSchema(t *testing.T) {
	examples := [][]string{
		{`{"db":{"url":"127.0.0.1:3306","database":"bookstore","time_out": 20}}`, `{
    "type": "object",
    "properties": {
        "db": {
            "type": "object",
            "properties": {
               "url": {
                  "type": "string"
               },
               "database": {
                  "type": "string"
               } ,
				"time_out": {
                  "type": "integer"
                } 
            }
        }
    }
}`},
		{`{"age":1111}`, `{"type":"object","properties":{"age":{"type":"integer"}}}`},
		{`{"age":1111,"name":"1111"}`, `{"type":"object","properties":{"age":{"type":"integer"}}}`},
		{``, ``},
	}
	for i := 0; i < len(examples); i++ {
		ok, err := CheckConfigSchema(examples[i][0], examples[i][1])
		if err != nil {
			log.Fatal(err)
		}
		log.Println(ok)
	}
}
