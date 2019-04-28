### memcache client

##### Version 1.5.1
> 1.修复parse reply时如果有err不关闭连接问题  

##### Version 1.5.0
> 1.支持cache和cache conn的写法
##### Version 1.4.0
> 1.add memcache mock conn

##### Version 1.3.2
> 1.修复判断是否合法key

##### Version 1.3.1
> 1.修复pool放回连接的bug

##### Version 1.3.0
> 1.修改memcache pool的实现方式，引用container/pool   
> 2.pool支持context传入超时以及Get connection WaitTimeout

##### Version 1.2.0
> 1. 增加pkg errors

##### Version 1.1.2
> 1. 修复gzip writer默认压缩level为0的bug

##### Version 1.1.1
> 1. fix populateOne error

##### Version 1.1.0
> 1. memcache添加largevalue支持

##### Version 1.0.0
> 1. 修改decode时protobuf bug,补全测试
