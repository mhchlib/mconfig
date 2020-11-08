package pkg

import (
	"github.com/go-acme/lego/v3/log"
	"testing"
)

func TestParseConfigJSONStr(t *testing.T) {
	v := `[
		{"id":"1222","config":"2121","schema":"dadsa","create_time":1604249335,"update_time":1604249335,"desc":"dasda","status":0},
		{"id":"1223","config":"2121","schema":"dadsa","create_time":1604249335,"update_time":1604249335,"desc":"dasda","status":0},
		{"id":"1224","config":"2121","schema":"dadsa","create_time":1604249335,"update_time":1604249335,"desc":"dasda","status":0},
		{"id":"1225","config":"2121","schema":"dadsa","create_time":1604249335,"update_time":1604249335,"desc":"dasda","status":0},
		{"id":"1226","config":"2121","schema":"dadsa","create_time":1604249335,"update_time":1604249335,"desc":"dasda","status":0},
		{}
	]`
	vs, err := ParseConfigJSONStr(ConfigJSONStr(v))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", vs)
}

func TestCheckConfigSchema(t *testing.T) {
	examples := [][]string{
		{`{"a":111}`, `{"type": "object","properties":{"a":{"type":"string"}}}`},
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
