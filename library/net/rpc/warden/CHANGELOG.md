### net/rpc/warden

##### Version 1.1.9
1. 增加NonBlock模式

##### Version 1.1.8
1. 新增appid mock

##### Version 1.1.7
1. 兼容cpu为0和wrr dt为0的情况

##### Version 1.1.6
1. 修改caller传递和获取方式
2. 添加error detail example

##### Version 1.1.5
1. 增加server端json格式支持

##### Version 1.1.4
1. 判断reosvler.builder为nil之后再注册

##### Version 1.1.3
1. 支持zone和clusters

##### Version 1.1.2
1. 业务错误日志记为 WARN

##### Version 1.1.1
1. server实现了返回cpu信息

##### Version 1.1.0
1. 增加ErrorDetail
2. 修复日志打印error信息丢失问题

##### Version 1.0.3
1. 给server增加keepalive参数

##### Version 1.0.2

1. 替代默认的timoue，使用durtaion.Shrink()来传递context
2. 修复peer.Addr为nil时会panic的问题

##### Version 1.0.1

1. 去除timeout的手动传递，改为使用grpc默认自带的grpc-timeout
2. 获取server address改为使用call option的方式，去除对balancer的依赖

##### Version 1.0.0

1. 使用NewClient来新建一个RPC客户端，并默认集成trace、log、recovery、moniter拦截器
2. 使用NewServer来新建一个RPC服务端，并默认集成trace、log、recovery、moniter拦截器
