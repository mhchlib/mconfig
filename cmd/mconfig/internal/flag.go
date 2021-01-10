package internal

import (
	"flag"
	"github.com/mhchlib/mconfig"
)

func InitFlag(mconfig *mconfig.MConfig) {
	mconfig.Namspace = flag.String("namespace", "com.github.mhchlib", "Input Your Namespace")
	mconfig.EnableRegistry = flag.Bool("registry", true, "enable use registry")
	mconfig.RegistryAddress = flag.String("registry_address", "127.0.0.1:2389", "Input Your Registry Address, multiple IP commas separate")
	mconfig.RegistryType = flag.String("registry_type", "etcd", "Input Your Registry Type, Such etcd ,...(support more soon)")
	mconfig.StoreType = flag.String("store_type", "etcd", "Input Your Store Type, Such etcd ,...(support more soon)")
	mconfig.ServerIp = flag.String("server_ip", "", "Input Your Server Ip, default local ip")
	mconfig.ServerPort = flag.Int("port", 8080, "Input Your Server Ip")
	flag.Parse()
}
