VERSION=$(shell git describe --tags --always --dirty --dirty="")

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -o mconfig-server cmd/mconfig-server/mconfig-server.go cmd/mconfig-server/plugin.go

image: build
	docker build -t dockerhcy/mconfig:${VERSION}   .

push: image
	docker push dockerhcy/mconfig:${VERSION}

dev:
	go run cmd/mconfig-server/mconfig-server.go cmd/mconfig-server/plugin.go cmd/mconfig-server/debug.go  --namespace=local_test --registry=etcd://etcd.u.hcyang.top:31770 --store=etcd://etcd.u.hcyang.top:31770 --expose :8081 --debug

clean:
	-rm mconfig-server

.PHONY: clean