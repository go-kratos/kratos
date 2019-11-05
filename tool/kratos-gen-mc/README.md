
#### genmc

> mc缓存代码生成

##### 项目简介

自动生成memcached缓存代码 和缓存回源工具kratos-gen-bts配合使用 体验更佳
支持以下功能:
- 常用mc命令(get/set/add/replace/delete)
- 多种数据存储格式(json/pb/raw/gob/gzip)
- 常用值类型自动转换(int/bool/float...)
- 自定义缓存名称和过期时间
- 记录pkg/error错误栈
- 记录日志trace id
- prometheus错误监控
- 自定义参数个数
- 自定义注释

##### 使用方式:
1. dao.go文件中新增 _mc interface
2. 在dao 文件夹中执行 go generate命令 将会生成相应的缓存代码
3. 示例见testdata/dao.go

##### 注意:
类型会根据前缀进行猜测
set / add 对应mc方法Set
replace 对应mc方法 Replace
del 对应mc方法 Delete
get / cache对应mc方法Get
mc Add方法需要用注解 -type=only_add单独指定

#### 注解参数:
| 名称        | 默认值              | 可用范围         | 说明                                                         | 可选值                       | 示例                       |
| ----------- | ------------------- | ---------------- | ------------------------------------------------------------ | ---------------------------- | -------------------------- |
| encode      | 根据值类型raw或json | set/add/replace  | 数据存储的格式                                               | json/pb/raw/gob/gzip         | json 或 json\|gzip 或gob等 |
| type        | 前缀推断            | 全部             | mc方法 set/get/delete...                                     | get/set/del/replace/only_add | get 或 replace 等          |
| key         | 根据方法名称生成    | 全部             | 缓存key名称                                                  | -                            | demoKey                 |
| expire      | 根据方法名称生成    | 全部             | 缓存过期时间                                                 | -                            | d.demoExpire            |
| batch       |                     | get(限多key模板) | 批量获取数据 每组大小                                        | -                            | 100                        |
| max_group   |                     | get(限多key模板) | 批量获取数据 最大组数量                                      | -                            | 10                         |
| batch_err   | break               | get(限多key模板) | 批量获取数据回源错误的时候 降级继续请求(continue)还是直接返回(break) | break 或 continue            | continue                   |
| struct_name | dao                 | 全部             | 用户自定义Dao结构体名称                                      |                              | MemcacheDao                |
|check_null_code||add/set|(和null_expire配套使用)判断是否是空缓存的代码 用于为空缓存独立设定过期时间||$.ID==-1 或者 $=="-1"等|
|null_expire|300(5分钟)|add/set|(和check_null_code配套使用)空缓存的过期时间||d.nullExpire|