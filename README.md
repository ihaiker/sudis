# distributed gosuv

在软件[gosuv](https://github.com/codeskyblue/gosuv)实现分布式管理。更多详细介绍请查阅[gosuv](https://github.com/codeskyblue/gosuv)。

![gosuv web](docs/gosuv.gif)

##分布式配置方式
在config.yaml的server节点下面配置主节点的地址 `master`
```
server:
  httpauth:
    enabled: false
    username:
    password:
  addr: :11313
  master: 172.16.10.100:11313
client:
  server_url: http://localhost:11313
```
注意：急群众所有应用的httpauth配置要一致。并且所有操作只能在主节点完成。

