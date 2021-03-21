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
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
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
	err, closeFunc := initRegister()
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

func initRegister() (error, func()) {
	//register service to register center
	if m.RegistryAddress != "" {
		regClient, err := register.InitRegister(
			register.Namespace(m.Namspace),
			register.ResgisterAddress(m.RegistryAddress),
			register.Instance(m.ServerIp+":"+strconv.Itoa(m.ServerPort)),
			register.Metadata("mode", store.GetStorePlugin().Mode),
		)
		if err != nil {
			return err, nil
		}
		m.RegistryType = string(regClient.RegisterType)
		demandSync := store.CheckNeedSyncData()
		if demandSync {
			err := store.SyncOtherMconfigData(regClient.Srv, SERVICE_NAME)
			if err != nil {
				return errors.New("sync store data fail:" + err.Error()), nil
			}
		}

		unRegisterFunc, err := regClient.Srv.RegisterService(SERVICE_NAME, nil)
		if err != nil {
			return err, nil
		}
		return nil, func() {
			unRegisterFunc()
		}
	}
	return nil, func() {}
}

func printMconfigDetail() {
	data := [][]string{
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
	})
}
