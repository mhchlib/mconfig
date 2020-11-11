## MConfig

```
				   action UPDATE
                                   / 
                                  /
                                 /
    app -> prefix + appid -> JSON ---- action ADD
                               | \
                               |  \
                               |   \
                            action  action
                            SELECT  DELETE
```

>客户端携带appid(暂时方案)以grpc服务端流的形式访问mconfig服务，mconfig为本次连接建立一个chan通道，并阻塞读取，考虑到网络抖动，此chan具有
 缓存功能缓存大小为5(暂定), mconfig-admin 服务为管理界面 并管理着etcd中的配置内容，mconfig服务会监听每一个连接他的服务需要的appid配置的
 key变化,收到变化后,diff mconfig本地的配置 并把修改内容推送给连接的服务

### Feature

 * 配置数据value支持复杂JSON，并且此JSON需要提前schema定义
 
 
#### config example

```go
type ConfigEntity struct {
	Id         string
	Schema     string
	Config     string
	Status     common.ConfigStatus
	Desc       string
	CreateTime int64
	UpdateTime int64
}
```
```json
{
    "100":[
        {"id":"1000","config":"{'name':'demo1','age':12}","schema":"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}","create_time":1604249335,"update_time":1604249335,"desc":"test","status":0},
        {"id":"1001","config":"{'name':'demo2','age':13}","schema":"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}","create_time":1604249335,"update_time":1604249335,"desc":"test","status":0},
        {"id":"1002","config":"{'name':'demo3','age':14}","schema":"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}","create_time":1604249335,"update_time":1604249335,"desc":"test","status":0},
        {"id":"1003","config":"{'name':'demo4','age':15}","schema":"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}","create_time":1604249335,"update_time":1604249335,"desc":"test","status":0},
        {"id":"1004","config":"{'name':'demo5','age':16}","schema":"{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}","create_time":1604249335,"update_time":1604249335,"desc":"test","status":0}
    ]
}
```
