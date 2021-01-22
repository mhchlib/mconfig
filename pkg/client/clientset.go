package client

import "sync"

type ClientSet struct {
	sync.RWMutex
	m map[ClientId]*Client
}

func (set *ClientSet) add(client Client) error {
	set.Lock()
	set.m[client.Id] = &client
	set.Unlock()
	return nil
}

func (set *ClientSet) remove(client Client) error {
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

func (set *ClientSet) GetClients() map[ClientId]*Client {
	return set.m
}
