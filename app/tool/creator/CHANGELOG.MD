## Creater 代码生成器

### Version 1.2.3
> 1.优化生成library单测代码逻辑
> 2.修正部分bug

### Version 1.2.2
> 1.重构生成monkey代码为独立mock文件

### Version 1.2.1
> 1.生成monkey代码排除返回为空的方法

### Version 1.2.0
> 1.增加生成monkey代码

### Version 1.1.10
> 1.修复生成测试用例格式错位问题  
> 2.修复service层TestMain模板问题

### Version 1.1.9
> 1.修复 goimport 失败文件为空的问题.  
> 2.修复自定义 chan 类型无法 parser 的问题.

### Version 1.1.8
> 1.gen方法兼容windows路径  
> 2.测试方法追加bug修复

### Version 1.1.7
> 1.增加自动生成Reset方法  
> 2.create everything...

### Version 1.1.6
> 1.支持 service 层 ut 生成.  
> 2.生成的 mock 代码迁移至 mock 文件夹.  
> 3.摘除 gock 相关代码.  

### Version 1.1.5
> 1.摘除 genmock 代码 这波优化很清真 音風真的很严格.  
> 2.魔改 genmock 工具 fork 为 github.com/otokaze/mock/gomock.  

### Version 1.1.4
> 1.修改模板context.TODO()为context.Background()
> 2.优化*ast.FuncType的处理逻辑
> 3.修改tpTestFunc模板增加根convey方便用例扩展

### Version 1.1.3
> 1.新增生成mock代码方法  
> 2.修改TestMain模板中发布版本号为"docker-1"  
> 3.修改genInterface中多包生成对应Interface问题  
> 4.修复genMock.go中pkgName冲突问题  

### Version 1.1.2
> 1.新增解析Dao层结构体方法  
> 2.新增动态增加httpMock方法  
> 3.修复Goimport方法  

### Version 1.1.1
> 1.优化convey.So  

### Version 1.1.0
> 1.增加生成interface功能  
> 2.改进类型解析功能  
> 3.改进生成imports功能  
> 4.修改Dao层测试模板至更加规范化  

### Version 1.0.2
> 1.删除 TestMain 无用 flag.  

### Version 1.0.1
> 1.添加 context.Context 类型支持.  
  
### Version 1.0.0
> 1.无差别传参. 支持同时传入目录和文件, 生成自如!  
> 2.--func 指定 func 生成测试用例. (当前存在 *_test.go 会以 append 方式追加)  
> 3.--m 支持选择工作模式. (目前支持的有 test 测试用例...等)  
> 4.支持生成 dao_test.go 测试入口文件, 初始化 Dao 一把梭  

### Version 0.0.1
> 1.初始化项目.  