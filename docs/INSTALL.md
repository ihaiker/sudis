# 安装说明

> 中文都没有搞转，就不写英文文档了，害人害己，哈哈😄。

## 依赖程序

- go 1.12+ [安装教程](https://www.runoob.com/go/go-environment.html)
- nodejs (npm,vue) [安装教程](https://www.runoob.com/nodejs/nodejs-install-setup.html)
- go-bindata [教程](https://github.com/shuLhan/go-bindata)
- make (非必须)

注：如果不明白如何安装的同学，请自行Google或百度。

## 安装

### 二进制安装

[下载二进制包](https://github.com/ihaiker/sudis/releases)

### 源码安装

#### 下载源代码

```shell script
git clone https://github.com/ihaiker/sudis.git
```

#### 编译程序（make方式）

```shell script
cd sudis
make release
```
编译完的程序在 当前文件夹bin目录下，编译完成。

#### 编译程序
第一步：编译生成前端页面
```shell script
$ cd sudis/webui 
webui$ npm i -g @vue/cli #安装vue cli
webui$ npm i #安装依赖
webui$ npm run build
```
> npm安装会很慢，安装[cnpm](https://npm.taobao.org/)会很快。

第二步：使用 go-bindata 把页面文件打包。

```shell script
$ go generate generator.go #执行此步需要安装go-bindata
```

第三步：下载go依赖包
```shell script
$ go mod download
```

第四步：编译
```shell script
$ go build 
```

### 程序配置
复制`conf/sudis.toml.example` 到 `bin/conf/sudis.toml`
```toml
[master]
  band = ":5983"
  http = ":5984"
  securityToken = "4E4AD35C6C0BEB20DC343A1E8F7E32D4"
  salt = "2CCAKYGBPTCET2S6"
  [master.database]
    type = "sqlite3"
    url = "/etc/sudis/sudis.db"

[server]
  dir = "/etc/sudis/programs"
  sock = "unix://etc/sudis/sudis.sock"
  master = "tcp://127.0.0.1:5983"
  securityToken = "4E4AD35C6C0BEB20DC343A1E8F7E32D4"
```
> 配置解释

`master`: 用来配置主控节点信息

`master.http`: 主控节点开放HTTP服务地址，默认：:5984

`master.salt`：管理端采用无状态控制用户登录，此值为用生成无状态验证串添加盐值。**（务必修改）**

`master.bind`：主控节点绑定的TCP端口地址，用于分布情况下程序server节点加入主控。默认：:5983

`master.securityToken`：master,server节点通信认证的安全串。

`master.database`：主控节点数据库配置。支持: sqlite3,mysql

```toml
[master]
  [master.database]
    type = "mysql"
    url = "sudis:passwd@127.0.0.1:3306/sudis?charset=utf8"
```

`server` : 程序控制节点

`server.dir`: 程序配置文件所在位置。默认：$PWD/conf/programs

`server.sock`: 节点控制sock服务连接地址。默认 $PWD/conf/sudis.sock

`server.master`: 连接主控节点地址。

`server.securityToken`：master,server节点通信认证的安全串。



## 运行程序

### 初始化中控节点：

```shell
$ sudis master init 
```

执行完成此步，会在数据库中建立相应的数据表结构，初始化管理用户。



### **启动中控节点(master)：**

```shell
$ sudis master
```



### **程序控制节点(server)：**

```shell
$ sudis server
```

程序控制节点可以分布在多台机器上，配置方式和启动完全一致。



###  单节点启动

> 上面的程序启动是分布式情况下的启动，如果您只是单机是使用可以使用独立模式运行

```shell
$ sudis single
或者
$ sudis # single 是默认命令
```


### 登录中控台

打开地址 http://master:5984 即可。`master`:为主控节点IP， 默认的登录用户为:admin，密码：12345678


### 开机启动项：

编译完程序后执行`autostart.sh`即可添加到开机启动脚本。你可以根据提示选择开机启动的节点

## 更多命令

程序启动参数和命令可以通过 -h 帮助方式查询例如：

```shell
$ ./bin/sudis -h
sudis, 一个分布式进程控制程序。

Usage:
  sudis [flags]
  sudis [command]

Available Commands:
  add         添加程序管理
  console     管理端命令
  delete      删除管理的程序
  detail      查看配置信息，JSON
  help        Help about any command
  list        查看程序列表
  master      管理控制端
  modify      修改程序
  server      守护进程管理端
  shutdown    关闭进程管理服务
  single      独立模式启动(默认命令)
  start       启动管理的程序
  status      查看运行状态
  stop        停止管理的程序
  tail        查看日志

Flags:
  -f, --conf string    配置文件
  -d, --debug          Debug模式
  -h, --help           help for sudis
  -l, --level string   日志级别 (default "info")
      --version        version for sudis
```

