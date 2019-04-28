#### feedback 反馈系统

###### Version 4.2.1 - 2018.12.21
##### Features
> 1.H5用户反馈接口兼容buvid获取方式  

###### Version 4.2.0 - 2018.11.28
##### Features
> 1.去掉net/ip的依赖  

###### Version 4.1.8 - 2018.09.28
##### Features
> 1.新增H5专用用户反馈接口(支持多图片上传)  

###### Version 4.1.7 - 2018.09.28
##### Features
> 1.新增H5专用tag接口  

###### Version 4.1.6 - 2018.09.25
##### Features
> 1.为移动端新反馈途径增加tag  

###### Version 4.1.5 - 2018.09.21
##### Features
> 1.新增云视听和创作中心反馈  

###### Version 4.1.4 - 2018.09.20
##### Bugfix
> 1.修复播放器投屏兼容逻辑bug  

###### Version 4.1.3 - 2018.08.27
##### Features
> 1.播放器投屏新增tagID并针对其入库数据status做特殊处理  

###### Version 4.1.2 - 2018.08.30
##### Features
> 1.替换RemoteIP方法  

###### Version 4.1.1 - 2018.08.23
##### Features
> 1.增加和完善Dao层UT测试用例  

###### Version 4.1.0 - 2018.08.08
##### Features
> 1.修改bm  
> 2.修改identiry  

###### Version 4.0.6 - 2018.07.13
##### Features
> 1.删除bfs方法中的多余Fock  

###### Version 4.0.5 - 2018.06.26
##### Features
> 1.补充完善粉版直播做兼容逻辑  

###### Version 4.0.4 - 2018.06.26
##### Features
> 1.对粉版直播做兼容逻辑，每条反馈都单独建立session  

###### Version 4.0.3 - 2018.06.22
##### Features
> 1.修改粉版直播tag文案  

###### Version 4.0.2 - 2018.06.08
##### Features
> 1.修改粉版直播tag顺序  

###### Version 4.0.1 - 2018.05.14
##### Features
> 1.增加粉版直播逻辑  
> 1.修改reply接口对buvid的判断  

###### Version 4.0.0 - 2018.04.26
##### Features
> 1.项目整体迁目录  
> 2.切bm  
> 3.新增文件上传接口(content-type: octet-stream)  

###### Version 3.8.0 - 2018.03.07
##### Features & Bug
> 1.ugc/reply增加mid参数和校验  

###### Version 3.7.0 - 2017.12.15
##### Features & Bug
> 1.对接创作姬，为创作姬增加入口字段  
> 2.修复了初次创建session时，如果content或imgURL为空导致logURL无法写入的问题  

###### Version 3.6.1 - 2017.11.27
##### Bug
> 1.恢复ecode(用到了具体错误码信息)  
> 2.删除main的syscall.SIGSTOP  

###### Version 3.6.0 - 2017.11.20
##### Features
> 1.添加播放器检测上报接口  
> 2.删除无用代码(比如ecode)  

###### Version 3.5.4
> 1.修复线上Tag下发  

###### Version 3.5.3 
> 1.调整tag位置  

###### Version 3.5.2
> 1.增加session事务处理  

###### Version 3.5.1 
> 1.支持管理海外版搜索  

###### Version 3.5.0 
> 1.feedback合入Kratos

###### Version 3.4.4 
> 1.兼容IOS播放器上报传参错误  

###### Version 3.4.3  
> 1.创作中心接入文章反馈  

##### Version 3.4.2  
> 1.修复seesion bug    

##### Version 3.4.1
> 1.修复第一次番剧反馈多张图片bug  

##### Version 3.4.0  
> 1.更新vendor  
> 2.兼容PGC反馈数据展示  

##### Version 3.3.3  

> 1.修复用户回复状态 

##### Version 3.3.2  

> 1.测试ci打tag  

##### Version 3.3.1 

> 1.修复emoji过滤  

##### Version 3.3.0  

> 1.修复内容过滤正则  
> 2.创作中心子tag显示  

##### Version 3.2.2 

> 1.更新vendor  

##### Version 3.2.0  
   
> 1.接入新的配置中心   
> 2.过滤content内容  
> 3.去掉回复状态更新 

##### Version 3.1.1

> 1.拆分mobi-feedback  

##### Version 3.1.0  

> 1.区分播放器与移动端反馈平台会话   
> 2.修复tag显示     

##### Version 3.0.2

> 1.兼容管理后台多tag  

##### Version 3.0.1 

> 1.过滤移动端session  

##### Version 3.0.0

> 1.增加ugc反馈平台   

##### Version 2.3.0  

> 1.视频详情页提交新增会话  

##### Version 2.2.0

> 1.升级vendor  
> 2.返回客户端tag列表  
> 3.增加tag字段写入  

##### Version 2.1.3

> 1.升级vendor  

##### Version 2.1.2

> 1.升级vendor  

##### Version 2.1.1

> 1.升级vendor  

##### Version 2.1.0

> 1.配置支持从环境变量获取配置  
> 2.更新go-common与go-business依赖  

##### Version 2.0.0

> 1.配置中心接入  
> 2.go-business的错误码替换  
> 3.更新go-common与go-business依赖  

##### Version 1.0.4

> 1.获取IP地址修复  
> 2.govendor代替glide  

##### Version 1.0.3

> 1.限制用户发送内容的字符长度，配置文件可配  
> 2.FeedbackSsnNotExist的错误过滤，不打印日志  

##### Version 1.0.2

> 1.更新会话去掉mid  
> 2.会话创建更新laster_time  

##### Version 1.0.1

> 1.登录用户的会话逻辑  
> 2.登录用户的拉去反馈逻辑  

##### Version 1.0.0

> 1.增加反馈  
> 2.拉取反馈列表  
> 3.上传文件  
