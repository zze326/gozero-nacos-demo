# 调用逻辑

demo 提供一个 http api 服务，其包含一个 `get /user/:id` 的接口，该接口会通过 rpc 调用 user 服务获取用户信息返回。

# 功能点

## 一行代码完成 nacos 配置加载

经过对 nacos sdk 的多层封装，此时 gozero 服务要获取 nacos 配置及其简单，以 `demo/api/demo.go` 为例，其只需要下面一行就可以完成 nacos 配置的加载：

```go
ctx := svc.NewServiceContext(c, commonConfig.MustLoad(*configFile, c))
```

当然，加载 nacos 配置前 nacos 本身的配置也是必不可少的，所以原来 `etc` 文件夹下现在仅有一个 `nacos.yaml` 配置：

```yaml
Addr: nacos-cs.ops.svc.cluster.local
Port: 8848
Group: DEFAULT_GROUP
DataID: demo-api
ExtDataIDs:
  - mysql
  - auth
NamespaceID: demo
```

没什么可说的，上面几个字段都和 nacos 的基本概念对应。
要补充的是：

- 该配置支持同时配置多个 `data-id`，`ExtDataIDs` 的优先级比 `DataID` 低，即 `DataID` 对应的配置内容如果和 `ExtDataIDs` 列表包含的配置内容有冲突的话，会优先使用 `DataID`
  中的配置；
- `ExtDataIDs` 是一个列表，其加载顺序是从上到下，下面的配置会覆盖前面的配置；

## 简单 RPC 服务的注册

既然配置中心已经使用了 nacos，注册中心再使用另一个中间件就没有必要了。

为了便于后续项目的使用，我对注册逻辑也做了简单的封装。以 `user/rpc/user.go` 为例，要注册其到 nacos 中，仅需要下面一行：

```go
commonConfig.MustRegister(commonConfig.MustLoad(*configFile, c), &c.RpcServerConf)
```

## 快捷的 RPC 调用

既然 user 服务已经注册到 nacos 了，那么其它服务如何调用它呢，直接看 `demo/api/internal/svc/serviceContext.go` 中的代码：

```go
UserRpc: user.NewUser(nc.NewZrpcClient("user.rpc", c.Name)),
```

> 这里的 `"user.rpc"` 为 user 服务注册到 nacos 中的项目名，该名称在 nacos 中进行配置，而 `c.Name` 其实也是取自 nacos 中的配置，即为当前服务的名称，如下图：
> [服务列表](https://raw.githubusercontent.com/zze326/gozero-nacos-demo/master/imgs/service-list.png)
> [服务消费者](https://raw.githubusercontent.com/zze326/gozero-nacos-demo/master/imgs/consumer.png)

也是只需要一行，就可以实例化一个 user 服务的 rpc 远程调用对象，将其置入 `ServiceContext` 中就可以在后续的 `logic` 中使用它完成 rpc 调用了。
如 `demo/api/internal/logic/getUserLogic.go` 的第 28 行：

```go
u, err := l.svcCtx.UserRpc.GetUser(l.ctx, &user.GetUserReq{
Id: uint32(req.ID),
})
```

# 补充点

## 配置内容

- 配置列表：
  [data id 列表](https://raw.githubusercontent.com/zze326/gozero-nacos-demo/master/imgs/config-list.png)

- demo-api：
  [demo-api](https://raw.githubusercontent.com/zze326/gozero-nacos-demo/master/imgs/demo-api.png)

- user-rpc
  [user-rpc](https://raw.githubusercontent.com/zze326/gozero-nacos-demo/master/imgs/user-rpc.png)

- auth
  [auth](https://raw.githubusercontent.com/zze326/gozero-nacos-demo/master/imgs/auth.png)

- mysql
  [mysql](https://raw.githubusercontent.com/zze326/gozero-nacos-demo/master/imgs/mysql.png)

## nacos sdk 版本问题

有一个问题是 gozero 社区仓库中使用 nacos 做注册中心的 demo 引用的 sdk 是适配 nacos 1.x 的，地址：

- <https://github.com/zeromicro/zero-contrib/tree/main/zrpc/registry/nacos>

而我使用的 nacos 版本为 2.1.0，所以用其接入 nacos 2.x 的使用过程中总有一些莫名其妙的问题，所以最后我 fork 了它并自行修改了引用的 sdk 支持 nacos 2.x，地址：

- <https://github.com/zze326/zero-contrib/tree/main/zrpc/registry/nacos>

所以在此项目的 `go.mod` 中可以看到我引用的库为：

```mod
github.com/zze326/zero-contrib/zrpc/registry/nacos v0.0.0-20220526111920-4c5f0ff42470
```
