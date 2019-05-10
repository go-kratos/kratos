
#### kratos-gen-mc

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

