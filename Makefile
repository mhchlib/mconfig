
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -o mconfig-server cmd/mconfig-server/mconfig-server.go cmd/mconfig-server/plugin.go

docker: build
	docker build -t dockerhcy/mconfig:v0.2  .

push: docker
	docker push dockerhcy/mconfig:v0.2

dev:
	go run cmd/mconfig-server/mconfig-server.go cmd/mconfig-server/plugin.go cmd/mconfig-server/debug.go  --namespace=local_test --registry=etcd://etcd.u.hcyang.top:31770 --store=etcd://etcd.u.hcyang.top:31770 --expose :8081 --debug

dev2:
	go run cmd/mconfig-server/mconfig-server.go cmd/mconfig-server/plugin.go --namespace=local_test --registry=etcd://etcd.u.hcyang.top:31770 --store=file://ttt

clear:
	rm mconfig

.PHONY: example