package pkg

import (
	"github.com/go-acme/lego/v3/log"
	"testing"
)

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
		ok, err := CheckConfigSchema(&Config{
			Config: examples[i][0],
			Schema: examples[i][1],
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Println(ok)
	}
}
