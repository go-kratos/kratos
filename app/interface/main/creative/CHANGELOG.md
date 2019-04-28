### 创作中心

#### Version 2.41.7
>1.荣誉周报: 新增取消订阅功能

##### Version 2.41.6
>1.h5创作学院支持关键词配置 

##### Version 2.41.5
>1.任务体系：奖励激活添加双写逻辑    

##### Version 2.41.4
>1.给弹幕反馈中的弹幕保护和弹幕举报添加str的辅助字段  
  
  * id_str for danmu/protect/list 
  * dmid_str for danmu/report


##### Version 2.41.3
>1.重构特殊用户组查询接口,使用grpc并且进行批量查询  

##### Version 2.41.2
>1.修改h5任务提示文案  
>2.修改关注创作中心账号任务，完成条件  

##### Version 2.41.1
>1.增加data api接口，mcn使用 

##### Version 2.41.0
>1.重构商单,上线喽   

### Version 2.40.6
>1.任务体系：领取/激活接口添加redis锁 

### Version 2.40.5
>1.任务体系：fix修复makeup补偿接口  
>2.删除原来无用的老的稿件搜索代码  
>3.优化日志  

### Version 2.40.4
>1.任务体系新建奖励领取表，将领取数据进行双写  

### Version 2.40.3
>1.update unittest  

### Version 2.40.2
>1.周报展示实时生成

### Version 2.40.1
>1.重构public service到其他service的注入 

### Version 2.40.0（App 537）
>1.红点逻辑, 包括滤镜(2),BGM(3),拍摄贴纸(5)  
>2.素材之滤镜 增加字段表示滤镜类型(filter_type)  
>3.Faq增加pad相关的数据提供给h5  
>4.支持联合投稿staff的标记位的返回,包括当前用户在稿件中的角色    
>5.app首页的icons的双重灰度策略的支持,需要增加字段whiteexp(json)字段  
>6.icons的subtitle当做提示文案返回给app  
>7.合拍素材添加download_url字段,下发给app 
>8.首页接口的portal_list支持ios的pad device  
>9.联合投稿相关:app支持查询添加coop字段,并且返回的时候archive里面带attrs属性  

### Version 2.39.2
>1.修改荣誉周报显示hid=10 'x十万'->'x0万'

##### Version 2.39.1
>1.添加联合投稿接口

##### Version 2.39.0
>1.联合投稿

##### Version 2.38.16
>1.修复删除稿件报错无code问题

##### Version 2.38.15
>1.升级账号API为grpc服务,简化配置  

### Version 2.38.14
>1.创作学院：  
  * 1 技能树（原主题课程）接口取消分页，最大限制100条
  * 2 移除老接口：recommend与theme/course
>2.投稿一级分区顺序调整: 数码(188)移动到科技的下面   

### Version 2.38.13
>1.水印设置内网接口去掉同步开关强校验  

### Version 2.38.12
>1.修复周报parse负数错误

### Version 2.38.11
>1.水印设置内网接口增加同步开关  

### Version 2.38.10
>1.任务体系：  
  * 1 修改新手头像挂件通知消息
  * 2 修改组内任务顺序

### Version 2.38.9
>1.任务体系  
  * 1 进阶任务隐藏
  * 2 新手任务调整
  * 3 头像挂件奖励调整

### Version 2.38.8
>1.记录up主点击荣誉周报次数

### Version 2.38.7
>1.H5创作学院-v3.0发布  
  * 1 添加“推荐”与“精选”接口v2版本
  * 2 “精选”视频添加推荐理由
  * 3 将热门推荐的视频加入“推荐”接口

### Version 2.38.6
>1.修复honorState缓存

### Version 2.38.5
>1.修复rank parse错误

### Version 2.38.4
>1.增加荣誉周报降级开关

### Version 2.38.3
>1.修复高并发闭包查询的写入错误accSvc.UpInfo  
>2.app功能: 主题+合拍+投票 的入口全量发布  

### Version 2.38.2
>1.h5任务体系，fix任务列表接口返回-400  

### Version 2.38.1
>1.重构app ios的siderbar的内网接口  

### Version 2.38.0
>1.h5任务系统上线  

### Version 2.37.9
>1.重新开启荣誉周报入口  

### Version 2.37.8
>1.荣誉周报评级生成错误  

### Version 2.37.7
>1.创作学院：fix视频标签错误  

### Version 2.37.6
>1.荣誉周报修改策略,且新策略只对新生成周报生效  

### Version 2.37.5
>1.创作学院：视频列表接口，添加rights字段  

### Version 2.37.4（App 536）
>1.添加支持need_vote参数和函数，来判断不同的端是否需要投票信息  
>2.每一个素材类型的红点最新提示，使用material/pre的latests 
>3.app首页豆腐块的跳转链接由服务端动态下发,参考block_intros   

### Version 2.37.3
>1.添加高能看点internal查询接口

### Version 2.37.2
>1.升级投币修改服务为grpc

### Version 2.37.1
>1.任务体系：添加“用户任务完成状态列表”接口，提供给电磁力计划  

### Version 2.37.0
>1.app view 支持bgm list

### Version 2.36.14
>1.web3.5加入专栏数据   

### Version 2.36.13
>1.升级windows的tag预览到大数据的v2版本  

### Version 2.36.12
>1.查询带投票的稿件列表去除uid字段

### Version 2.36.11
>1.查询带投票的稿件列表
>2.修改荣誉周报文案

### Version 2.36.10
>1.上传大小白名单增加新档位，能上传8GB以上的单P视频，但是不能超过16GB 

### Version 2.36.9
>1.fix update: 过滤掉已经删除的分类,允许离散的有效绑定关系存在   

### Version 2.36.8
>1.update and fix: panic write map for ExtSidHotMapAndSort  

### Version 2.36.7
>1.update and fix: 兄弟sid集合应该从已有分类所属的推荐库里获取，而不是从基础bgm库里面索引，否则会造成没有aid关联数据     

### Version 2.36.6
>1.任务体系：添加挂件领取成功后发送消息   

### Version 2.36.5
>1.update and fix : 以音频加入到投稿背景音乐分类列表里面为准, music_with_category.ctime as jointime   

### Version 2.36.4
>1.合拍入口暂时关闭百分比灰度，仅限85个mid， from fangzhang   

### Version 2.36.3
>1.任务体系挂件接口由行为日志切换为调用业务方接口  
>2.添加粉丝数>=100时，自动解锁进阶任务  

### Version 2.36.2
>1.h5创作学院搜索和我的课程优化  
>2.upgrade bgm第一个Tag(NEW)的判断时间为mtime 
>3.h5数据中心的播放退出点增加0点计算

### Version 2.36.1
>1.h5数据中心的播放退出点均匀切分

### Version 2.36.0(App 535)
>1.无脑上传合拍以及主题的数据  
>2.合拍相关  
  * 1 CameraCfg.{$coo_max_sec|$coo_min_sec}
  * 2 AppModuleShowMap.cooperate  
  * 3 添加合拍的素材内容暴露,material/pre
  * 4 添加合拍独立的两个接口 cooperate/pre + cooperate/view   
  * 5 合拍计算热度   
  * 6 合拍接搜索，填充其他附属字段的数据信息  
>3.素材通用相关
  * 1.所有的素材都受到白名单的约束,type==18,暂时拍摄的贴纸上线白名单 
  * 2.添加theme主题下发的接口 type == 10  
  * 3.获取10种素材类型中的最新一个，素材使用ctime|BGM使用pubtime 
>4.LBS相关
  * 1.poi object view的时候app需要吐出
  * 2.need_poi 参数view，防止web查询携带不必要的数据  
>5.活动的相关代码迁移到public service 
>6.BGM扩展页面的数据需要对接搜索，提供bgm/ext接口给h5页面使用  
>7.App创作中心首页接入用户的任务系统   

### Version 2.35.10
>1.fix bug:字幕的数据丢了，而且数据被滤镜占用  

### Version 2.35.9
>1.任务系统，添加任务状态更新补偿接口，防止扫表遗漏导致任务状态无法更新  

### Version 2.35.8
>1.H5页面的跳转自定义化，后期产品和h5的调整会越来越频繁  

### Version 2.35.7
>1.fix：任务系统，任务和奖励缓存用法  
>2.TODO：防止重复奖励重复领取  

### Version 2.35.6
>1.荣誉周报部分文案修改    

### Version 2.35.5
>1.fix：任务系统，任务和奖励缓存用法    

### Version 2.35.4
>1.修复白名单map，并发读写异常    

### Version 2.35.3
>1.兼容editor上传null，然后过滤掉，防止进入垃圾数据    

### Version 2.35.2
>1.和服务端下发相关的素材类型兼容app上报数据集合或者单个数据, string/int   

### Version 2.35.1
>1.修复观看创作学院任务发消息  

### Version 2.35.0
>1.web任务系统上线  

### Version 2.34.1
>1.编辑器贴纸功能入口双端放量调为全量(videoup_sticker)  

### Version 2.34.0
>1.UGC内容付费相关  
  1. 最新协议的版本落地到配置文件，按照运营的需求强制配置更新 
  2. 付费稿件过审之后30天内不允许删除,落配置   
  3. 列表里面应该给出当前能否编辑和删除，且不能操作的理由 
  4. x/web/white接口里面还得给出ugcpay的标记位  

>2.fix：带id的素材上报的时候，不需要json格式化处理，去除多余的双引号  

### Version 2.33.13
>1.服务端开始下放ugcpay的标记位，from mochi， 为ugc内容付费做准备    

### Version 2.33.12
>1.消息助手兼容ios粉色和蓝色版本,mobi_app: iphone_b|iphone  

### Version 2.33.11
>1.荣誉周报文案修改 与->于
>2.荣誉周报hbase数据缓存到mc

### Version 2.33.10
>1.优化数据中心播放退出点。APP充电改为12个月

### Version 2.33.9
>1.APP弹幕也添加发送人的账号；关联关系；充电信息  
>2.去掉弹幕无用的接口和代码：web端的最新弹幕+web端的弹幕过滤相关  
>3.水印相关外网接口的h5化 

### Version 2.33.8
>1.UP荣誉周报文案修改

### Version 2.33.7
>1.UP荣誉周报 fix

### Version 2.33.6
>1.APP Index拉消息ios build改为>8220

### Version 2.33.5
>1.UP荣誉周报

### Version 2.33.4
>1.支持grpc server

### Version 2.33.3
>1.APP534之后up小助手升级，ios:8820;android:5332000之后就不需要调用消息接口了，app独立调用  
>2.fake数据参数配置化  

### Version 2.33.2
>1.减少app首页的banner，调整为5个  

### Version 2.33.1
>1.播放开关全量

### Version 2.33.0 (for APP5.34)
>1.素材相关：投稿贴纸需要新增加分类，参考滤镜
>2.给官方号推荐关注
>3.app首页index页面：给活动增加new和comment的字段
>4.增加BGM合作投稿的标记位和跳转URL
>5.APP付费标记ugcpay的支持和展现

### Version 2.32.14
>1.hot patch:创作中心添加创作学院的跳转url

### Version 2.32.13
>1.创作学院h5我的课程增加删除接口

### Version 2.32.12
>1.app数据中心相关的接口路由H5化----简单的稿件查询定位(archives/simple)     

### Version 2.32.11
>1.app数据中心相关的8个接口路由H5化，老的不影响    

### Version 2.32.10
>1.增加绿洲计划

### Version 2.32.9
>1.稿件字幕提交全量开放  

### Version 2.32.8
>1.去除成长计划状态查询接口

### Version 2.32.7
>1.修复申诉分页start小于total的边界case,三连判断内存分页算法，顺序不能换    

### Version 2.32.6
>1.增加创作学院h5去重功能

### Version 2.32.5
>1.APP首页的创作学院banner随机三个返回  

### Version 2.32.4
>1.增加专栏部分mock test

### Version 2.32.3
>1.更新bgm搜索的逻辑，仅支持搜索Up主的昵称，并且返回的时候把musicans字段替换成当前mid的昵称  

### Version 2.32.2
>1.创作学院h5部分接口去掉登录验证支持分享  

### Version 2.32.1
>1.创作学院修复视频计数

### Version 2.32.0
>1.创作学院h5上线  

### Version 2.31.4
>1.web端灰度特殊用户组：允许上传视频超过4G但是不允许超过8G, type=16   
>2.评论搜索添加状态:5,先发后审   

### Version 2.31.3
>1.app的icon状态后端返回使之动态化，archive/pre+archive/view   
>2.专场Transitions的素材下发支持pre和view  
>3.增加素材的上报,支持Transition,美颜，闪光灯，摄像头旋转，倒计时 
>4.app端native粉丝勋章相关H5化  

### Version 2.31.2
>1.去除成长计划查询相关

### Version 2.31.1
>1.评论搜索除以100，粉丝和勋章的H5上线，换个路由 

### Version 2.31.0
>1.3.5专栏数据需求  

### Version 2.30.15
>1.range里面的end默认统一传空字符串，减少ES索引创建的压力  

### Version 2.30.14
>1.oid有效的时候就不需要传o_mid了，索引建立跟不上   

### Version 2.30.13
>1.默认情况下(filter==-1)搜索所有的评论，不传ctime时间区间

### Version 2.30.12
>1.评论搜索重构，迁移到ES  

### Version 2.30.11
>1.faq加到redis的缓存,过期时间两分钟，防止智齿被打炸    
>2.app端native粉丝数据H5化  

### Version 2.30.10
>1.创作中心播放器开关迁移到whitelist
>2.开放5w粉丝up

### Version 2.30.9
>1.type为16的手机拍摄活动，如果活动提供有效的act_url,那么就直接用，活动负责页面对端的适配   (from mochi)
>2.增加h5充电相关的接口  (from app533)
>3.为防止消息助手提示重复，app端暂时只展现最新的一条  (from mochi)

### Version 2.30.8
>1.fix bug: 字幕和抽奖灰度策略互换 

### Version 2.30.7
>1.增加bgm的搜索接口，对接elastic search  

### Version 2.30.6
>1.给app的投稿贴纸的百分比设置为20%，暂时的这块运营不是很稳定, 3d人脸拍摄贴纸100%灰度    
>2.独立充电web端的router,为后期充电h5独立做准备   

### Version 2.30.5
>1.fix bug: 选择分区后默认强制获取所有平台的活动  

### Version 2.30.4
>1.app和web端的myinfo添加lottery字段，表示当前用户是否拥有抽奖资格  
>2.重构web端的数据预览(pre)接口   

### Version 2.30.3
>1.全端活动列表信息展现，包含type==16的手机拍摄活动，web和client可以兼容  

### Version 2.30.2
>1.Web投稿支持字幕的查询开关和个人的灰度开关  

### Version 2.30.1
>1.单个素材查询支持2D的投稿编辑器贴纸，也就是gif自定义水印  
>2.编辑器使用信息上报支持2d的投稿贴纸信息, videoup_stickers  

### Version 2.30.0
>1.dao单元测试更新  

### Version 2.29.7
>1.fix app material view without data 

### Version 2.29.6
>1.创作学院过滤up主删除的专栏  

### Version 2.29.5
>1.APP添加素材-投稿贴纸(videoup_sticker), type为7 

### Version 2.29.4
>1.添加动态抽奖的灰度字段 
>2.添加抽奖是否可修改的rules字段 

### Version 2.29.3
>1.创作学院支持按最热值排序，筛选已开发专栏   

### Version 2.29.2
>1.修改创作学院tag列表 

### Version 2.29.1
>1.给手摄活动的ActUrl代理转换一下,友情操作 

### Version 2.29.0(APP 532)
>1.ShowAcademy for APP(android interface+ios internal)

>2.Sticker

  - module_show_map
  - pre add sticker,hotword,single intro
  - view for 字幕，字体，滤镜，贴纸

>3.add CustomManagerTip(投稿贴士&问题反馈) for 投稿预览和编辑 

>1.去掉弹幕所属稿件的所有权判断，交给弹幕服务自己判断,涉及:修改弹幕池，操作弹幕属性，弹幕举报 

### Version 2.28.13
>1.取消APP投稿调用摄像头拍摄的灰度策略,全量开放  

### Version 2.28.12
>1.fix bug: 修复错误被覆盖的bug
>2.对-661的错误进行忽略，并增加对应的业务日志

### Version 2.28.11
>1.添加播放引导白名单  
>2.修改AddUpInfoCache的超时时间，使用模板查询的缓存时间，默认线上配置一个小时 

### Version 2.28.10
>1.fix bug: 删除稿件的时候，还需要清理mc里面对mid和title的定时限制，允许用户不断的删除&&提交相同名字的稿件 

### Version 2.28.9
>1.添加活动的protocol提示字段  

### Version 2.28.8
>1.新稿件搜索上线  

### Version 2.28.7
>1.新搜索回滚  

### Version 2.28.6
>1.创作学院标签支持排序  

### Version 2.28.5
>1.稿件列表接入新搜索   

### Version 2.28.4
>1.重命名，素材相关数据上报的接口使用内网接口:/upload/material 

### Version 2.28.3 (APP 5.31)
>1.FAQ
  - 视频编辑器H5的地址以及接口判断

>2.用户BGM反馈数据上报,需要对接对应的行为日志上报接口

>3.BGM的投稿行为和拍摄行为分类单独排序，需要增加新字段

##### Version 2.28.2
> 1.增加弹幕dao的超时时间，使用slow config 

##### Version 2.28.1
> 1.update 独立出errgroup的context，防止被cancel掉 

##### Version 2.28.0
> 1.增加up主开关功能

### Version 2.27.3
>1.feature:faq先上，配合产品需求

### Version 2.27.2
>1.fix bug:拍摄处于灰度中，热门活动中暂时排除【手摄】活动 

### Version 2.27.1
>1.myinfo时候兼容账号的rpc错误，把定时发布提示信息提前初始化,防止app奔溃 
>2.简化并删除不必要的cookie和access_key的传递 

### Version 2.27.0
>1.粉丝活跃度：接音频、专栏、动态、直播  
>2.视频列表新增点赞、分享数据  
>3.粉丝勋章页面添加领取总数和当前佩戴总数  

### Version 2.26.9
>1.H5的mission/type接口带上用户from标记位，区分web和app，加上首摄的活动 

### Version 2.26.8
>1.tag为空的时候，data不返回null，需要默认初始化为{}   

### Version 2.26.7
>1.创作学院接入搜索  

### Version 2.26.6
>1.弹幕管理相关的Post请求添加独立的限速handler 

### Version 2.26.5
>1.升级ip获取的方式，使用metadata.RemoteIP

### Version 2.26.4
>1.添加不允许投稿的黑名单,需要配置支持
>2.添加账号信息为空的时候的日志记录  

### Version 2.26.3
>1.增加whitelist nil判断

### Version 2.26.2
>1.弹幕查询和最新弹幕接口内部接口升级  

### Version 2.26.1
>1.移动端530数据需求 

### Version 2.26.0
>1.App/Index页面新增和优化
  - 4个活动，最多四个，rank排序asc 
  - UpNotify, 获取最新两个,时间转成int时间戳（接入消息通知，注意超时时间） 
  - 首页添加优先级最高的四个活动, /x/app/index接口返回四个推荐活动 (注意：首页这里添加的是征稿活动，不是普通的活动列表)

>2.App/Index Portal增加more字段  
  - 增加带有子类别的more标记位的scheme入口 

>3.滤镜素材在兼容原来的基础上，增加素材分类
  - 原有的排序规则不变
  - 增加素材分类,存储分类的到素材主键的映射关系，分类本身有权重，映射关系中本身也有权重 （增加第二套排序方案： 分类排序和素材自身排序）

>4.Bgm
  - 增加新字段Recommend，指代是否被运营推荐
  - 新增view接口 

>5.活动为【手摄】合并到app做准备

### Version 2.25.10
>1.1.5.30的H5版本，活动热门标记提前发布，线上接口和数据已经准备好了 

### Version 2.25.9
>1.rebuild library bm's logger 

### Version 2.25.8
>1.identify服务拆分为authSvc和VerfySvc 

### Version 2.25.7
>1.fix danmu dao API res.Object is nil

### Version 2.25.6
>1.xints迁移到model

### Version 2.25.5
>1.为5.30和活动线上资源联调而提前开启app的拍摄模块灰度功能 

### Version 2.25.4
>1.优化tag的返回，活动tag优先校验
>2.完善tag不允许使用的nil判断 
>3.弹幕过滤操作的时候，如果dmids长度为0则直接返回  

### Version 2.25.3
>1.fix bug:Article Capture CJSON return

### Version 2.25.2
>1.fix bug:music的重命名的锅   
>2.feature:根据用户昵称精确查询用户mid   

### Version 2.25.1
>1.创作学院教程的反馈和建议数据报表需增加提供建议的用户ID  
>2.更新抓图接口log  

### Version 2.25.0
>1.音频支持别名展现,如果有别名就替换成别名 
>2.删除稿件aid推荐关联的相关旧代码 

### Version 2.24.6
>1.修改成长计划入口逻辑,兼容大数据炸了的情况
>2.APP投稿滤镜和录音模块全量开放 

### Version 2.24.5
>1.老版hbase查询数据为空不做拦截处理,防止影响上层逻辑    

### Version 2.24.4
>1.升级hbase client  

### Version 2.24.3
>1.首页添加专栏分享数据(rpc自动获取)  
>1.recover水印设置panic  

### Version 2.24.2
>1.fix:弹幕反馈列表里面解析时间的时候需要带上Localtion信息，使用time.ParseInLocation  

### Version 2.24.1
>1.添加独立为添加稿件使用的geetest验证接口，标记位为true的时候前端才使用验证服务  

### Version 2.24.0
>1.移动端增加专栏数据  
>2.优化创作学院分页搜索  
>3.修复视频播放留存率  

### Version 2.23.5
>1.修复弹幕report的username 

### Version 2.23.4
>1.up-service切换成rpc模式
>2.修复弹幕举报从举报人变为弹幕发送人mid

### Version 2.23.3
>1.mid投稿限速的计算标记位提示返回
>2.投稿推荐tag上报给推荐引擎client_type, 0:web;1:app

### Version 2.23.2
>1.创作学院稿件返回时长    

### Version 2.23.1
>1.创作学院专栏列表返回简介  

### Version 2.23.0
>1.增加滤镜素材的输出和灰度的标记位
>2.增加bgm.music的tags和时间轴标记信息 
>3.素材数据的请求对接build号的大于小于限制和platform类型 
>4.APP5.28

### Version 2.22.1
>1.修复创作学院一级标签筛选问题  
>2.修复移动端粉丝排行关注状态    

### Version 2.22.0
>1.创作学院  

### Version 2.21.20
>1.bug fix:水印type默认为0是非法值,修改默认为1,用户昵称类型 

### Version 2.21.19
>1.修复弹幕专业错误返回

### Version 2.21.18
>1.rpc 聚合
>2.活动tag为不带逗号的字符串，并且tag校验的适合自我校验是否未活动tag 

### Version 2.21.17
>1.跳过

### Version 2.21.16
>1.移动投稿入口全量开放，并提示LV1以下的用户请先答题 

### Version 2.21.15
>1.放ssrf增加禁止30x跳转

### Version 2.21.14
>1.App视频添加字幕模块全量开放

### Version 2.21.13
>1.投稿入口下线二级分区:tid=175 

### Version 2.21.12
>1.capture接口防ssrf

### Version 2.21.11
>1.修复首页点赞数据  

### Version 2.21.10
>1.升级http组件到baldemaster

### Version 2.21.9
>1.Web端首页返回分享数据，30天增量数据返回点赞

### Version 2.21.8
>1.Web端的一级分区进行自定义排序 
>2.大数据的http请求都添加上mid的参数内容 
>3.来自产品的更新策略，app暂时没有完整支持，所以APP端的字幕效果暂时只接受白名单 

### Version 2.21.7
>1.投稿web端预览的时候添加用户投稿习惯数据:热门分区 
>2.web端提示添加type=7的带跳转的投稿小贴士 
>3.app端提示添加type=8的带跳转的H5引导小贴士 
>4.fix:app粉丝数据在ios平台中返回null的错误  

### Version 2.21.6 
>1.移动端最新评论只返回视频的

### Version 2.21.5
>1.移动端最新评论过滤音频  

### Version 2.21.4
>1.修复app数据概览返回daytime  

### Version 2.21.3
>1.patch:app查看稿件的时候也需要module_show    

### Version 2.21.2
>1.patch:app字幕模块的灰度兼容投稿的mid白名单列表   

### Version 2.21.1
>1.接入评论音频和首页展示文章音频评论  

### Version 2.21.0
>1.app粉丝管理和数据分析  

### Version 2.20.0
>1.主APP粉版本提供弹幕操作的一系列接口 
>2.投稿预览接口里面添加客户端字幕subtitle灰度的标记位 
>3.提供素材库查询的接口 
>4.弹幕列表和最近弹幕都添加封面字段 

### Version 2.19.11
>1.patch:APP投稿预览的时候如果是未初始化的新用户账号的水印信息，就通过url判断来默认开启水印状态 

### Version 2.19.10
>1.修复app首页5天数据panic    
>2.异步routine使用自己生成的context   

### Version 2.19.9
>1.优化app首页数据接口  

### Version 2.19.8
>1.patch:5.26之前的app版本在查看单个稿件的时候默认去除活动tag 

### Version 2.19.7
>1.给H5和APP端提供接口，通过tid分区id来获取当前有效的活动列表 

### Version 2.19.6
>1.向前端兼容活动列表，返回为空的时候返回空的json数组 

### Version 2.19.5
>1.修复文章评论标题  

### Version 2.19.4
>1.update bazel BUILD files

### Version 2.19.3
>1.移动端分区优化：改变UGC分区排序,生活、游戏、舞蹈、音乐、时尚、娱乐、科技、动画、影视、鬼畜、国创 
>2.增加隐藏分区：纪录片、电视剧、番剧 
>3.添加错误码:20052,当前稿件的私单已被禁止在前端展现 

### Version 2.19.2
>1.根据评论返回的oid查询标题和封面，根据mid查询用户名和图像  

### Version 2.19.1
>1.推荐tag列表去除活动tag,包括活动自带的tag集合和活动当前的名称 
>2.投稿预览的时候tags永远只出现一个，无tag情况下未活动名称  

### Version 2.19.0
>1.App Index Banner接入resource service的统一管理和投放 
>2.fix bug: slice遍历中不允许修改元素内容之后再重新计算,会造成越界 

### Version 2.18.10
>1.fix bug:当某一个策略组的gbm全部被下架的时候，这个分区不应该返回给前端  
>2.feature:musicians字段改成up主的昵称 

### Version 2.18.9
>1.pre接口里面活动信息要展现活动的tags，给h5前端做判断，在tags未空的情况下，默认把活动名字当做活动tag 

### Version 2.18.8
>1.删除web端的bgm路由 
>2.增加获取稿件游戏信息接口

### Version 2.18.7
>1.粉版BGM相关的接口支持:预览和查询分区两个接口 
>2.APP H5接口投稿预览的时候，返回活动相关的信息，用于做和native的tag操作相关联 
>3.单稿件查询结果的字段编辑规则进行了更新，新增dynamic和mission_tag:前者和描述字段一致，后者只有打回的稿件允许修改 

### Version 2.18.6
>1.粉丝管理增加活跃度数据    

### Version 2.18.5
>1.更新聚合白名单接口返回字段   

### Version 2.18.4
>1.提供目前所有白名单聚合接口  

### Version 2.18.3
>1.播放分析数据返回稿件标题  

### Version 2.18.2
>1.优化app稿件查询接口  

### Version 2.18.1
>1.申诉详情返回用户信息  
>2.app粉丝数据容错处理  

### Version 2.18.0
>1.数据中心三期需求 

### Version 2.17.7
>1.水印支持预览和View查询,在APP和Web两端
>2.添加waterMarkSetInternal内网生成水印的接口 

### Version 2.17.6
>1.AccountRPC升级使用Discovery,需要简化配置 

### Version 2.17.5
>1.Fix Bug For Tip：首页和稿件列表页面的提示信息，在下架了所有的内容之后，没有做对应的清理操作 
>2.私单的添加Warn日志 

### Version 2.17.4
>1.提供APP内网查询稿件的商单下的游戏信息的接口，通过aid查询对应的游戏信息，调用链路: order->game 

### Version 2.17.3
>1.兼容hbase脏数据的读取  

### Version 2.17.2
>1.app粉丝总量和增量获取实时数据  

### Version 2.17.1
>1.增加绿洲计划统计分析结果的数据接口, /cm/oasis/stat 

### Version 2.17.0
>1.增加移动端粉丝数据  
>2.申诉评分不判断申诉的状态    

### Version 2.16.1
>1.修复添加申诉bug   

### Version 2.16.0
>1.creative迁移到interface/main目录      

### Version 2.15.8
>1.申诉限制     

### Version 2.15.7
> 1.去处对数据库依赖，从up-service读取specialup的用户信息
> 2.去掉无用的对client端的编码率2500的内测组,现在已经全部6k码率 

### Version 2.15.6
> 1. support for MainAPP 2.25添加首页和稿件列表的紧急通知和公告 

### Version 2.15.5
> 1.去掉Relation3 RPC无效的日志输出 

### Version 2.15.4
> 1.account Profile双重判断是否被屏蔽  

### Version 2.15.3
> 1.使用account-service v7  

### Version 2.15.2
> 1.增加搜索ps限制

### Version 2.15.1
> 1.update: 添加cache的时候使用context.TODO() 

### Version 2.15.0
> 1.增加up/info 内网接口  

### Version 2.14.0
> 1.粉丝关注数据增加按月筛选  

### Version 2.13.4
> 1.去掉老的无用的投稿封面查询的配置和代码，线上30002端口的大数据服务已经下线 

### Version 2.13.3
> 1.Fix APP Banner的remark字段是否有值,之前的sql scan没获取出来  

### Version 2.13.2
> 1.APP5.24开始Banner的显示需要manager后台相关负责人确认之后才能给展现机会(operation.remark字段) 

### Version 2.13.1
> 1.修改up-service接口uri  

### Version 2.13.0
> 1.成长计划增加返回字段  

### Version 2.12.5
> 1.fix 主APP的白名单强制mid，和当前等级没有关系,应该去掉等级限制 
- case :mid不再灰度名单, level == 5,可以在老稿件里面添加视频

### Version 2.12.4
> 1.断言的类型必须一致，否则断言失败，value会为空值，white一直为0 

### Version 2.12.3
> 1.主App投稿查询，先判断是不是在白名单里面，如果不在白名单里面，则强制不允许在已有的稿件中添加分P 

### Version 2.12.2
> 1.PortalConfig add mc cache

### Version 2.12.1
> 1.-6,-7,-8三种待审状态的稿件不允许修改定时发布时间 
> 2.投稿提示文案默认为: "成为UP主，分享你的创作" 
> 3.格式化文案修复逻辑bug

### Version 2.12.0
> 1.支持主App添加稿件，查询稿件，检查tag等操作

### Version 2.11.12
> 1.去掉statsd的无用配置 
> 2.增加缓存开关

### Version 2.11.11
> 1.关闭数据缓存

### Version 2.11.10
> 1.增加商单信息的订单序号字段id_code 

### Version 2.11.9
> 1.修复申诉安全漏洞

### Version 2.11.8
> 1.开启数据缓存
> 2.弹幕举报增加mid判断
> 3.反馈增加mid

### Version 2.11.7
> 1.优化数据问题

### Version 2.11.6
> 1.修复数据问题

### Version 2.11.5
> 1.增加首页数据白名单。增加相应日志

### Version 2.11.4
> 1.优化协管的弹幕操作接口，取消弹幕状态和池类别修改的老接口，统一使用internal/v2/dm/edit/{state|pool} 

### Version 2.11.3
> 1.前端视频审核压力提示增加最高级别提示: 阻塞, Level值为5  

### Version 2.11.2
> 1.更新弹幕编辑接口uri 

### Version 2.11.1
> 1.去掉Covers的使用，使用AICovers

### Version 2.11.0
> 1.申诉重构  

### Version 2.10.22
> 1.编辑之前查看单个稿件信息，包括充电，修改为异步充电查询 

### Version 2.10.21
> 1.弹幕状态设置增加mid判断

### Version 2.10.20
> 1.修改文案： "成长计划" => "创作激励"  

### Version 2.10.19
> 1. 优化私单缓存命中为空的问题，当aid尚未参与私单计划，然后给一个默认初始化对象的cache值 

### Version 2.10.18
> 1.update and fix 成长计划 bug

        case: 允许以前投稿过，但是现在全删除,或者官方原因被全部下降，BUT粉丝数和累积的播放量满足条件

### Version 2.10.17
> 1.hbase timeout设置

### Version 2.10.16
> 1.creator 数据中心添加白名单  

### Version 2.10.15
> 1.app 创作中心入口添加白名单  

### Version 2.10.14
> 1.fix 安全bug，删除稿件之前需要校验稿件的mid所属  

### Version 2.10.13
> 1.web的pre预览接口添加video_jam字段，指示当前审核和转码的压力状态 

### Version 2.10.12
>1.稿件列表只展示中文分区    

### Version 2.10.11
>1.修复抓图接口判断图片类型    

### Version 2.10.10
>1.多语言支持 之 分区支持日文 
>2.多语言支持 之 简介格式化支持日文 

### Version 2.10.9
> 1.私单游戏列表查询支持一级索引带上首字母   

### Version 2.10.8
> 1.专栏逻辑迁移  

### Version 2.10.7
> 1.弹幕列表接口为CID失效进行容错处理  

### Version 2.10.6
> 1.投稿等级白名单用户可以绕过帐号等级的验证 
> 2.游戏列表按照ID倒序排列 
> 3.去掉了老的pre里面的porder相关的mid列表 

### Version 2.10.5
> 1.私单增加缓存策略 

### Version 2.10.4
> 1.水印db操作去掉stmt    
> 2.check MD5Sum    

### Version 2.10.3
> 1.更新单元测试以符合saga的要求  
> 2.添加昵称修改触发水印修改的手机短信告警  

### Version 2.10.2 
> 1.申诉兼容  

### Version 2.10.1
> 1.修復申訴bug

### Version 2.10.0
> 1.对接私单项目之游戏广告交易平台 

### Version 2.9.10
> 1.feature：加强申诉限制，相同aid在处理期间不能进行二次申诉操作 

### Version 2.9.9 
> 1.创作姬的简介编辑支持提示文字长度限制  
> 2.创作姬特定分区和创作类型下支持简介扩充到2k字符    

### Version 2.9.8 
> 1.添加活动协议的查询接口: /x/web/mission/protocol   

### Version 2.9.7 
> 1.修复appeal list panic

#### Version 2.9.6 
> 1.创作姬公测改为内测  
> 2.申诉回复图片修复  

#### Version 2.9.5 
> 1.申诉接入workflow  

#### Version 2.9.4
>1.tag终于找到产品了，操作和主站一致,只有活动tag不允许，其他的都允许，分类TAG也允许，终于知道如何上热门了  
>2.弹幕反馈的report视频id标记错误，给前端的还是得用dm_inid  

#### Version 2.9.3
>1.数据正确性问题,内容tag不能禁止 

#### Version 2.9.2
>1.Fix过滤创作姬的tag状态和类型 
>2.tag推荐增加mid

#### Version 2.9.1
>1.弹幕反馈列表和弹幕反馈保护列表重构完成 


#### Version 2.9.0
>1.增加私单内网查询接口

#### Version 2.8.5
>1.log基础库更新

#### Version 2.8.4
>1.添加视频审核时长等级接口
>2.去掉了CheckIsFriend方法调用

#### Version 2.8.3
>1.成长计划的用户限制添加到配置文件中 

#### Version 2.8.2
>1.文章专栏提交增加三个配置项限制,内容最大,字数Max,字数Min    
>2.评论列表支持白名单 

#### Version 2.8.1
>1.粉丝限制改成3个   

#### Version 2.8.0
>1.对接商业产品成长计划的7个接口

#### Version 2.7.6
>1.专栏增加原图字段

#### Version 2.7.5
>1.专栏投稿支持活动合作关联  

#### Version 2.7.4
>1.弹幕列表需要带上aid, 返回的时候增加分P标题和稿件的aid以及标题   

#### Version 2.7.3
>1.抓图接口增加jpg类型校验

#### Version 2.7.2
>1.弹幕分布间隔时间设置为1秒  

#### Version 2.7.1
>1.danmu/filter/list 弹幕列表接口重构,添加三个上中逆向行为的弹幕标记位 

#### Version 2.7.0
>1.专栏增加动态推荐语

#### Version 2.6.1
>1.过滤type=2自动解析filter字符串的值,如果是数字,强转成hash; 如果是字符串,就原样传输  

#### Version 2.6.0
>1.创作姬v4    

#### Version 2.5.0
>1.专栏返回图片空间大小   

#### Version 2.4.13
>1.弹幕转移支持参数offset,弹幕偏移量(单位为秒),浮点数,2bit精度,支持负数 

#### Version 2.4.12
>1.upinfo mc prom优化

#### Version 2.4.11
>1.透传弹幕转移操作不当的逻辑错误,提示给用户
>2.增加up主信息缓存时间

#### Version 2.4.10
>1.弹幕分布接口 查询ArchiveService RPC接口的时候对NotFound进行容错处理  
>2.弹幕分布接口 过滤非用户主动上传的cid,节省查询弹幕分布数据的次数  

#### Version 2.4.9
>1.弹幕相关支持替换弹幕分布 

#### Version 2.4.8
>1.必须先初始化返回值 

#### Version 2.4.7
>1.最新弹幕重构: 支持 分页;按照类型搜索;支持排序;支持批量编辑操作   

#### Version 2.4.6
>1.替换账号未登录错误码  

#### Version 2.4.5
>1.默认强制展示第1页的前1000个最新弹幕  

#### Version 2.4.4
>1.修改稿件列表描述字段   
>2.文章和草稿列表增加日志  

#### Version 2.4.3
>1.最新弹幕支持分页   

#### Version 2.4.2
>1.修改文章内容大小限制为1M  

#### Version 2.4.1
>1.弹幕相关第二期:弹幕过滤，新增两个接口：过滤列表(/danmu/filter/list)和添加屏蔽词(/danmu/filter/edit)   

#### Version 2.4.0
>1.创作姬校验文章作者

#### Version 2.3.1
>1.修复创作姬视频退出点数 

#### Version 2.3.0
>1.创作姬v3版本   

#### Version 2.2.8
>1.实名制账号查询使用Pb3的rpc接口，能快不少，稳定性也提升      
>2.稿件封面推荐完全去除redis依赖

#### Version 2.2.7
>1.HotFix:最新弹幕添加Aid，方便跳转到对应详情页面  

#### Version 2.2.6
>1.Creative添加自定义Prometheus的监控指标，主要是三个方面：
    
    article和data的Cache命中和失效计数；
    极验的第三方业务处理失败；
    删除稿件的业务计数打点优化；

#### Version 2.2.5
>1.Recent接口添加 FontSize,Color,Mode三个字段  

#### Version 2.2.4
>1.智能推荐封面重构、去除bfs依赖

#### Version 2.2.3
>1.添加tag查询接口,过滤tag不存在的错误情况   

#### Version 2.2.2
>1.弹幕管理的5个接口, 最近弹幕，弹幕列表，弹幕池修改，弹幕状态修改，弹幕转移    

#### Version 2.2.1
>1.修复创作姬宣发列表  

#### Version 2.2.0
>1.创作姬v2版本 

#### Version 2.1.4
>1.errgroup替换包

#### Version 2.1.3
>1.移动端入口可配置

#### Version 2.1.2
>1.UpInfo鉴权的时候忽略错误ArtCreationNoPrivilege,但是保留其他类型权限判断的错误  

#### Version 2.1.1
>1.增加upinfo缓存时间

#### Version 2.1.0
>1.增加up主状态查询接口

#### Version 2.0.5
>1.增加Service服务层的普罗米修斯监控的详细打点数据  
>2.修复极验的timeout参数配置    

#### Version 2.0.4
>1.增加极验timeout

#### Version 2.0.3
>1.数据中心添加白名单功能  

#### Version 2.0.2
>1.修复创作姬数据中心

#### Version 2.0.1
>1.修复首页专栏最新评论  

#### Version 2.0.0
>1.创作姬APP接口上线

#### Version 1.45.1
>1.修复关注数据为每日更新  

#### Version 1.45.0
>1.【宣发中心】Web端宣发消息兼容全平台的platform类型,(0||2)  
>2.【宣发中心】App端宣发消息兼容全平台的platform类型,(0||1)  

#### Version 1.44.6
>1.禁止Hbase的信息发生修改而发送日志到容器的stdout   

#### Version 1.44.5
>1.添加支持稿件的动态字段Dynamic展现  

#### Version 1.44.4
>1.修复Bug， RPC3的调用严格以返回值错误为标准， 并且service里预先初始化返回值，并且也按照错误码为判断标准做后续计算  

#### Version 1.44.3
>1.修复文章提交的mc缓存设置的格式类型，BigEndian之后的Object按照FlagJSON存储,可以参照videoup    

#### Version 1.44.2
>1.推荐的日志之前误写成ERROR,会影响日志查询  
>2.fix ecode.To 到 ecode.Int(res.Code)的转换  

#### Version 1.44.1
>1.fix 数据中心的mc获取数据添加引用符号

#### Version 1.44.0
>1.creative合并大仓库

#### Version 1.43.3
>1.修复抓图bug  

#### Version 1.43.2
>1.校验抓取图片的url   
>2.下载图片设置800ms超时控制    

#### Version 1.43.1
>1.推荐的查询和更新入口，进行重构，去掉对dede数据库的依赖  

#### Version 1.43.0
>1.新增获取弹幕分p列表接口

#### Version 1.42.0
>1.新增获取弹幕分p列表接口

#### Version 1.41.0
>1.新增高级弹幕管理接口

#### Version 1.40.1
>1.创作中心接入Protobuffer的RPC3接口， 单稿件信息查询，多稿件信息查询，多稿件Stats信息查询   

#### Version 1.40.0
>1.智能推荐封面接口

#### Version 1.39.1
>1.创作中心接口实名制开启 文章submit/update  
>2.article/pre 添加实名信息字段  

#### Version 1.39.0
>1.增加粉丝勋章重命名功能

#### Version 1.38.1
>1.充电留言透传cookie, ak用于实名支持
>2.评论切换到internal接口

#### Version 1.38.0
>1.提供图片链接抓取上传bfs接口

#### Version 1.37.5
>1.调整APP文章列表格式

#### Version 1.37.4
>1.开放新版APP小视频、相簿入口

#### Version 1.37.3
>1.实现网安添加Nid和MD5值的外部接口
>2.异步反馈给网安的地址:https://wax.gtloadbalance.cn:38080/InterfaceInfo/GetResult   

#### Version 1.37.2
>1.去除河童子数据接口依赖-更新

#### Version 1.37.1
>1.去除河童子数据接口依赖

#### Version 1.37.0
>1.数据中心二期  
 
#### Version 1.36.5
>1.重构清理商单和私单对于MID的判断属性，porders和orders都只表示当前用户有资格参与对应的商单或者私单    

#### Version 1.36.4
>1.更改app专栏入口的不带new icon  

#### Version 1.36.3
>1.迁移稿件模板    
>2.更改app专栏入口的icon    

#### Version 1.36.2
>1.operations(公告和话题)添加过滤未来即将发布的数据, stime <= now

#### Version 1.36.1
>1.web获取最近充电列表   

#### Version 1.36.0
>1.添加宣发中心两个接口，分别获取创作中心的各个平台的版本更新列表,话题和公告   

#### Version 1.35.2
>1.app专栏入口上线 

#### Version 1.35.1
>1.安卓开放app专栏入口    

#### Version 1.35.0
>1.app专栏入口    

#### Version 1.34.1
>1.优化: 简介格式化的数据从videoup-service的接口获取  
>2.优化: 封面数据直接从redis集群里面获取  
>3.综上所述，creative就已经完成去除对arc数据库的依赖  

#### Version 1.34.0
>1.文章列表展示硬币计数  

#### Version 1.33.7
>1.数据缓存清理dao, 线上数据源替换

#### Version 1.33.6
>1.稿件推荐增加profile校验

#### Version 1.33.5
>1.申诉库迁移到creative库中

#### Version 1.33.4
>1.投稿封面自动截图的重构，先查Redis，再查MySQL DB，等过一定时间之后，再全量切换到Redis  
>2.重构redis的配置，发布前需要把videoup-job的redis tw配置加上  

#### Version 1.33.3
>1.fix推荐列表的作者Auth字段，另调用账号接口批量查询当前mid对应的用户名称  

#### Version 1.33.2
>1.稿件列表容错，增加videoup-service数据源
>2.关联稿件推荐service拆分

#### Version 1.33.1
>1.web首页数据源替换，和APP统一

#### Version 1.33.0
>1.骑士2.0.
>2.可选择添加主站或者直播的骑士. 扩展添加协管接口
>3.骑士列表，增加封禁字段，增加类型字段{main:1, live:1}
>4.接入其他直播骑士相关的接口，查询、撤销等

#### Version 1.32.0
>1.app首页展示文章
>2.app粉丝勋章    

#### Version 1.31.2
>1.文章追加日志  
>2.app首页评论数据容错  

#### Version 1.31.1
>1.充电查询的DAO重构: ArchiveState 和 UserState  

#### Version 1.31.0
>1.用户加入退出充电返回之前状态

#### Version 1.30.1
>1.创作中心对接商业产品部门的成长计划，生成对应标记位grow_up  

#### Version 1.30.0
>1.增加粉丝勋章排行榜功能

#### Version 1.29.1
>1.HotFixAllowCommercial 包括商业产品部门的用户     

#### Version 1.29.0
>1.支持部分特殊分区的简介格式化    

#### Version 1.28.8
>1.app view兼容desc format

#### Version 1.28.7
>1.app首页下掉用户调研

#### Version 1.28.6
>1.修改水印图片保存路径

#### Version 1.28.5
>1.文章图片上传debug 

#### Version 1.28.4
>1.增加当没有5天数据时的处理

#### Version 1.28.3
>1.稿件列表size限制20
>2.identity切换到internal

#### Version 1.28.2
>1.APP首页基础数据替换成hbase来源
>2.APP5天数据增量使用真实数据

#### Version 1.28.1 
>1.申诉评论限速      
>2.app用户调研地址更新  

#### Version 1.28.0
> 1.图文数据统计上线

#### Version 1.27.0  
> 1.配置是否开启job消费databus消息更新水印

#### Version 1.26.0  
> 1.app1.3 首页 

#### Version 1.25.2 
>1.用户查看申诉更新为已读    
>2.修改校验文章作者权限日志    

#### Version 1.25.1
>1.文章白名单权限校验

#### Version 1.25.0  
> 1.私单上线  
> 2.私单灰度策略: 私单报备用户组  

#### Version 1.24.0
>1.专栏问题反馈

#### Version 1.23.2
>1.查看稿件支持白名单

#### Version 1.23.1
>1.图文防重复缓存优化

#### Version 1.23.0
>1.水印重构  

#### Version 1.22.0
>1.粉丝勋章V1.1 校验名称

#### Version 1.21.1
>1.文章支持上传gif图片,配置信息写入配置文件

#### Version 1.21.0
> 1.创作中心增长、删除、查看协管
> 1.创作中心协管日志相关接口

#### Version 1.20.2
> 1.修改mc format

#### Version 1.20.1
> 1.修改set mc raw

#### Version 1.20.0
> 1.实名认证手机或者身份证有一个有效就可以  

#### Version 1.19.1
> 1.文章投稿限制,添加字数字段

#### Version 1.19.0
> 1.充电留言

#### Version 1.18.1
> 1.实名认证全量

#### Version 1.18.0
> 1.实名认证上线, 灰度十分之一   

#### Version 1.17.7
> 1.添加bfs超时code

#### Version 1.17.6
> 1.去掉article无用日志，添加rpc的token

#### Version 1.17.5
> 1.修复最近充电获取用户信息panic问题

#### Version 1.17.4
> 1.修复稿件列表必须返回分区列表

#### Version 1.17.3
> 1.修复app view接口稿件充电状态返回null

#### Version 1.17.2
> 1.修复app、client稿件列表

#### Version 1.17.1
> 1.粉丝勋章增加超时时间

#### Version 1.17.0
> 1.粉丝勋章查看、开通、领取列表

#### Version 1.16.12
> 1.appeal更新稿件查看接口，不在使用archive service

#### Version 1.16.11
> 1.appeal add log

#### Version 1.16.10
> 1.更新go-business到2.28.1 (.0的identity mc坑

#### Version 1.16.9
> 1.更新memcache package

#### Version 1.16.8
> 1.更新vendor

#### Version 1.16.7
> 1.更新提交文章带草稿id

#### Version 1.16.6  
> 1.修改REST风格的请求 

#### Version 1.16.5  
> 1.稿件状态迁移到json配置文件中  

#### Version 1.16.4
> 1.bfs上传设置超时2s并限速 

#### Version 1.16.3
> 1.修复获取白名单 

#### Version 1.16.2
> 1.提交文章不过滤韩文，藏文，阿拉伯  
> 2.针对rpc服务返回的非正常errcode做转换  

#### Version 1.16.1
> 1.支持文章评论搜索查看父评论

#### Version 1.16.0
> 1.增加文章专栏管理

#### Version 1.15.6
> 1.数据中心优化 小于100关闭地区分布

#### Version 1.15.5
> 1.投稿pre接口默认添加音频码率字段192(单位kbps), 之后会有192和320区分
> 2.删除已经无用的分区配置文件

#### Version 1.15.4
> 1.升级go-common到v6.17.1和go-businessv2.24.1， 为logagent做准备  
> 2.重构Prom监控的接入方式

#### Version 1.15.3
> 1.重构评论列表中的粉丝状态显示,使用新接口异步批量查询粉丝状态

#### Version 1.15.2
> 1.数据中心增加读取前一天数据的逻辑容错

#### Version 1.15.1
> 1.打开数据中心地区排行

#### Version 1.15.0
> 1.获取模板列表接口  
> 2.添加模板接口  
> 3.更新模板接口  
> 4.删除模板接口  

#### Version 1.14.3
> 1.对应videoBycid接口重构，更新使用方式

#### Version 1.14.2
> 1.app首页改为并行请求

#### Version 1.14.1
> 1.prom在配置文件中初始化逻辑修复

#### Version 1.14.0
> 1.新增稿件列表接口

#### Version 1.13.5  
> 1.update go-common with tag v6.12.1, fix db transaction bug

#### Version 1.13.4  
> 1.prometheus监控代码植入

#### Version 1.13.3  
> 1.app最近充电用户和账单接口默认返回pager结构

#### Version 1.13.2  
> 1.重构创作中心投稿分区数据， 脱离对老库arctype的依赖

#### Version 1.13.1 
> 1.app只返回开通充电的最近充电用户列表

#### Version 1.13.0  
> 1.创作中心新分区数据结构的重构实现  

#### Version 1.12.1  
> 1.修复最近充电用户列表返回  

#### Version 1.12.0  
> 1.app首页接口显示4条最近充电用户   
> 2.app获取最近充电用户列表   
> 3.app获取电池余额和最近两个月结算记录  

#### Version 1.11.1  
> 1.稿件推荐绑定的重构实现

#### Version 1.11.0
> 1.web稿件列表追加消息通知  

#### Version 1.10.3
> 1.app 修复获取用户充电状态结构体  

#### Version 1.10.2
> 1.web/app 首页屏蔽大数据地区信息  

#### Version 1.10.1
> 1.当历史中的tag为空时使用现有的tag  

#### Version 1.10.0
> 1.web获取首页数据  
> 2.web获取工具配置位  
> 3.web获取首页聚合[运营、通知、成长之路]  
> 4.web查看生日面板显示状态  
> 5.web设置生日面板显示状态  
> 6.web最近充电排行  
> 7.web本月充电排行  
> 8.web总充电排行  
> 9.web电池余额接口  
> 10.web每日结算单接口  

#### Version 1.9.6
> 1.web/client pre添加overrate字段  

#### Version 1.9.5
> 1.重构稿件历史回溯  

#### Version 1.9.4
> 1.根据运营需求, 分区自制转载简介补充  

#### Version 1.9.3
> 1.根据运营需求, 分区自制转载简介更新  

#### Version 1.9.2
> 1.真.新分区自制转载简介更新  

#### Version 1.9.1
> 1.新分区自制转载简介更新  

#### Version 1.9.0
> 1.推荐tag V2上线  

#### Version 1.8.1
> 1.真的合并配置中心  

#### Version 1.8.0
> 1.接入新配置中心  

#### Version 1.7.0
> 1.app评论接口  
> 2.app稿件数据接口  
> 3.app视频播放退出数据接口 

#### Version 1.6.0
> 1.稿件列表改用videoup-service接口  

#### Version 1.5.2
> 1.支持windows client对接商单ID

#### Version 1.5.1
> 1.播放退出点数据算法更新  

#### Version 1.5.0
> 1.数据中心相关接口上线  

#### Version 1.4.8
> 1.升级govendor 

#### Version 1.4.7
> 1.需要把已经下架的活动的信息也加入到活动列表  
> 2.单个查询稿件会增加mission_name字段  

#### Version 1.4.6
> 1.增加搜索超时时间  

#### Version 1.4.5
> 1.恢复分区parent字段  

#### Version 1.4.4
> 1.恢复移动端的分区描述  

#### Version 1.4.3
> 1.根据运营需求，修改分区描述  

#### Version 1.4.2
> 1.用户反馈API增加超时时间  

#### Version 1.4.1
> 1.新增反馈时, 使用对方的错误码  

#### Version 1.4.0
> 1.新增用户反馈相关接口  

#### Version 1.3.0
> 1.接入了极验SDK  
> 2.新增了web稿件删除接口扣除硬币解绑活动  

#### Version 1.2.4
> 1.仅当商单id有绑定的稿件，才调用接口  

#### Version 1.2.3
> 1.商单修改了线上URL (同步修改)  

#### Version 1.2.2
> 1.去除查看稿件、投稿pre的商单依赖  

#### Version 1.2.1
> 1.更新分区名称和返回结构  

#### Version 1.2.0
> 1.商业产品部门的商单绿洲计划  

#### Version 1.1.0
> 1.web端创作中心的评论接口重构版本完成  

#### Version 1.0.6
> 1.稿件列表为空时不走RPC  
> 2.修改分区文案和打回理由文案  

#### Version 1.0.5
> 1.申诉重构接口优化  
> 1.稿件模板接口  

#### Version 1.0.4  
> 1.增加国创分区文案  

#### Version 1.0.3  
> 1.app StatePanel 增加-40  
> 2.app StatePanel 增加-11  

#### Version 1.0.2

> 1.查看稿件增加mid限制  

#### Version 1.0.1

> 1.添加稿件打回理由聚合  
> 2.封面不返回默认图  

#### Version 1.0.0

> 1.申诉功能  
> 2.移动端初版接口  
> 3.稿件相关查询接口
