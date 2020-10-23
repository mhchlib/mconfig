package service

import (
	"log"
	"testing"
)

func TestParseConfigJSONStr(t *testing.T) {
	v := `[
		{"status":1,"key":"1111","value":"121211","schema":"1aaaa","objects":"dsadda"},
		{"status":2,"key":"333","value":"wwww","schema":"ewew","objects":"dsad111da"},
		{"status":1,"key":"4444","value":"ssss","schema":"ewwe","objects":"dsadda"},
		{}
	]`
	vs, err := ParseConfigJSONStr(ConfigJSONStr(v))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", vs)
}

func TestCheckConfigsSchema(t *testing.T) {

}

func TestCheckConfigSchema(t *testing.T) {

	examples := [][]string{
		{`{"a":"b"}`, `{"type": "object","properties":{"a":{"type":"string"}}}`},
		{`{"age":1111}`, `{"type":"object","properties":{"age":{"type":"integer"}}}`},
		{`{"age":1111,"name":"1111"}`, `{"type":"object","properties":{"age":{"type":"integer"}}}`},
		{``, ``},
	}
	for i := 0; i < len(examples); i++ {
		config := &Mconfig{
			Value:  examples[i][0],
			Schema: examples[i][1],
		}
		ok, err := CheckConfigSchema(config)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(ok)
	}
}
