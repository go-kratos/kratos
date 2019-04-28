## .rider

###Version 2.20
1. 去除megacheck, unused, gosimple三项

###Version 2.19
1. ut开放library目录执行

###Version 2.18
1. 更新unit_test.sh的ReadDir()逻辑

###Version 2.17
1. sagacheck 添加bgr工具作为规则引擎

###Version 2.16
1. bazel coverage --config=ci 

###Version 2.15
1. ut增加自定义lint检查
2. ut 覆盖率文件过滤monkey.go文件
3. ut magic流程增加评论用户校验

###Version 2.14
1. ut执行增加 convey json.
2. ut执行增加过滤器

###Version 2.13
1. ut脚本增加覆盖率原始文件上传

### Version 2.12
1. ut脚本去除包含多重子目录测试结果冗余问题
 
### Version 2.11
1. ut脚本修复未建mr跳过测试

### Version 2.10
1. ut脚本magic方法bug修复

### Version 2.9
1. ut脚本bazel test增加超时及测试报告文件检查
2. ut脚本start方法增加请求commit statues获取实际commit author
3. ut脚本upload方法增加入参校验

### Version 2.8
1. ut脚本加入app/(admin|interface)/main的测试
2. ut脚本加入结束执行后saga评论报告

### Version 2.7
1. 单元测试改为bazel test
2. 去除部分debug打印信息

### Version 2.6
1. 添加跳过check检测开关

### Version 2.5.3
1. upload接口异常打印日志

### Version 2.5.2
1. 加入bazel test相关处理

### Version 2.5.1
1. 优化ut脚本流程控制

### Version 2.5
1. ut脚本只执行service/main下dao层和library下的单测
2. 增加达标标准判断
    2.1 通过率=100% && 覆盖率>=30% && 当前pkg同比上次执行的单测覆盖增长率>=0
3. 输出当前执行概况日志

### Version 2.4
1. 删除编译失败的时候打包变化文件的机制
2. 测试权限

### Version 2.3
1.修正检测changelog大小写的问题
2.修改全量编译的job名称

### Version 2.2
1. 重构pipeline代码
2. 修改ut为不允许失败

### Version 2.1
1. cyclo函数复杂度调整到50
2. gometalinter过滤掉vendor包

### Version 2.0
1. 合并优化代码到lint.sh
2. 修改匹配变化文件的逻辑

### Version 1.0
1. 将compile重试次数改为1次
2. 修正lint_if_changed.sh匹配文件的bug