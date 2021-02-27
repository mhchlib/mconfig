### START
 
#### INSTALL

1 INSTALL MCONFIG SERVER

```shell
    docker pull dockerhcy/mconfig-server
    docker run -it dockerhcy/mconfig-server:v1.0.0-70-g4cb575f  --registry=etcd://etcd.u.hcyang.top:31770 --store=etcd://etcd.u.hcyang.top:31770 --expose :8081 --debug
```

2 INSTALL MCONFIG ADMIN

```shell
    docker pull dockerhcy/mconfig-admin
```
