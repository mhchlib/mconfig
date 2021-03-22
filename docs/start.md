### START

1 install mconfig server
```shell
docker pull dockerhcy/mconfig-server
```

2 run mconfig server
```shell
docker run -it dockerhcy/mconfig-server:v1.0.0-70-g4cb575f \
  --registry=etcd://etcd.u.hcyang.top:31770 \
  --store=etcd://etcd.u.hcyang.top:31770 \
  --expose :8081 \
  --debug
```

3 install mconfig admin
```shell
docker pull dockerhcy/mconfig-admin
```

4 run mconfig admin


5 open mconfig admin website on the browser

