package file

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"log"
	"testing"
)

func TestIterator(t *testing.T) {
	db, err := leveldb.OpenFile("test", nil)
	if err != nil {
		t.Fatal(err)
	}
	prefix1 := "prefix1/"
	prefix2 := "prefix2/"
	for i := 0; i < 100; i++ {
		db.Put([]byte(prefix1+fmt.Sprintf("%d", i)), []byte(prefix1+fmt.Sprintf("%d", i)), nil)
		db.Put([]byte(prefix2+fmt.Sprintf("%d", i)), []byte(prefix2+fmt.Sprintf("%d", i)), nil)
	}
	customerRange := &util.Range{
		Start: []byte(prefix1),
		Limit: []byte(prefix1 + "z"),
	}
	iterator := db.NewIterator(customerRange, nil)
	for iterator.Next() {
		log.Println(fmt.Sprintf("%s", iterator.Value()))
	}
}
