module github.com/mhchlib/mconfig

go 1.14

require (
	github.com/ChenHaoHu/ExpressionParser v0.0.0-20200730123550-c11f86762d52
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.2 // indirect
	github.com/hashicorp/golang-lru v0.5.3
	github.com/mhchlib/logger v0.0.0-20201023050446-420de20374cc
	github.com/mhchlib/mconfig-api v0.0.0-20201023050446-420de20374cc
	github.com/mhchlib/register v0.0.0-20201023050446-420de20374cc
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.9.1
	github.com/mitchellh/mapstructure v1.1.2
	github.com/mkevac/debugcharts v0.0.0-20191222103121-ae1c48aa8615
	github.com/pyroscope-io/pyroscope v0.0.26
	github.com/stretchr/testify v1.4.0
	github.com/syndtr/goleveldb v1.0.0
	github.com/yuin/gopher-lua v0.0.0-20200816102855-ee81675732da
	go.etcd.io/etcd v3.3.25+incompatible // indirect
	google.golang.org/grpc v1.27.0
)

replace github.com/mhchlib/logger v0.0.0-20201023050446-420de20374cc => ../logger

replace github.com/mhchlib/mconfig-api v0.0.0-20201023050446-420de20374cc => ../mconfig-api

replace github.com/mhchlib/register v0.0.0-20201023050446-420de20374cc => ../register

//replace go.etcd.io/etcd v3.3.25+incompatible => go.etcd.io/etcd v3.3.25+incompatible
//
//replace github.com/coreos/etcd v3.3.25+incompatible => go.etcd.io/etcd v3.3.25+incompatible

replace google.golang.org/grpc v1.27.0 => google.golang.org/grpc v1.26.0
