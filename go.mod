module github.com/mhchlib/mconfig

go 1.14

require (
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/deckarep/golang-set v1.7.1
	github.com/go-acme/lego/v3 v3.4.0
	github.com/golang/protobuf v1.4.0
	github.com/google/uuid v1.1.2 // indirect
	github.com/mhchlib/logger v0.0.0-20201023050446-420de20374cc
	github.com/mhchlib/mconfig-api v0.0.0-20201023050446-420de20374cc
	github.com/mhchlib/register v0.0.0-20201023050446-420de20374cc
	github.com/micro/go-micro/v2 v2.9.1
	github.com/syndtr/goleveldb v1.0.0
	github.com/xeipuuv/gojsonschema v1.1.0
	go.etcd.io/etcd v3.3.25+incompatible // indirect
	golang.org/x/net v0.0.0-20201016165138-7b1cca2348c0 // indirect
	google.golang.org/grpc v1.26.0
)

replace github.com/mhchlib/logger v0.0.0-20201023050446-420de20374cc => ../logger

replace github.com/mhchlib/mconfig-api v0.0.0-20201023050446-420de20374cc => ../mconfig-api

replace github.com/mhchlib/register v0.0.0-20201023050446-420de20374cc => ../register
