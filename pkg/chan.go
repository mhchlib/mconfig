package pkg

import (
	"sync"
)

type ClientChanMap struct {
	sync.RWMutex
	m map[AppId]map[ClientId]chan interface{}
}

var (
	//用于收到store配置改变事件后通知更新客户端
	configChangeChan chan AppId
	clientChanMap    *ClientChanMap
)

func init() {
	configChangeChan = make(chan AppId, 10)
	clientChanMap = &ClientChanMap{
		m: make(map[AppId]map[ClientId]chan interface{}),
	}
}
