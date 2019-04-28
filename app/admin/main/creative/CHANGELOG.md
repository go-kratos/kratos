#### creative-admin

### Version 1.4.0
>1.创作学院支持搜索关键词配置  

### Version 1.3.1
>1.update account api rpc

### Version 1.3.0  
>1.任务后台1.0

### Version 1.2.8(for app 5.37)
>1.fixbug 合拍素材的download_url地址校验,必须带有acgvideo.com

### Version 1.2.7(for app 5.37)
>1.入口的新增,查询和更新增加两个字段(subtitle和whiteexp)
>2.增加平台platform参数的说明, 0:全平台,1:Android,2:iOS,3:iPad  
>3.滤镜素材添加渲染方式的下发,material.filter_type   
>4.增加合拍素材的download_url地址校验,必须带有acgvideo.com  

### Version 1.2.6
>1.创作学院v1.0后台调整  

### Version 1.2.5
>1.waitGroup使用规范化

### Version 1.2.4(for app 5.35)
>1.素材库支持合拍
>2.素材库支持主题库
>3.素材库分类支持 new tag
>4.贴纸素材配置支持白名单

### Version 1.2.3
>1.素材库分类支持多类素材管理 调整升级去重逻辑

### Version 1.2.2
>1.素材库投稿贴纸支持贴纸分类
>2.素材库拍摄贴纸新增画面效果类型
>3.bgm新增合拍入口

### Version 1.2.1
>1.素材库新增转场，贴纸新增贴纸类型按bitmask存储
>2.bgm绑定素材分类支持覆盖操作

### Version 1.2.0
>1.创作学院添加技能树管理

### Version 1.1.10
>1.创作学院过滤up主删除的专栏 

##### Version 1.1.9
>1.新增投稿贴纸

### Version 1.1.8
>1.新增tag设置rank值

### Version 1.1.7
>1.修改创作学院tag列表

##### Version 1.1.6
>1.新增bgm收录 站内信通知
>2.素材库支持贴纸及贴纸ICON

##### Version 1.1.5
>1.fix 稿件批量添加check new arc    

##### Version 1.1.4
>1.创作学院批量添加稿件逻辑，1)新老都有的，会自动过滤老的，提交新的 2) 全是新的，正常提交 3) 全是老的，提示已存在

##### Version 1.1.3
>1.tag支持排序  
>2.添加稿件支持热值计算  

##### Version 1.1.2
>1,升级了context ip
>2,bgm 批量tag
>3,素材库支持热词，贴纸
>4,bgm分类支持拍摄sort排序
>5,bgm收录通知（未完结feature）

##### Version 1.1.1
>1.fix bug bgm推荐列表展示frontname

##### Version 1.1.0
>1.bgm推荐列表展示frontname

##### Version 1.0.0
>1.app入口支持平台检索,支持统一的版本控制,支持 more
>2.素材库新增分类管理并支持素材归类及分类下排序
>3.bgm支持时间轴记录副歌起点
>4.素材库（包含bgm）统一日志采集  buiness=6  type 按照子业务分发

##### Version 0.6.1
>1.默认展示全部稿件   

##### Version 0.6.0
>1.创作学院接入搜索  
>2.单个稿件支持绑定多个分类标签    

##### Version 0.5.5
>1.bgm1.2支持frontname展示逻辑及新增其他未绑定tid检索

##### Version 0.5.4
>1.创作学院搜索不到稿件返回空  

##### Version 0.5.3
>1.优化搜索分页超时   

##### Version 0.5.2
>1.修复不存在business触发scan error问题  

##### Version 0.5.1
>1.修复不存在oid触发scan error问题   
>2.录入重复oid增加提示  

##### Version 0.5.0
>1.迁移音频库 并支持音频打点和tag管理
>2.新增 素材库的滤镜支持
>3.素材库支持platform 和 build exp 配置


##### Version 0.4.0
>1.创作学院管理  

##### Version 0.3.0
>1.新增字幕库和字体库

##### Version 0.2.2
>1.upgrade account rpc to version 3

##### Version 0.2.1
>1.fix GORM查询结果的时间字段的转换  

##### Version 0.2.0
>1.数据录入支持operations新类型collect_arc， 指代业务含义为：征稿启示,而非普通的公告   
>2.入口管理支持类型字段，区分创作中心和个人中心   
>3.白名单支持创作姬和粉版主APP   
>4.分区提示添加移动端APP专用的分区提示  

##### Version 0.1.0
>1.Update: upgrade Web Component to BM 

##### Version 0.0.7
>1.Update:init card by rpc and rewrite err to nil 

##### Version 0.0.6
>1.fix 504 close channel, 关闭错了    

##### Version 0.0.5
>1.fix panic 用户信息    

##### Version 0.0.4
>1.添加创作中心白名单管理后台的API接口      

##### Version 0.0.3
>1.修复插入新记录ctime和mtime为空问题,关闭入口更新mtime    
>2.app入口列表ctime/mtime/ptime仍返回时间戳格式  

##### Version 0.0.2
>1.更新Portal的时候同步更新ptime字段，预留ptime字段当做发布时间  
>2.有空值更新的时候，使用map[string]interface{}来更新数据，否则struct会过滤掉字段的空值  

##### Version 0.0.1
>1.初始化creative-admin
>2.宣发区分平台CRUD
>3.移动布局页面数据操作CRUD