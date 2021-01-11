package pkg

import (
	"sync"
)

// ClientChanMap ...
type ClientChanMap struct {
	sync.RWMutex
	m map[Appkey]map[ClientId]chan interface{}
}

var (
	//用于收到store配置改变事件后通知更新客户端
	ConfigChangeChan chan Appkey
	ClientChans      *ClientChanMap
)

func init() {
	ConfigChangeChan = make(chan Appkey, 10)
	ClientChans = &ClientChanMap{
		m: make(map[Appkey]map[ClientId]chan interface{}),
	}
}
