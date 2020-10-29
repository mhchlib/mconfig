## MConfig

```
                                   / action UPDATE
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