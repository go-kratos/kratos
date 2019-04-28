### tools/cache
##### Version 1.6.4
1. 修复某些参数下多make一次的问题

##### Version 1.6.3
1. 使用fanout替换cache包

##### Version 1.6.2
1. 改为使用errgroup提供的GOMAXPROCS方法 替换channel

##### Version 1.6.1
1. 弃用errgroup改用channel进行批量操作 防止线程饥饿
##### Version 1.6.0
1. 增加对metadata.WithContext的支持

##### Version 1.5.3
1. 优化gofmt提示

##### Version 1.5.2
1. 补充返回部分数据时的测试
2. 增加两种空缓存错误参数的检测
3. 支持// cache: 这样语法

##### Version 1.5.1
1. 批量模板中分批回源失败时候 返回部分数据

##### Version 1.5.0
1. 批量模板中改增加对数字类型0值返回的支持

##### Version 1.4.2
1. 修复回源失败 缓存数据未返回的问题

##### Version 1.4.1
1. 修复Hit计算问题
2. 由于mc已经有pkg/errors了 因此不再warp
3. 修复变量类型省略解析失败的问题

##### Version 1.4.0
1. 增加自定义注释和忽略参数的支持

##### Version 1.3.0
1. 增加batch_err选项 用于在分批发生错误的时候是否降级

##### Version 1.2.1
1. 回源错误的时候返回部分数据

##### Version 1.2.0
1. 解决saga提示无用代码的问题

##### Version 1.1.0
1. 去掉生成代码中的Cp前缀

##### Version 1.0.0

1. 添加基础模块与测试：
  - 代码生成组件
