#### Macross 各种管理~

#### BGM
> http://girigiri.love/

### Version 3.0.9 - 2018.12.11
#### Features
> 1.修复 tools/unzip.go 的 bug：使用 “zip -r xxxx.zip xxxx/*”命令时，产物中不包含文件夹信息导致的解压失败

### Version 3.0.8 - 2018.12.4
#### Features
> 1.如果 apk 的渠道包生成失败，说明包有问题，删除整个 apk 包

### Version 3.0.7 - 2018.11.29
##### Features
> 1.unzip 下沉到 tools 中
> 2.取包服务部分业务逻辑下沉到 service 中
> 3.邮件新增附件功能

### Version 3.0.6 - 2018.11.15
##### Features
> 1.取包服务改为 GET

### Version 3.0.5 - 2018.11.14
##### Features
> 1.邮件、包上传、取包服务出错返回具体错误原因

### Version 3.0.4 - 2018.11.12
##### Features
> 1.macross 新增包上传 & apk 重签服务
> 2.macross 新增取包服务

### Version 3.0.3 - 2018.10.30
##### Features
> 1.macross 邮件发送服务新增抄送和密送功能，部分字段名称标准化

### Version 3.0.2 - 2018.10.29
##### Features
> 1.新增 macross 邮件发送服务

### Version 3.0.1 - 2018.09.17
##### Features
> 1.publish接口去掉res_size为0的校验  

### Version 3.0.0 - 2018.08.20
##### Features
> 1.初始化项目，全新逻辑  

### Version 2.2.0 - 2018.01.23
##### Features
> 1.config使用新SDK  
> 2.增加定期load model的逻辑  

### Version 2.3.0 - 2018.02.26
##### Features
> 1.创建、修改依赖图添加extra_info字段  
> 2.依赖图查询接口返回值添加extra_info字段  
> 3.暂时注释修改弱依赖图的model冲突判断  
> 4.上传framework时，校验模糊依赖关系的符号  
> 5.重新修改模糊依赖model的逻辑  
> 6.修改模糊依赖model的返回值去掉图中已存在的model  
> 7.增加查询和修改framwork依赖信息的接口  

### Version 2.2.0 - 2018.01.23
##### Features
> 1.依赖图相关功能  

### Version 2.1.2 - 2017.12.08
##### Features
> 1.精简逻辑，优化逻辑  

### Version 2.1.1 - 2017.12.05
##### Features & BUG_Fix
> 1.修复ios的fremawork上传接口参数model与当前文件不匹配的问题  

### Version 2.1.0 - 2017.12.05
##### Features & BUG_Fix
> 1.ios的fremawork上传接口添加用户和model权限校验  
> 2.ios的fremawork上传接口添加code和model唯一性校验  

### Version 2.0.0 - 2017.11.29
##### Features
> 1.ios的fremawork上传和列表查看接口  
> 2.ios的version管理接口  
> 3.ios的model、角色管理、权限查询接口  

##### Version 1.0.0
> 1.添加权限管理模块接口  
