
build:
	go build cmd/mconfig.go

dev:
	go run cmd/mconfig.go cmd/plugin.go --registry=etcd --registry_address=etcd.u.hcyang.top:31770

example01:
	go run example/continuous/main.go --registry=etcd --registry_address=etcd.u.hcyang.top:31770

example02:
	go run example/concurrent/main.go --registry=etcd --registry_address=etcd.u.hcyang.top:31770



.PHONY: example