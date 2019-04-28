#### tag-admin

#### Version 1.2.4
> 1.分区热门tag 显示tag为空时不更新rank表
> 2.迁移使用fanout，代替cache.Cache

#### Version 1.2.3
> 1.二级分区热门tag异步刷稿件资源.
> 2.迁移ES搜索查询稿件到SDK V3接口.
> 3.迁移账号info接口到GRPC info.

#### Version 1.2.2
> 1.频道海外版限制.
> 2.tag名称搜索排序以及相关度展示.

#### Version 1.2.1
> 1.增加相似频道.
> 2.更新分区tag展示逻辑.

#### Version 1.2.0
> 1.增加频道头图，长短简介，去除频道优先级.

#### Version 1.1.20
> 1.迁移es update tag接口到es sdk.

#### Version 1.1.19
> 1.分类下频道排序不生效问题;
> 2，频道创建时间一致的问题;
> 3，取消频道排序，置顶 修改最后编辑人问题;
> 4，无置顶频道 接口报错问题

#### Version 1.1.18
> 1.频道分类管理后台
> 2.置顶与排序频道

#### Version 1.1.17
> 1.删除c.RemoteIP()

#### Version 1.1.16
> 1.manager-search  search_type=arc

#### Version 1.1.15
> 1.search 升级    

#### Version 1.1.14
> 1.修复business appkey冲突

#### Version 1.1.13
> 1.去除verify
> 2.修复name==0判断问题

#### Version 1.1.12
> 1.新增business

#### Version 1.1.11
> 1.接口/x/admin/channel/info支持tname查询.
> 2.ES服务增加config配置文件.

#### Version 1.1.10
> 1.按照tag名称获取举报列表.

#### Version 1.1.9
> 1.升级elastic sdk

##### Version 1.1.8
> 1.测试使用bind default.

##### Version 1.1.7
> 1.更改模糊匹配按照匹配数返回.

##### Version 1.1.6
> 1.使用identify公共配置.

##### Version 1.1.5
> 1.摒弃主站搜索API,接入主站搜索SDK.
> 2.修复update tag接口问题.

##### Version 1.1.4
> 1.删除频道缓存变更清除和更新.

##### Version 1.1.3
> 1.增加频道缓存清除和更新.

##### Version 1.1.2
> 1.修复频道规则长度限制bug

##### Version 1.1.1
> 1.修复分区热门tag的id显示

##### Version 1.1.0
> 1.频道后台管理业务

##### Version 1.0.1
> 1.迁移tag搜索接口

##### Version 1.0.0
> 1.接入动态产品的管理后台审核功能
> 2.接入blademaster
> 3.调整后台权限点
> 4.增加清除tag cache的操作
> 5.增加同义词列表tname显示
> 6.迁移到business/admin/main目录
> 7.增加行为日志记录模块

##### Version 0.0.2
> 1.增加批量删除和批量处理
> 2.增加举报有效分

##### Version 0.0.1
> 1.初始化tag-admin
> 2.合并大仓库
