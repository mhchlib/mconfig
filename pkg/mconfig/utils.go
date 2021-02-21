package mconfig

import (
	"crypto/md5"
	"fmt"
	"time"
)

func createDataVersion() int64 {
	return time.Now().UnixNano()
}

func GetInterfaceMd5(d interface{}) string {
	str := fmt.Sprintf("%v", d)
	return fmt.Sprintf("%x", (md5.Sum([]byte(str))))
}
