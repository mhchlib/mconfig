package client

import (
	"github.com/mhchlib/mconfig/test"
	"sync"
	"testing"
)

func TestAddClient(t *testing.T) {
	InitClientManagement()
	_, err := NewClient(&MetaData{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestBuildClientConfigRelation(t *testing.T) {
	InitClientManagement()
	client, err := NewClient(&MetaData{})
	if err != nil {
		t.Fatal(err)
	}
	err = client.BuildClientConfigRelation(test.MockAppkey(), test.MockConfigkeys(100))
	if err != nil {
		t.Fatal(err)
	}
}

func TestRemoveClient01(t *testing.T) {
	InitClientManagement()
	client, err := NewClient(&MetaData{})
	if err != nil {
		t.Fatal(err)
	}
	err = client.BuildClientConfigRelation(test.MockAppkey(), test.MockConfigkeys(100))
	if err != nil {
		t.Fatal(err)
	}
	err = client.RemoveClient()
	if err != nil {
		t.Fatal(err)
	}
}

func TestRemoveClient02(t *testing.T) {
	InitClientManagement()
	appKey := test.MockAppkey()
	configKeys := test.MockConfigkeys(100)
	client, err := NewClient(&MetaData{})
	if err != nil {
		t.Fatal(err)
	}
	err = client.BuildClientConfigRelation(appKey, configKeys)
	if err != nil {
		t.Fatal(err)
	}
	err = client.RemoveClient()
	if err != nil {
		t.Fatal(err)
	}
	set := GetOnlineClientSet(appKey, configKeys[3])
	if set.contains(*client) {
		t.Fatal()
	}
}

func TestGetClientSet01(t *testing.T) {
	InitClientManagement()
	appKey := test.MockAppkey()
	configKeys := test.MockConfigkeys(100)
	client, err := NewClient(&MetaData{})
	if err != nil {
		t.Fatal(err)
	}
	err = client.BuildClientConfigRelation(appKey, configKeys)
	if err != nil {
		t.Fatal(err)
	}
	set := GetOnlineClientSet(appKey, configKeys[3])
	if set != nil {
		if !set.contains(*client) {
			t.Fatal("lose client data")
		}
	} else {
		t.Fatal("no client set")
	}
}

func TestGetClientSet02(t *testing.T) {
	count := 100
	InitClientManagement()
	clients := make([]Client, 0)
	appKey := test.MockAppkey()
	configKeys := test.MockConfigkeys(100)

	for i := 0; i < count; i++ {
		client, err := NewClient(&MetaData{})
		if err != nil {
			t.Fatal(err)
		}
		err = client.BuildClientConfigRelation(appKey, configKeys)
		if err != nil {
			t.Fatal(err)
		}
		clients = append(clients, *client)
	}
	set := GetOnlineClientSet(appKey, configKeys[3])
	for _, client := range clients {
		if !set.contains(client) {
			t.Fatal("lose client with id: ", client.Id)
		}
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
			if err != nil {
				t.Fatal(err)
			}
			err = client.BuildClientConfigRelation(test.MockAppkey(), test.MockConfigkeys(10))
			if err != nil {
				t.Fatal(err)
			}
			group.Done()
		}(&swg)
	}
	swg.Wait()
	client, err := NewClient(&MetaData{})
	if err != nil {
		t.Fatal(err)
	}
	err = client.BuildClientConfigRelation(appKey, configKeys)
	if err != nil {
		t.Fatal(err)
	}
	set := GetOnlineClientSet(appKey, configKeys[2])
	if !set.contains(*client) {
		t.Fatal()
	}
}
