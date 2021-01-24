package main

import (
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig-api/api/v1/sdk"
	"github.com/mhchlib/mconfig/cmd/mconfig-server/internal"
	"github.com/mhchlib/mconfig/pkg"
	mconfig "github.com/mhchlib/mconfig/pkg/mconfig"
	"github.com/mhchlib/mconfig/pkg/rpc"
	_ "github.com/mhchlib/mconfig/pkg/store/plugin/etcd"
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

var m *mconfig.MConfig

func init() {
	m = mconfig.NewMConfig()
	internal.ShowBanner()
	internal.ParseFlag(m)
}

func main() {
	done := make(chan os.Signal, 1)
	defer pkg.InitMconfig(m)()
	if m.EnableRegistry {
		reg, err := register.InitRegister(m.RegistryType, func(options *mregister.Options) {
			options.Address = strings.Split(m.RegistryAddress, ",")
			options.NameSpace = m.Namspace
			if m.ServerIp == "" {
				ip, err := common.GetClientIp()
				if err != nil {
					log.Fatal("get client ip error")
				}
				m.ServerIp = ip
			}
			options.ServerInstance = m.ServerIp + ":" + strconv.Itoa(m.ServerPort)
		})
		if err != nil {
			log.Fatal(err)
		}
		err = reg.RegisterService("mconfig-server-cli")
		if err != nil {
			log.Fatal(err)
		}
		err = reg.RegisterService("mconfig-server-sdk")
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			err = reg.UnRegisterService("mconfig-server-sdk")
			if err != nil {
				log.Error(err)
			}
			err = reg.UnRegisterService("mconfig-server-cli")
			if err != nil {
				log.Error(err)
			}
		}()
	}
	listener, err := net.Listen("tcp", "0.0.0.0"+":"+strconv.Itoa(m.ServerPort))
	log.Info("mconfig-server listen :" + strconv.Itoa(m.ServerPort) + " success")

	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer()
	defer func() {
		_ = listener.Close()
		server.Stop()
	}()
	sdk.RegisterMConfigServer(server, rpc.NewMConfigSDK())
	//cli.RegisterMConfigCliServer(server, rpc.NewMConfigCLI())
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
