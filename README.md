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

[TOC]



### Feature

- [x] 配置数据value支持复杂JSON，并且此JSON需要提前Schema定义 
- [x] 灰度发布配置
- [x] 基本功能达到

### config example

```json
{
    "1000-100": {
        "ABFilters": {
            "ip": "192.168.1.12"
        },
        "configs": {
            "entry": {
                "0": {
                    "config": "{'name':'demo13','age':12}",
                    "create_time": 1604249335,
                    "schema": "{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}",
                    "update_time": 1604249335
                },
                "1": {
                    "config": "{'name':'demo1','age':24}",
                    "create_time": 1604249335,
                    "schema": "{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}",
                    "update_time": 1604249335
                },
                "2": {
                    "config": "{'name':'demo1','age':23}",
                    "create_time": 1604249335,
                    "schema": "{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}",
                    "update_time": 1604249335
                }
            }
        },
        "create_time": 1604249335,
        "desc": "test",
        "update_time": 1604249335
    },
    "1000-101": {
        "ABFilters": {
            "ip": "192.0.0.1"
        },
        "configs": {
            "entry": {
                "0": {
                    "config": "{'name':'demo1','age':12}",
                    "create_time": 1604249335,
                    "schema": "{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}",
                    "update_time": 1604249335
                },
                "1": {
                    "config": "{'name':'demo1','age':12}",
                    "create_time": 1604249335,
                    "schema": "{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}",
                    "update_time": 1604249335
                },
                "2": {
                    "config": "{'name':'demo1','age':12}",
                    "create_time": 1604249335,
                    "schema": "{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}",
                    "update_time": 1604249335
                }
            }
        },
        "create_time": 1604249335,
        "desc": "test",
        "update_time": 1604249335
    },
    "1000-102": {
        "ABFilters": {
            "ip": "192.0.0.1"
        },
        "configs": {
            "entry": {
                "0": {
                    "config": "{'name':'demo1','age':12}",
                    "create_time": 1604249335,
                    "schema": "{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}",
                    "update_time": 1604249335
                },
                "1": {
                    "config": "{'name':'demo1','age':13}",
                    "create_time": 1604249335,
                    "schema": "{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}",
                    "update_time": 1604249335
                },
                "2": {
                    "config": "{'name':'demo1','age':12}",
                    "create_time": 1604249335,
                    "schema": "{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}",
                    "update_time": 1604249335
                }
            }
        },
        "create_time": 1604249335,
        "desc": "test",
        "update_time": 1604249335
    },
    "1000-103": {
        "ABFilters": {
            "ip": "192.0.0.1"
        },
        "configs": {
            "entry": {
                "0": {
                    "config": "{'name':'demo1','age':12}",
                    "create_time": 1604249335,
                    "schema": "{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}",
                    "update_time": 1604249335
                },
                "1": {
                    "config": "{'name':'demo1','age':22}",
                    "create_time": 1604249335,
                    "schema": "{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}",
                    "update_time": 1604249335
                },
                "2": {
                    "config": "{'name':'demo1','age':12}",
                    "create_time": 1604249335,
                    "schema": "{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}",
                    "update_time": 1604249335
                }
            }
        },
        "create_time": 1604249335,
        "desc": "test",
        "update_time": 1604249335
    },
    "1000-104": {
        "ABFilters": {
            "ip": "192.0.0.1"
        },
        "configs": {
            "entry": {
                "0": {
                    "config": "{'name':'demo1','age':12}",
                    "create_time": 1604249335,
                    "schema": "{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}",
                    "update_time": 1604249335
                },
                "1": {
                    "config": "{'name':'demo1','age':12}",
                    "create_time": 1604249335,
                    "schema": "{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}",
                    "update_time": 1604249335
                },
                "2": {
                    "config": "{'name':'demo1','age':12}",
                    "create_time": 1604249335,
                    "schema": "{'type': 'object','properties':{'name':{'type':'string'},'age':{'type':'integer'}}}",
                    "update_time": 1604249335
                }
            }
        },
        "create_time": 1604249335,
        "desc": "test",
        "update_time": 1604249335
    }
}
```
