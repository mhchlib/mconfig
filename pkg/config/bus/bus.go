package bus

import "github.com/mhchlib/mconfig/pkg/mconfig"

type ConfigChangeItem struct {
	AppKey    mconfig.Appkey
	ConfigKey mconfig.ConfigKey
}

var bus chan ConfigChangeItem

const LENGTH_MAX_BUS = 20

func init() {
	bus = make(chan ConfigChangeItem, LENGTH_MAX_BUS)
}

func AddConfigChange(item ConfigChangeItem) {
	bus <- item
}

func GetConfigChangeBus() chan ConfigChangeItem {
	return bus
}
