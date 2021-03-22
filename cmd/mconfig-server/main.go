package main

import (
	"errors"
	"fmt"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/cmd/mconfig-server/internal"
	"github.com/mhchlib/mconfig/core"
	"github.com/mhchlib/mconfig/core/mconfig"
	"github.com/mhchlib/mconfig/core/store"
	_ "github.com/mhchlib/mconfig/core/store/plugin/etcd"
	"github.com/mhchlib/mconfig/rpc"
	"github.com/mhchlib/register"
	"github.com/olekukonko/tablewriter"
	_ "go.uber.org/automaxprocs"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
)

var m *mconfig.MConfigConfig

func init() {
	m = mconfig.NewMConfig()
	internal.ShowBanner()
	internal.ParseFlag(m)
}

const SERVICE_NAME = "mconfig-server"

func main() {
	//set log level
	log.SetDebugLogLevel()
	done := make(chan os.Signal, 1)
	defer core.InitMconfig(m)()

	listener, err := net.Listen("tcp", "0.0.0.0"+":"+strconv.Itoa(m.ServerPort))
	log.Info(fmt.Sprintf("mconfig-server listen :%d success", m.ServerPort))
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	defer func() {
		s.GracefulStop()
		err = listener.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	rpc.InitRpc(s)
	go func() {
		err = s.Serve(listener)
		if err != nil {
			log.Error(err)
			done <- syscall.SIGTERM
			return
		}
	}()
	closeFunc, err := initRegister()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		closeFunc()
	}()
	//print some useful data with ASCII
	printMconfigDetail()

	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
}

func initRegister() (func(), error) {
	//register service to register center
	if m.RegistryAddress != "" {
		regClient, err := register.InitRegister(
			register.Namespace(m.Namspace),
			register.ResgisterAddress(m.RegistryAddress),
			register.Instance(m.ServerIp+":"+strconv.Itoa(m.ServerPort)),
			register.Metadata("mode", store.GetStorePlugin().Mode),
			register.Metadata("plugin", store.GetStorePlugin().Name),
		)
		if err != nil {
			return nil, err
		}
		m.RegistryType = string(regClient.RegisterType)
		storePluginMode := store.GetStorePluginModel()
		storePluginName := store.GetStorePluginName()
		if storePluginMode == store.MODE_LOCAL {
			//register center only have one instance
			services, err := regClient.Srv.ListAllServices(SERVICE_NAME)
			if err != nil {
				return nil, err
			}
			if len(services) != 0 {
				return nil, errors.New("Store local mode only can register one instance in register center")
			}
		}
		if storePluginMode == store.MODE_SHARE {
			//register center only have one class instance
			services, err := regClient.Srv.ListAllServices(SERVICE_NAME)
			if err != nil {
				return nil, err
			}
			if len(services) != 0 {
				//check one service plugin
				service := services[0]
				if service.Metadata["plugin"] != storePluginName {
					return nil, errors.New("Store share mode only can register one class plugin in register center")
				}
			}
		}

		unRegisterFunc, err := regClient.Srv.RegisterService(SERVICE_NAME, nil)
		if err != nil {
			return nil, err
		}
		return func() {
			unRegisterFunc()
		}, nil
	}
	return func() {}, nil
}

func printMconfigDetail() {
	data := [][]string{
		[]string{"Process Num", fmt.Sprintf("%v", runtime.GOMAXPROCS(0))},
		[]string{"Namespace", m.Namspace},
		[]string{"Store Type", m.StoreType},
		[]string{"Store Address", m.StoreAddress},
		[]string{"Store Mode", fmt.Sprintf("%s", store.GetStorePlugin().Mode)},
		[]string{"Store Plugin", fmt.Sprintf("%s", store.GetStorePlugin().Name)},
		[]string{"Register Type", m.RegistryType},
		[]string{"Register Address", m.RegistryAddress},
		[]string{"Register Server Address", m.ServerIp + ":" + strconv.Itoa(m.ServerPort)},
	}
	headers := []string{"Name", "Val"}
	log.PrintDataTable(data, headers, "print some useful data about mconfig ↓ ↓ ↓ ↓", func(table *tablewriter.Table) {
		table.SetAlignment(1)
	})
}
