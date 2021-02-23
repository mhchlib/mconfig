module github.com/mhchlib/mconfig

go 1.14

require (
	github.com/ChenHaoHu/ExpressionParser v0.0.0-20200730123550-c11f86762d52
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/golang/protobuf v1.4.0
	github.com/google/uuid v1.1.2 // indirect
	github.com/hashicorp/golang-lru v0.5.3
	github.com/mhchlib/logger v0.0.0-20201023050446-420de20374cc
	github.com/mhchlib/mconfig-api v0.0.0-20201023050446-420de20374cc
	github.com/mhchlib/register v0.0.0-20201023050446-420de20374cc
	github.com/mitchellh/mapstructure v1.1.2
	github.com/mkevac/debugcharts v0.0.0-20191222103121-ae1c48aa8615
	github.com/stretchr/testify v1.4.0
	github.com/yuin/gopher-lua v0.0.0-20200816102855-ee81675732da
	go.etcd.io/etcd v3.3.25+incompatible // indirect
	golang.org/x/net v0.0.0-20201016165138-7b1cca2348c0 // indirect
	google.golang.org/grpc v1.26.0
)

replace github.com/mhchlib/logger v0.0.0-20201023050446-420de20374cc => ../logger

replace github.com/mhchlib/mconfig-api v0.0.0-20201023050446-420de20374cc => ../mconfig-api

replace github.com/mhchlib/register v0.0.0-20201023050446-420de20374cc => ../register
