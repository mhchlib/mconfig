package internal

import (
	"errors"
	"flag"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/pkg"
	"strings"
)

type MconfigFlag struct {
	Namspace    *string
	RegistryStr *string
	StoreStr    *string
	ExposeStr   *string
}

var ConfigSeparateSymbol = "://"

func NewMconfigFlag() *MconfigFlag {
	return &MconfigFlag{}
}

func ParseFlag(mconfig *pkg.MConfig) {
	mconfigFlag := initFlagConfig()
	flag.Parse()
	err := parseFlagData(mconfigFlag, mconfig)
	if err != nil {
		log.Fatal(err)
	}
}

func parseFlagData(mconfigFlag *MconfigFlag, mconfig *pkg.MConfig) error {
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
	return nil
}

func parseAddressFlag(str string) (string, string, error) {
	splits := strings.Split(str, ConfigSeparateSymbol)
	if len(splits) != 2 {
		return "", "", errors.New(str + " is invalid")
	}
	return splits[0], splits[1], nil
}

func initFlagConfig() *MconfigFlag {
	mconfigFlag := NewMconfigFlag()
	mconfigFlag.Namspace = flag.String("namespace", "com.github.mhchlib", "input your namespace")
	mconfigFlag.RegistryStr = flag.String("registry", "", "input registry address like etcd://127.0.0.1:2389")
	mconfigFlag.StoreStr = flag.String("store", "file://t_file/", "input store address like file://t_file/")
	mconfigFlag.ExposeStr = flag.String("expose", "0.0.0.0:8080", "input server ip, default local ip")
	return mconfigFlag
}
