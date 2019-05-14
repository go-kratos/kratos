### kratos tool genmc

> 缓存代码生成

在internal/dao/dao.go中添加mc缓存interface定义，可以指定对应的[注解参数](../../tool/kratos-gen-mc/README.md)；  
并且在接口前面添加`go:generate kratos tool genmc`；  
然后在当前目录执行`go generate`，可以看到自动生成的mc.cache.go代码。  

### 缓存模板
```go
//go:generate kratos tool genmc
type _mc interface {
	// mc: -key=demoKey
	CacheDemos(c context.Context, keys []int64) (map[int64]*Demo, error)
	// mc: -key=demoKey
	CacheDemo(c context.Context, key int64) (*Demo, error)
	// mc: -key=keyMid
	CacheDemo1(c context.Context, key int64, mid int64) (*Demo, error)
	// mc: -key=noneKey
	CacheNone(c context.Context) (*Demo, error)
	// mc: -key=demoKey
	CacheString(c context.Context, key int64) (string, error)

	// mc: -key=demoKey -expire=d.demoExpire -encode=json
	AddCacheDemos(c context.Context, values map[int64]*Demo) error
	// mc: -key=demo2Key -expire=d.demoExpire -encode=json
	AddCacheDemos2(c context.Context, values map[int64]*Demo, tp int64) error
	// 这里也支持自定义注释 会替换默认的注释
	// mc: -key=demoKey -expire=d.demoExpire -encode=json|gzip
	AddCacheDemo(c context.Context, key int64, value *Demo) error
	// mc: -key=keyMid -expire=d.demoExpire -encode=gob
	AddCacheDemo1(c context.Context, key int64, value *Demo, mid int64) error
	// mc: -key=noneKey
	AddCacheNone(c context.Context, value *Demo) error
	// mc: -key=demoKey -expire=d.demoExpire
	AddCacheString(c context.Context, key int64, value string) error

	// mc: -key=demoKey
	DelCacheDemos(c context.Context, keys []int64) error
	// mc: -key=demoKey
	DelCacheDemo(c context.Context, key int64) error
	// mc: -key=keyMid
	DelCacheDemo1(c context.Context, key int64, mid int64) error
	// mc: -key=noneKey
	DelCacheNone(c context.Context) error
}

func demoKey(id int64) string {
	return fmt.Sprintf("art_%d", id)
}

func demo2Key(id, tp int64) string {
	return fmt.Sprintf("art_%d_%d", id, tp)
}

func keyMid(id, mid int64) string {
	return fmt.Sprintf("art_%d_%d", id, mid)
}

func noneKey() string {
	return "none"
}
```

### 参考

也可以参考完整的testdata例子：kratos/tool/kratos-gen-mc/testdata

-------------

[文档目录树](summary.md)
