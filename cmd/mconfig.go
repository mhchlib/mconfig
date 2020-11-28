package main

import (
	"errors"
	"flag"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig-api/api/v1/cli"
	"github.com/mhchlib/mconfig-api/api/v1/sdk"
	"github.com/mhchlib/mconfig/pkg"
	"github.com/mhchlib/register"
	"github.com/mhchlib/register/common"
	etcd_kit "github.com/mhchlib/register/etcd-kit"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

var namspace = flag.String("namespace", "com.github.mhchlib", "Input Your Namespace")
var registry_address = flag.String("registry_address", "127.0.0.1:2389", "Input Your Registry Address, multiple IP commas separate")
var registry_type = flag.String("registry_type", "etcd", "Input Your Registry Type, Such etcd ,...(support more soon)")
var store_address = flag.String("store_address", "127.0.0.1:2389", "Input Your Store Address, multiple IP commas separate")
var store_type = flag.String("store_type", "etcd", "Input Your Store Type, Such etcd ,...(support more soon)")
var server_ip = flag.String("server_ip", "", "Input Your Server Ip, default local ip")
var server_port = flag.Int("port", 8080, "Input Your Server Ip")

func init() {
	flag.Parse()
}

func main() {
	done := make(chan os.Signal, 1)
	defer pkg.InitMconfig(*store_type, *store_address)()
	reg, err := InitRegister(*registry_type, *registry_address)
	if err != nil {
		log.Fatal(err)
	}

	reg.RegisterService("mconfig-cli")
	reg.RegisterService("mconfig-sdk")
	defer func() {
		reg.UnRegisterService("mconfig-sdk")
	}()

	listener, err := net.Listen("tcp", "0.0.0.0"+":"+strconv.Itoa(*server_port))
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer()
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

func InitRegister(registry_type, registry_address string) (register.Register, error) {
	if registry_type == "etcd" {
		reg := etcd_kit.EtcdRegister{}
		reg.Init(func(options *register.Options) {
			options.Address = strings.Split(registry_address, ",")
			options.NameSpace = *namspace
			if *server_ip == "" {
				ip, err := common.GetClientIp()
				if err != nil {
					log.Fatal("get client ip error")
				}
				*server_ip = ip
			}
			options.ServerInstance = *server_ip + ":" + strconv.Itoa(*server_port)
		})
		return &reg, nil
	} else {
	}
	return nil, errors.New("registry type: " + registry_type + " can not be supported, you can choose: etcd")
}
