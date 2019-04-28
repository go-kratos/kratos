
# manager-admin

### 1.4.4
1. permit grpc client start noblock

### 1.4.3
1. 支持 user state 配置

### 1.4.2
1. fix 用户没有权限点
2. fix manager用户不存在的情况

### 1.4.1
1. 支持grpc

### 1.4.0
1. 获取用户所有角色

### v1.3.9
1. manager_business_user_role表增加state状态处理

### v1.3.8
1. 增加internal获取控件接口

### v1.3.7
1. 修改角色列表排序方式

### v1.3.6
1. replace cfg name auth to permit

### v1.3.5
1. internal增加tag接口

### v1.3.4
1. tag控件增加业务id

### v1.3.3
1. permission接口增加权限组信息，角色信息

### v1.3.2
1. manager角色管理

### v1.3.1
1. manager-tag管理

### v1.3.0
1. 重构为使用BM框架

### v1.2.12
1. 在perm接口添加admin字段，返回是否是admin

### v1.2.11
1. 修改Auth接口，避免重复权限点

### v1.2.10
1. 增加批量查询用户部门接口

### v1.2.9
1. 静态资源zip文件上传接口 - 相同文件名覆盖的问题修复
2. Manager-admin迁移到Main路径下

### v1.2.8
1. 新增静态资源zip文件上传接口

### v1.2.7
1. 新增批量查询用户名接口 - Post TO Get

### v1.2.6
1. auth perms

### v1.2.5
1.fix auth sql

### v1.2.4
1. 加用户活跃心跳接口

### v1.2.1
1. 更新go-common

### v1.2.0
1. 稿件权限组相关

### v1.1.6
1. 新增分页获取用户信息接口

### v1.1.3
1. 提供auth

### v1.4.1
1. 新增理由管理

### v1.4.2
1. 新增常用字段

### v1.4.2
1. 移除role限制required

### v1.4.4
1. 修正角色id映射关系r.RID

### v1.4.5 
1. 添加行为日志
2. 添加具体部门和角色下属用户
3. 修正角色表auth_item
4. 修正下拉列表漏数据和配置为每次编辑人员
5. 调整lib和bid
