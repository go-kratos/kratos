#### breaker

##### 项目简介
> 提供熔断器功能，供各种client（如rpc、http、msyql）等进行熔断
> 提供Go方法供业务在breaker熔断前后进行回调处理

##### 编译环境
> 1. 请只用golang v1.8.x以上版本编译执行。

##### 依赖包
> 1. 无

##### 配置说明
> 1. NewGroup(name string,c *Config)当c==nil时则采用默认配置
> 2. 可通过breaker.Init(c *Config)替换默认配置
> 3. 可通过group.Reload(c *Config)进行配置更新
> 4. 默认配置如下所示：
    _conf = &Config{
            Window:  xtime.Duration(3 * time.Second),
            Sleep:   xtime.Duration(100 * time.Millisecond),
            Bucket:  10,
            Ratio:   0.5,
            Request: 100,
    }

##### 测试
> 1. 执行当前目录下所有测试文件，测试所有功能
