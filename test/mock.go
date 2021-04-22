package test

import (
	"github.com/mhchlib/mconfig/core/mconfig"
	"math/rand"
	"strconv"
)

// MockAppkey ...
func MockAppkey() mconfig.AppKey {
	return mconfig.AppKey("appkey_" + strconv.Itoa(rand.Intn(10000)))
}

// MockConfigkey ...
func MockConfigkey() mconfig.ConfigKey {
	return mconfig.ConfigKey("configkey_" + strconv.Itoa(rand.Intn(10000)))
}

// MockConfigkeys ...
func MockConfigkeys(n int) []mconfig.ConfigKey {
	configKeys := []mconfig.ConfigKey{}
	for i := 0; i < n; i++ {
		configKeys = append(configKeys, MockConfigkey())
	}
	return configKeys
}
