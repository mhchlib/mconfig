module github.com/mhchlib/mconfig

go 1.14

require (
	    github.com/coreos/etcd v3.3.18+incompatible
    	github.com/google/uuid v1.1.2 // indirect
    	github.com/micro/go-micro/v2 v2.9.1
    	github.com/xeipuuv/gojsonschema v1.1.0
    	golang.org/x/net v0.0.0-20201016165138-7b1cca2348c0 // indirect
    	github.com/mhchlib/logger v0.0.0-20201023050446-420de20374cc
    	github.com/mhchlib/mconfig-api v0.0.0-20201023050446-420de20374cc
)

replace github.com/mhchlib/logger v0.0.0-20201023050446-420de20374cc => ../logger
replace github.com/mhchlib/mconfig-api v0.0.0-20201023050446-420de20374cc => ../mconfig-api
