package client

import (
	"fmt"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/core/mconfig"
	"sync"
)

// ClientSet ...
type ClientSet struct {
	sync.RWMutex
	m map[ClientId]*Client
}

// NewClientSet ...
func NewClientSet() *ClientSet {
	clientSet := &ClientSet{
		m: make(map[ClientId]*Client),
	}
	return clientSet
}

func (set *ClientSet) add(client *Client) error {
	set.Lock()
	set.m[client.Id] = client
	set.Unlock()
	return nil
}

func (set *ClientSet) remove(client *Client) error {
	set.Lock()
	delete(set.m, client.Id)
	set.Unlock()
	return nil
}

func (set *ClientSet) contains(client *Client) bool {
	set.RLock()
	_, ok := set.m[client.Id]
	set.RUnlock()
	return ok
}

func (set *ClientSet) count() int {
	set.RLock()
	defer set.RUnlock()
	return len(set.m)
}

// SendMsg ...
func (set *ClientSet) SendMsg(data *mconfig.ConfigChangeNotifyMsg) error {
	set.RLock()
	defer set.RUnlock()
	for _, c := range set.m {
		err := c.SendConfigChangeNotifyMsg(data)
		if err != nil {
			return err
		}
	}
	return nil
}

// ReCalEffectEnv ...
func (set *ClientSet) ReCalEffectEnv() error {
	set.RLock()
	defer set.RUnlock()
	for _, c := range set.m {
		go func(c *Client) {
			err := c.ReCalEffectEnv()
			if err != nil {
				log.Error(err, "client:", fmt.Sprintf("%v", c))
			}
		}(c)
	}
	return nil
}
