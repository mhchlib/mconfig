package client

import (
	"github.com/mhchlib/mconfig/test"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestAddClient01(t *testing.T) {
	_, err := NewClient(&MetaData{})
	assert.NotNil(t, err)
}

func TestAddClient02(t *testing.T) {
	InitClientManagement()
	_, err := NewClient(&MetaData{})
	assert.Nil(t, err)
}

func TestBuildClientConfigRelation(t *testing.T) {
	InitClientManagement()
	client, err := NewClient(&MetaData{})
	assert.Nil(t, err)
	err = client.BuildClientConfigRelation(test.MockAppkey(), test.MockConfigkeys(100), "")
	assert.Nil(t, err)
}

func TestRemoveClient01(t *testing.T) {
	InitClientManagement()
	client, err := NewClient(&MetaData{})
	assert.Nil(t, err)
	err = client.BuildClientConfigRelation(test.MockAppkey(), test.MockConfigkeys(100), "")
	assert.Nil(t, err)
	err = client.RemoveClient()
	assert.Nil(t, err)
}

func TestRemoveClient02(t *testing.T) {
	InitClientManagement()
	appKey := test.MockAppkey()
	configKeys := test.MockConfigkeys(100)
	client, err := NewClient(&MetaData{})
	assert.Nil(t, err)
	err = client.BuildClientConfigRelation(appKey, configKeys, "")
	assert.Nil(t, err)
	err = client.RemoveClient()
	assert.Nil(t, err)
	set := GetOnlineClientSet(appKey, configKeys[3], "")
	assert.False(t, set.contains(client))
}

func TestGetClientSet01(t *testing.T) {
	InitClientManagement()
	appKey := test.MockAppkey()
	configKeys := test.MockConfigkeys(100)
	client, err := NewClient(&MetaData{})
	assert.Nil(t, err)
	err = client.BuildClientConfigRelation(appKey, configKeys, "")
	assert.Nil(t, err)
	set := GetOnlineClientSet(appKey, configKeys[3], "")
	assert.True(t, set.contains(client))
}

func TestGetClientSet02(t *testing.T) {
	count := 100
	InitClientManagement()
	clients := make([]*Client, 0)
	appKey := test.MockAppkey()
	configKeys := test.MockConfigkeys(100)

	for i := 0; i < count; i++ {
		client, err := NewClient(&MetaData{})
		assert.Nil(t, err)
		err = client.BuildClientConfigRelation(appKey, configKeys, "")
		assert.Nil(t, err)
		clients = append(clients, client)
	}
	set := GetOnlineClientSet(appKey, configKeys[3], "")
	for _, client := range clients {
		assert.True(t, set.contains(client))
	}
}

func TestGetClientSet03(t *testing.T) {
	var swg sync.WaitGroup
	count := 100
	InitClientManagement()
	appKey := test.MockAppkey()
	configKeys := test.MockConfigkeys(100)
	swg.Add(count)
	for i := 0; i < count; i++ {
		go func(group *sync.WaitGroup) {
			client, err := NewClient(&MetaData{})
			assert.Nil(t, err)
			err = client.BuildClientConfigRelation(test.MockAppkey(), test.MockConfigkeys(10), "")
			assert.Nil(t, err)
			group.Done()
		}(&swg)
	}
	swg.Wait()
	client, err := NewClient(&MetaData{})
	assert.Nil(t, err)
	err = client.BuildClientConfigRelation(appKey, configKeys, "")
	assert.Nil(t, err)
	set := GetOnlineClientSet(appKey, configKeys[3], "")
	assert.True(t, set.contains(client))
}
