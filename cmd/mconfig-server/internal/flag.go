package internal

import (
	"errors"
	"flag"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/core/mconfig"
	"net/http"
	"strconv"
	"strings"
)

type MconfigFlag struct {
	Namspace    *string
	RegistryStr *string
	StoreStr    *string
	ExposeStr   *string
	EnableDebug *bool
}

var (
	ConfigSeparateSymbol = "://"
	DefaultExposePort    = 8080
	DefaultExposeIp      = ""
)

func NewMconfigFlag() *MconfigFlag {
	return &MconfigFlag{}
}

func ParseFlag(mconfig *mconfig.MConfigConfig) {
	mconfigFlag := initFlagConfig()
	flag.Parse()
	err := parseFlagData(mconfigFlag, mconfig)
	if err != nil {
		log.Fatal(err)
	}
}

func parseFlagData(mconfigFlag *MconfigFlag, mconfig *mconfig.MConfigConfig) error {
	//namespace
	mconfig.Namspace = *mconfigFlag.Namspace
	//registry
	if *mconfigFlag.RegistryStr != "" {
		mconfig.EnableRegistry = true
		registryStr := *mconfigFlag.RegistryStr
		registerType, registerAddress, err := parseAddressFlag(registryStr)
		if err != nil {
			return err
		}
		mconfig.RegistryAddress = registerAddress
		mconfig.RegistryType = registerType
	}
	//store
	if *mconfigFlag.StoreStr != "" {
		storeStr := *mconfigFlag.StoreStr
		storeSType, storeAddress, err := parseAddressFlag(storeStr)
		if err != nil {
			return err
		}
		mconfig.StoreType = storeSType
		mconfig.StoreAddress = storeAddress
	}
	//expose
	mconfig.ServerIp = DefaultExposeIp
	mconfig.ServerPort = DefaultExposePort
	if *mconfigFlag.ExposeStr != "" {
		ip, port, err := parseExposeFlag(*mconfigFlag.ExposeStr)
		if err != nil {
			return err
		}
		mconfig.ServerIp = ip
		mconfig.ServerPort = port
	}
	//debug
	if *mconfigFlag.EnableDebug {
		//debug
		//-------------
		go func() {
			log.Info(http.ListenAndServe("localhost:6060", nil))
		}()
		log.Info("you can now open http://localhost:6060/debug/charts/ in your browser for debug, support ppprof")
		//------------------
	}

	return nil
}

func parseExposeFlag(exposeStr string) (string, int, error) {
	splits := strings.Split(exposeStr, ":")
	if len(splits) != 2 {
		return "", 0, errors.New(exposeStr + " is invalid Expose Address")
	}
	ip := splits[0]
	port, err := strconv.Atoi(splits[1])
	if err != nil {
		return "", 0, errors.New(exposeStr + " is invalid Expose Address")
	}
	return ip, port, nil
}

func parseAddressFlag(str string) (string, string, error) {
	splits := strings.Split(str, ConfigSeparateSymbol)
	if len(splits) != 2 {
		return "", "", errors.New(str + " is invalid Address")
	}
	return splits[0], splits[1], nil
}

func initFlagConfig() *MconfigFlag {
	mconfigFlag := NewMconfigFlag()
	mconfigFlag.Namspace = flag.String("namespace", "com.github.mhchlib", "input your namespace")
	mconfigFlag.RegistryStr = flag.String("registry", "", "input registry address like etcd://127.0.0.1:2389")
	mconfigFlag.StoreStr = flag.String("store", "file://mconfigData/", "input store address like file://mconfigData/")
	mconfigFlag.ExposeStr = flag.String("expose", ":8080", "input server ip, default local ip")
	mconfigFlag.EnableDebug = flag.Bool("debug", false, "enable debug mode")
	return mconfigFlag
}
