# go-common/sync/pipeline/fanout

以前的library/cache包改为pipeline/fanout

增加使用范围 不止由于异步增加缓存 也可以用在其他地方

功能:

* 支持定义Worker 数量的goroutine，进行消费
* 内部支持的元数据传递（library/net/metadata）
* 后续会作废library/cache，以及统一收敛Go并行里面的扇出模型

示例:
```golang
//名称为cache 执行线程为1 buffer长度为1024
cache := fanout.New("cache", fanout.Worker(1), fanout.Buffer(1024))
cache.Do(c, func(c context.Context) { SomeFunc(c, args...) })
cache.Close()
```