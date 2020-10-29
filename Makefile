
build:
	go build cmd/mconfig.go

dev:
	go run cmd/mconfig.go cmd/plugin.go --registry=etcd --registry_address=127.0.0.1:2379