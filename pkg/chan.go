package pkg

import (
	"sync"
)

type ClientChanMap struct {
	sync.RWMutex
	m map[Appkey]map[ClientId]chan interface{}
}

var (
	//用于收到store配置改变事件后通知更新客户端
	configChangeChan chan Appkey
	clientChanMap    *ClientChanMap
)

func init() {
	configChangeChan = make(chan Appkey, 10)
	clientChanMap = &ClientChanMap{
		m: make(map[Appkey]map[ClientId]chan interface{}),
	}
}
