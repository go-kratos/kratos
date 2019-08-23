### kratos tool genbts

> 缓存回源代码生成

在internal/dao/dao.go中添加mc缓存interface定义，可以指定对应的[注解参数](../../tool/kratos-gen-bts/README.md)；  
并且在接口前面添加`go:generate kratos tool genbts`；  
然后在当前目录执行`go generate`，可以看到自动生成的dao.bts.go代码。  

### 回源模板
```go
//go:generate kratos tool genbts
type _bts interface {
	// bts: -batch=2 -max_group=20 -batch_err=break -nullcache=&Demo{ID:-1} -check_null_code=$.ID==-1
	Demos(c context.Context, keys []int64) (map[int64]*Demo, error)
	// bts: -sync=true -nullcache=&Demo{ID:-1} -check_null_code=$.ID==-1
	Demo(c context.Context, key int64) (*Demo, error)
	// bts: -paging=true
	Demo1(c context.Context, key int64, pn int, ps int) (*Demo, error)
	// bts: -nullcache=&Demo{ID:-1} -check_null_code=$.ID==-1
	None(c context.Context) (*Demo, error)
}
```

### 参考

也可以参考完整的testdata例子：kratos/tool/kratos-gen-bts/testdata

-------------

[文档目录树](summary.md)
