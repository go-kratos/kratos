# 前提

服务注册与发现最简单的就是`direct`固定服务端地址的直连方式。也就是服务端正常监听端口启动不进行额外操作，客户端使用如下`target`：

```url
direct://default/127.0.0.1:9000,127.0.0.1:9091
```

> `target`就是标准的`URL`资源定位符[查看WIKI](https://zh.wikipedia.org/wiki/%E7%BB%9F%E4%B8%80%E8%B5%84%E6%BA%90%E5%AE%9A%E4%BD%8D%E7%AC%A6)

其中`direct`为协议类型，此处表示直接使用该`URL`内提供的地址`127.0.0.1:9000,127.0.0.1:9091`进行连接，而`default`在此处无意义仅当做占位符。

# gRPC Resolver

gRPC暴露了服务发现的接口`resolver.Builder`和`resolver.ClientConn`和`resolver.Resolver`，[官方代码位置](https://github.com/grpc/grpc-go/blob/master/resolver/resolver.go)：

```go
// Builder creates a resolver that will be used to watch name resolution updates.
type Builder interface {
	// Build creates a new resolver for the given target.
	//
	// gRPC dial calls Build synchronously, and fails if the returned error is
	// not nil.
	Build(target Target, cc ClientConn, opts BuildOption) (Resolver, error)
	// Scheme returns the scheme supported by this resolver.
	// Scheme is defined at https://github.com/grpc/grpc/blob/master/doc/naming.md.
	Scheme() string
}

// ClientConn contains the callbacks for resolver to notify any updates
// to the gRPC ClientConn.
//
// This interface is to be implemented by gRPC. Users should not need a
// brand new implementation of this interface. For the situations like
// testing, the new implementation should embed this interface. This allows
// gRPC to add new methods to this interface.
type ClientConn interface {
	// UpdateState updates the state of the ClientConn appropriately.
	UpdateState(State)
	// NewAddress is called by resolver to notify ClientConn a new list
	// of resolved addresses.
	// The address list should be the complete list of resolved addresses.
	//
	// Deprecated: Use UpdateState instead.
	NewAddress(addresses []Address)
	// NewServiceConfig is called by resolver to notify ClientConn a new
	// service config. The service config should be provided as a json string.
	//
	// Deprecated: Use UpdateState instead.
	NewServiceConfig(serviceConfig string)
}

// Resolver watches for the updates on the specified target.
// Updates include address updates and service config updates.
type Resolver interface {
	// ResolveNow will be called by gRPC to try to resolve the target name
	// again. It's just a hint, resolver can ignore this if it's not necessary.
	//
	// It could be called multiple times concurrently.
	ResolveNow(ResolveNowOption)
	// Close closes the resolver.
	Close()
}
```

下面依次分析这三个接口的作用：

* `Builder`用于gRPC内部创建`Resolver`接口的实现，但注意声明的`Build`方法将接口`ClientConn`作为参数传入了
* `ClientConn`接口有两个废弃方法不用管，看`UpdateState`方法需要传入`State`结构，看代码可以发现其中包含了`Addresses []Address // Resolved addresses for the target`，可以看出是需要将服务发现得到的`Address`对象列表告诉`ClientConn`的对象
* `Resolver`提供了`ResolveNow`用于被gRPC尝试重新进行服务发现

看完这三个接口就可以明白gRPC的服务发现实现逻辑，通过`Builder`进行`Reslover`的创建，在`Build`的过程中将服务发现的地址信息丢给`ClientConn`用于内部连接创建等逻辑。主要逻辑可以按下面顺序来看源码理解：

* 当`client`在`Dial`时会根据`target`解析的`scheme`获取对应的`Builder`，[官方代码位置](https://github.com/grpc/grpc-go/blob/master/clientconn.go#L242)
* 当`Dial`成功会创建出结构体`ClientConn`的对象[官方代码位置](https://github.com/grpc/grpc-go/blob/master/clientconn.go#L447)(注意不是上面的`ClientConn`接口)，可以看到结构体`ClientConn`内的成员`resolverWrapper`又实现了接口`ClientConn`的方法[官方代码位置](https://github.com/grpc/grpc-go/blob/master/resolver_conn_wrapper.go)
* 当`resolverWrapper`被初始化时就会调用`Build`方法[官方代码位置](https://github.com/grpc/grpc-go/blob/master/resolver_conn_wrapper.go#L89)，其中参数为接口`ClientConn`传入的是`ccResolverWrapper`
* 当用户基于`Builder`的实现进行`UpdateState`调用时，则会触发结构体`ClientConn`的`updateResolverState`方法[官方代码位置](https://github.com/grpc/grpc-go/blob/master/resolver_conn_wrapper.go#L109)，`updateResolverState`则会对传入的`Address`进行初始化等逻辑[官方代码位置](https://github.com/grpc/grpc-go/blob/master/clientconn.go#L553)

如此整个服务发现过程就结束了。从中也可以看出gRPC官方提供的三个接口还是很灵活的，但也正因为灵活要实现稍微麻烦一些，而`Address`[官方代码位置](https://github.com/grpc/grpc-go/blob/master/resolver/resolver.go#L79)如果直接被业务拿来用于服务节点信息的描述结构则显得有些过于简单。

所以`warden`包装了gRPC的整个服务发现实现逻辑，代码分别位于`pkg/naming/naming.go`和`warden/resolver/resolver.go`，其中：

* `naming.go`内定义了用于描述业务实例的`Instance`结构、用于服务注册的`Registry`接口、用于服务发现的`Resolver`接口
* `resolver.go`内实现了gRPC官方的`resolver.Builder`和`resolver.Resolver`接口，但也暴露了`naming.go`内的`naming.Builder`和`naming.Resolver`接口

# warden Resolver

接下来看`naming`内的接口如下：

```go
// Resolver resolve naming service
type Resolver interface {
	Fetch(context.Context) (*InstancesInfo, bool)
	Watch() <-chan struct{}
	Close() error
}

// Builder resolver builder.
type Builder interface {
	Build(id string) Resolver
	Scheme() string
}
```

可以看到封装方式与gRPC官方的方法一样，通过`Builder`进行`Resolver`的初始化。不同的是通过封装将参数进行了简化：

* `Build`只需要传对应的服务`id`即可：`warden/resolver/resolver.go`在gRPC进行调用后，会根据`Scheme`方法查询对应的`naming.Builder`实现并调用`Build`将`id`传入，而`naming.Resolver`的实现即可通过`id`去对应的服务发现中间件进行实例信息的查询
* 而`Resolver`则对方法进行了扩展，除了简单进行`Fetch`操作外还多了`Watch`方法，用于监听服务发现中间件的节点变化情况，从而能够实时的进行服务实例信息的更新

在`naming/discovery`内实现了基于[discovery](https://github.com/bilibili/discovery)为中间件的服务注册与发现逻辑。如果要实现其他中间件如`etcd`|`zookeeper`等的逻辑，参考`naming/discovery/discovery.go`内的逻辑，将与`discovery`的交互逻辑替换掉即可（后续会默认将etcd/zk等实现，敬请期待）。

# 使用discovery

因为`warden`内默认使用`direct`的方式，所以要使用[discovery](https://github.com/bilibili/discovery)需要在业务的`NewClient`前进行注册，代码如下：

```go
package dao

import (
	"context"

	"github.com/bilibili/kratos/pkg/naming/discovery"
	"github.com/bilibili/kratos/pkg/net/rpc/warden"
	"github.com/bilibili/kratos/pkg/net/rpc/warden/resolver"

	"google.golang.org/grpc"
)

// AppID your appid, ensure unique.
const AppID = "demo.service" // NOTE: example

func init(){
	// NOTE: 注意这段代码，表示要使用discovery进行服务发现
	// NOTE: 还需注意的是，resolver.Register是全局生效的，所以建议该代码放在进程初始化的时候执行
	// NOTE: ！！！切记不要在一个进程内进行多个不同中间件的Register！！！
	// NOTE: 在启动应用时，可以通过flag(-discovery.nodes) 或者 环境配置(DISCOVERY_NODES)指定discovery节点
	resolver.Register(discovery.Builder())
}

// NewClient new member grpc client
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (DemoClient, error) {
	client := warden.NewClient(cfg, opts...)
	conn, err := client.Dial(context.Background(), "discovery://default/"+AppID)
	if err != nil {
		return nil, err
	}
	// 注意替换这里：
	// NewDemoClient方法是在"api"目录下代码生成的
	// 对应proto文件内自定义的service名字，请使用正确方法名替换
	return NewDemoClient(conn), nil
}
```

> 注意：`resolver.Register`是全局行为，建议放在包加载阶段或main方法开始时执行，该方法执行后会在gRPC内注册构造方法

`target`是`discovery://default/${appid}`，当gRPC内进行解析后会得到`scheme`=`discovery`和`appid`，然后进行以下逻辑：

1. `warden/resolver.Builder`会通过`scheme`获取到`naming/discovery.Builder`对象（靠`resolver.Register`注册过的）
2. 拿到`naming/discovery.Builder`后执行`Build(appid)`构造`naming/discovery.Discovery`
3. `naming/discovery.Discovery`对象基于`appid`就知道要获取哪个服务的实例信息

# 服务注册

客户端既然使用了[discovery](https://github.com/bilibili/discovery)进行服务发现，也就意味着服务端启动后必须将自己注册给[discovery](https://github.com/bilibili/discovery)知道。

相对服务发现来讲，服务注册则简单很多，看`naming/discovery/discovery.go`内的代码实现了`naming/naming.go`内的`Registry`接口，服务端启动时可以参考下面代码进行注册：

```go
// 该代码可放在main.go，当warden server进行初始化之后
// 省略...

ip := "" // NOTE: 必须拿到您实例节点的真实IP，
port := "" // NOTE: 必须拿到您实例grpc监听的真实端口，warden默认监听9000
hn, _ := os.Hostname()
dis := discovery.New(nil)
ins := &naming.Instance{
    Zone:     env.Zone,
    Env:      env.DeployEnv,
    AppID:    "your app id",
    Hostname: hn,
    Addrs: []string{
        "grpc://" + ip + ":" + port,
    },
}
cancel, err := dis.Register(context.Background(), ins)
if err != nil {
    panic(err)
}

// 省略...

// 特别注意！！！
// cancel必须在进程退出时执行！！！
cancel()
```



# 使用ETCD

和使用discovery类似,只需要在注册时使用etcd naming即可。

```go
package dao

import (
	"context"

	"github.com/bilibili/kratos/pkg/naming/etcd"
	"github.com/bilibili/kratos/pkg/net/rpc/warden"
	"github.com/bilibili/kratos/pkg/net/rpc/warden/resolver"

	"google.golang.org/grpc"
)

// AppID your appid, ensure unique.
const AppID = "demo.service" // NOTE: example

func init(){
	// NOTE: 注意这段代码，表示要使用etcd进行服务发现 ,其他事项参考discovery的说明
    // NOTE: 在启动应用时，可以通过flag(-etcd.endpoints) 或者 环境配置(ETCD_ENDPOINTS)指定etcd节点
    // NOTE: 如果需要自己指定配置时 需要同时设置DialTimeout 与 DialOptions: []grpc.DialOption{grpc.WithBlock()}
	resolver.Register(etcd.Builder(nil))
}

// NewClient new member grpc client
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (DemoClient, error) {
	client := warden.NewClient(cfg, opts...)
  	// 这里使用etcd scheme
	conn, err := client.Dial(context.Background(), "etcd://default/"+AppID)
	if err != nil {
		return nil, err
	}
	// 注意替换这里：
	// NewDemoClient方法是在"api"目录下代码生成的
	// 对应proto文件内自定义的service名字，请使用正确方法名替换
	return NewDemoClient(conn), nil
}
```

etcd的服务注册与discovery基本相同,可以传入详细的etcd配置项, 或者传入nil后通过flag(-etcd.endpoints)/环境配置(ETCD_ENDPOINTS)来指定etcd节点。

### 其他配置项

etcd默认的全局keyPrefix为kratos_etcd,当该keyPrefix与项目中其他keyPrefix冲突时可以通过flag(-etcd.prefix)或者环境配置(ETCD_PREFIX)来指定keyPrefix。



# 扩展阅读

[warden快速开始](warden-quickstart.md) [warden拦截器](warden-mid.md) [warden基于pb生成](warden-pb.md) [warden负载均衡](warden-balancer.md)

-------------

[文档目录树](summary.md)
