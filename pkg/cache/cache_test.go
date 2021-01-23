package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetConfigFromCache(t *testing.T) {
	key := "AAA"
	value := "BBB"

	InitCacheManagement()
	err := PutConfigToCache(key, value)
	assert.Nil(t, err)
	val, err := GetConfigFromCache(key)
	assert.Nil(t, err)
	assert.Equal(t, val, value)
}
