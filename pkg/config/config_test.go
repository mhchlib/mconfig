package config

import (
	"github.com/mhchlib/mconfig/pkg/mconfig"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"testing"
)

func TestParseConfigEventMetadata(t *testing.T) {
	m := ConfigEventMetadata{
		AppKey:    "appKey",
		ConfigKey: "configKey",
		Env:       "dev",
		Val:       mconfig.ConfigVal("66666" + strconv.Itoa(rand.Intn(1000))),
	}
	_, err := parseConfigEventMetadata(m)
	assert.Nil(t, err)
}
