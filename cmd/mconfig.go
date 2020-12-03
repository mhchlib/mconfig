package main

import (
	"flag"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig-api/api/v1/cli"
	"github.com/mhchlib/mconfig-api/api/v1/sdk"
	"github.com/mhchlib/mconfig/pkg"
	"github.com/mhchlib/register"
	"github.com/mhchlib/register/common"
	"github.com/mhchlib/register/mregister"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

var mconfig *pkg.MConfig

func init() {
	mconfig = &pkg.MConfig{}
	mconfig.Namspace = flag.String("namespace", "com.github.mhchlib", "Input Your Namespace")
	mconfig.RegistryAddress = flag.String("registry_address", "127.0.0.1:2389", "Input Your Registry Address, multiple IP commas separate")
	mconfig.RegistryType = flag.String("registry_type", "etcd", "Input Your Registry Type, Such etcd ,...(support more soon)")
	mconfig.StoreAddress = flag.String("store_address", "127.0.0.1:2389", "Input Your Store Address, multiple IP commas separate")
	mconfig.StoreType = flag.String("store_type", "etcd", "Input Your Store Type, Such etcd ,...(support more soon)")
	mconfig.ServerIp = flag.String("server_ip", "", "Input Your Server Ip, default local ip")
	mconfig.ServerPort = flag.Int("port", 8080, "Input Your Server Ip")
	flag.Parse()
}

func main() {
	done := make(chan os.Signal, 1)
	defer pkg.InitMconfig(mconfig)()
	reg, err := register.InitRegister(*mconfig.RegistryType, func(options *mregister.Options) {
		options.Address = strings.Split(*mconfig.RegistryAddress, ",")
		options.NameSpace = *mconfig.Namspace
		if *mconfig.ServerIp == "" {
			ip, err := common.GetClientIp()
			if err != nil {
				log.Fatal("get client ip error")
			}
			*mconfig.ServerIp = ip
		}
		options.ServerInstance = *mconfig.ServerIp + ":" + strconv.Itoa(*mconfig.ServerPort)
	})
	if err != nil {
		log.Fatal(err)
	}
	_ = reg.RegisterService("mconfig-cli")
	_ = reg.RegisterService("mconfig-sdk")
	defer func() {
		_ = reg.UnRegisterService("mconfig-sdk")
		_ = reg.UnRegisterService("mconfig-cli")
	}()
	listener, err := net.Listen("tcp", "0.0.0.0"+":"+strconv.Itoa(*mconfig.ServerPort))
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer()
	defer func() {
		_ = listener.Close()
		server.Stop()
	}()
	sdk.RegisterMConfigServer(server, pkg.NewMConfigSDK())
	cli.RegisterMConfigCliServer(server, pkg.NewMConfigCLI())
	go func() {
		err = server.Serve(listener)
		if err != nil {
			log.Fatal(err)
		}
	}()
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
