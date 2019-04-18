#### mcgen

> 缓存代码生成

##### 项目简介
本工具用于自动生成memcached缓存代码.

### 支持以下功能

1. 常用mc命令(get/set/add/replace/delete)
2. 多种数据存储格式(json/pb/raw/gob/gzip)
3. 自定义缓存名称和过期时间
4. 常用值类型自动转换(int/bool/float...)
5. 分批获取数据 
6. 记录pkg/error错误栈
7. 记录日志trace id
8. prometheus错误监控
9. 自定义参数个数
10. 自定义注释

### 使用方式举例
dao.go文件中新增: (缓存声明不要删除 保留在源码里)

### 要求:
需要安装cachegen

```go
//go:generate kratos tool mcgen
type _mc interface {
   // mc: -key=articleKey
   CacheArticles(c context.Context, keys []int64) (map[int64]*Article, error)
   // mc: -key=articleKey
   CacheArticle(c context.Context, key int64) (*Article, error)
   // mc: -key=keyMid
   CacheArticle1(c context.Context, key int64, mid int64) (*Article, error)
   // mc: -key=noneKey
   CacheNone(c context.Context) (*Article, error)

   // mc: -key=articleKey -expire=d.articleExpire -encode=json
   AddCacheArticles(c context.Context, values map[int64]*Article) error
   // 这里也支持自定义注释 会替换默认的注释
   // mc: -key=articleKey -expire=d.articleExpire -encode=json|gzip
   AddCacheArticle(c context.Context, key int64, value *Article) error
   // mc: -key=keyMid -expire=d.articleExpire -encode=gob
   AddCacheArticle1(c context.Context, key int64, value *Article, mid int64) error
   // mc: -key=noneKey
   AddCacheNone(c context.Context, value *Article) error

   // mc: -key=articleKey
   DelCacheArticles(c context.Context, keys []int64) error
   // mc: -key=articleKey
   DelCacheArticle(c context.Context, key int64) error
   // mc: -key=keyMid
   DelCacheArticle1(c context.Context, key int64, mid int64) error
   // mc: -key=noneKey
   DelCacheNone(c context.Context) error
}
```
执行go generate就会生成相应的代码了.

## MC方法类型

类型会根据前缀进行猜测

set / add 对应mc方法Set

replace 对应mc方法 Replace

del 对应mc方法 Delete

get / cache对应mc方法Get

mc Add方法需要用注解 -type=only_add单独指定

### 注解参数

格式: -key=value

| 名称        | 默认值              | 可用范围         | 说明                                                         | 可选值                       | 示例                       |
| ----------- | ------------------- | ---------------- | ------------------------------------------------------------ | ---------------------------- | -------------------------- |
| encode      | 根据值类型raw或json | set/add/replace  | 数据存储的格式                                               | json/pb/raw/gob/gzip         | json 或 json\|gzip 或gob等 |
| type        | 前缀推断            | 全部             | mc方法 set/get/delete...                                     | get/set/del/replace/only_add | get 或 replace 等          |
| key         | 根据方法名称生成    | 全部             | 缓存key名称                                                  | -                            | articleKey                 |
| expire      | 根据方法名称生成    | 全部             | 缓存过期时间                                                 | -                            | d.articleExpire            |
| batch       |                     | get(限多key模板) | 批量获取数据 每组大小                                        | -                            | 100                        |
| max_group   |                     | get(限多key模板) | 批量获取数据 最大组数量                                      | -                            | 10                         |
| batch_err   | break               | get(限多key模板) | 批量获取数据回源错误的时候 降级继续请求(continue)还是直接返回(break) | break 或 continue            | continue                   |
| struct_name | Dao                 | 全部             | 用户自定义Dao结构体名称                                      |                              | MemcacheDao                |

Q&A
  Q: 为什么add 不叫add叫only_add? 
  A: mc中的add方法指的是当key不存在的时候才新加 绝大多数人其实需要的是set方法 为了有所区别 防止弄混 所以换了一个名称
