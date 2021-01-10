
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -o mconfig cmd/mconfig.go cmd/plugin.go

docker: build
	docker build -t dockerhcy/mconfig:v0.2  .

push: docker
	docker push dockerhcy/mconfig:v0.2

dev:
	go run cmd/mconfig/mconfig.go cmd/mconfig/plugin.go --namespace=local_test --registry_type=etcd --registry_address=etcd.u.hcyang.top:31770 --store_etcd=etcd.u.hcyang.top:31770 --store_type=etcd

example01:
	go run example/continuous/main.go --registry=etcd --registry_address=etcd.u.hcyang.top:31770

example02:
	go run example/concurrent/main.go --registry=etcd --registry_address=etcd.u.hcyang.top:31770

clear:
	rm mconfig

.PHONY: example