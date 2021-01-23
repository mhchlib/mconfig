package client

import (
	"github.com/mhchlib/mconfig/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddClient01(t *testing.T) {
	_, err := NewClient(&MetaData{}, nil, nil)
	assert.NotNil(t, err)
}

func TestAddClient02(t *testing.T) {
	InitClientManagement()
	_, err := NewClient(&MetaData{}, nil, nil)
	assert.Nil(t, err)
}

func TestBuildClientConfigRelation(t *testing.T) {
	InitClientManagement()
	client, err := NewClient(&MetaData{}, nil, nil)
	assert.Nil(t, err)
	err = client.BuildClientConfigRelation(test.MockAppkey(), test.MockConfigkeys(100), "")
	assert.Nil(t, err)
}

func TestRemoveClient01(t *testing.T) {
	InitClientManagement()
	client, err := NewClient(&MetaData{}, nil, nil)
	assert.Nil(t, err)
	err = client.BuildClientConfigRelation(test.MockAppkey(), test.MockConfigkeys(100), "")
	assert.Nil(t, err)
	err = client.RemoveClient()
	assert.Nil(t, err)
}
