package main

import (
	log "github.com/mhchlib/logger"
	mconfig2 "github.com/mhchlib/mconfig"
	"github.com/mhchlib/mconfig-api/api/v1/cli"
	"github.com/mhchlib/mconfig-api/api/v1/sdk"
	"github.com/mhchlib/mconfig/cmd/mconfig/internal"
	"github.com/mhchlib/mconfig/pkg"
	"github.com/mhchlib/mconfig/pkg/rpc"
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

var mconfig *mconfig2.MConfig

func init() {
	mconfig = mconfig2.NewMConfig()
	internal.ShowBanner()
	internal.InitFlag(mconfig)
}

func main() {
	done := make(chan os.Signal, 1)
	defer pkg.InitMconfig(mconfig)()
	if *mconfig.EnableRegistry {
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
	}
	listener, err := net.Listen("tcp", "0.0.0.0"+":"+strconv.Itoa(*mconfig.ServerPort))
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer()
	defer func() {
		_ = listener.Close()
		server.Stop()
	}()
	sdk.RegisterMConfigServer(server, rpc.NewMConfigSDK())
	cli.RegisterMConfigCliServer(server, rpc.NewMConfigCLI())
	go func() {
		err = server.Serve(listener)
		if err != nil {
			log.Error(err)
			done <- syscall.SIGTERM
			return
		}
	}()
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
