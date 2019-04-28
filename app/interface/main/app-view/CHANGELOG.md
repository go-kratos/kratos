#### App view 移动端稿件详情相关

#### Version 2.6.20

> 1.haslike错误日志

#### Version 2.6.19

> 1.三连errgroup改v2 & 点赞接口迁grpc

#### Version 2.6.18

> 1.优化点赞三连

#### Version 2.6.17

> 1.点赞过滤已点赞错误码

#### Version 2.6.16

> 1.视频详情页增加staff逻辑 

#### Version 2.6.15

> 1.相关推荐is_dalao & bangumi rating

#### Version 2.6.13

> 1.新加plat iphone_b，修改相应的版本限制
> 2.相关推荐上报id过滤字符串

#### Version 2.6.12

> 1.region表换成新表

#### Version 2.6.11

> 1.相关推荐修改上报id

#### Version 2.6.10

> 1.审核后台特殊用户组生产环境详情页接口不可见稿件

#### Version 2.6.9

> 1.拜年祭接口只返回过审的稿件

#### Version 2.6.8

> 1.不感兴趣上报buvid和mid二选一

#### Version 2.6.7

> 1.fix context

#### Version 2.6.6

> 1.高能看点和bgm

#### Version 2.6.5

> 1.广告不感兴趣上报

#### Version 2.6.4

> 1.关闭分享双写

#### Version 2.6.3

> 1.广告不感兴趣上报增加cm_reason_id

#### Version 2.6.2

> 1.拜年祭视频打开充电

#### Version 2.6.1

> 1.拜年祭接口

#### Version 2.6.0

> 1.投币接口proxy

#### Version 2.5.9

> 1.分享增加ip

#### Version 2.5.8

> 1.bnj增加tid

#### Version 2.5.7

> 1.增加嘉定proxy转发

#### Version 2.5.6

> 1.去除view接口中的upCount调用

#### Version 2.5.5

> 1.拜年祭视频可下载

#### Version 2.5.4

> 1.去掉zlimit，接入location

#### Version 2.5.3

> 1.隐藏详情页接口中的UP主投稿数字段,本周末内无问题则删除调用

#### Version 2.5.2

> 1.安卓国际版 pgc详情页

#### Version 2.5.1

> 1.接入直播的拜年祭配置

#### Version 2.5.0

> 1.拜年祭接口上线

#### Version 2.4.40

> 1.AI相关推荐ERROR特殊处理

#### Version 2.4.39

> 1.野版tv pgc详情页
> 2.搜索推荐up主卡片增加ipad

#### Version 2.4.38

> 1.修复相关推荐代码

#### Version 2.4.37

> 1.fix panic

#### Version 2.4.36

> 1.活动标签+活动tag入口+view上报增加spmid&from_spmid

#### Version 2.4.35

> 1.支持相关推荐按百分比灰度

#### Version 2.4.34

> 1.分享双写

#### Version 2.4.33

> 1.相关推荐展示日志trackid上报修改

#### Version 2.4.32

> 1.点赞动效文案修改

#### Version 2.4.31

> 1.付费稿件up主修复

#### Version 2.4.30

> 1.view接口增加fnver、fnval

#### Version 2.4.29

> 1.搜索UP主卡片返回21个

#### Version 2.4.28

> 1.相关推荐家长模式
> 2.投币且点赞的点赞上报act变更为acttolike

#### Version 2.4.27

> 1.watch config
> 1.版本限制 533以及没qn的走480p

#### Version 2.4.26

> 1.详情页点踩上报up主mid
> 2.MonitorInfo请求修改

#### Version 2.4.24

> 1.详情页点赞上报up主mid

#### Version 2.4.23

> 1.相关推荐上报加source_page,修改user_feature

#### Version 2.4.22

> 1.相关推荐秒开

#### Version 2.4.21

> 1.关注完成后的推荐关注

#### Version 2.4.20

> 1.三连进默认收藏前先判断收藏状态

#### Version 2.4.19

> 1.完善三连推荐引导条件

#### Version 2.4.18

> 1.点赞动效+点赞三连+投币同时点赞

#### Version 2.4.17

> 1.相关推荐增加pgc卡片

#### Version 2.4.16

> 1.推荐列表增加trackid

#### Version 2.4.15

> 1.fix follow panic

#### Version 2.4.14

> 1.音频换回老接口

#### Version 2.4.13

> 1.音频换新接口

#### Version 2.4.12

> 1.http请求修改

#### Version 2.4.11

> 1.音频换新接口

#### Version 2.4.10

> 1.音频换回老接口

#### Version 2.4.9

> 1.pgc迁移新接口（bangumi.bilibili.co/api/inner/pgc ---> api.bilibili.co/pgc/internal/season/appview）

#### Version 2.4.8

> 1.是否是弹幕广告字段
> 2.商单监控上报

#### Version 2.4.7

1.播放器关注引导粉丝数限制改为从配置读取

#### Version 2.4.6

1.增加播放器关注引导接口

#### Version 2.4.5

> 1.相关推荐游戏修复上报来源

#### Version 2.4.4

> 1.详情页直播间返回uri

#### Version 2.4.3

> 1.相关推荐过滤mid和buvid都为空的请求&删除灰aid的判断

#### Version 2.4.2

> 1.详情页广告和相关推荐并发处理

#### Version 2.4.1

> 1.上报av_feature格式改为json

#### Version 2.4.0

> 1.删除详情页中的elec调用

#### Version 2.3.3

> 1.视频广告去除ipad

#### Version 2.3.2

> 1.是否是弹幕广告字段

#### Version 2.3.1

> 1.给新tv增加访问innerPGC权限

#### Version 2.3.0

> 1.升级grpc identify

#### Version 2.2.15

> 1.蒙版弹幕透传

#### Version 2.2.14

> 1.加infoc日志

#### Version 2.2.13

> 1.更改errorlog日志

#### Version 2.2.12

> 1.相关推荐接口日志上报中增加goto、from信息
> 2.相关推荐走ai新接口，上报新增source av_feature build return_code user_feature

#### Version 2.2.11

> 1.ipadhd pgc build limit

#### Version 2.2.10

> 1.商单游戏
> 2.视频详情页请求商业产品数据字段增加
> 3.视频详情页相同推荐视频过滤

#### Version 2.2.9

> 1.fix avHandler

#### Version 2.2.8

> 1.异步获取直播状态超时调整

#### Version 2.2.4

> 1.异步获取直播状态调整超时时间

#### Version 2.2.3

> 1.修复广告位判断

#### Version 2.2.2

> 1.详情页接口增加自动播放的上报

#### Version 2.2.1

> 1.详情页广告接口聚合调用

#### Version 2.2.0

> 1.详情页分批长宽比
> 2.详情页相关推荐长宽比
> 3.FillURI重构

#### Version 2.1.3

> 1.干掉踩的总数

#### Version 2.1.2

> 1.详情页增加用户是否已经投币
> 2.稿件不存在返回10008

#### Version 2.1.1

> 1./vip/playurl增加aid和cid判断
> 2.部分接口indentify从user改为userMobile

#### Version 2.1.0

> 1.update infoc sdk

#### Version  2.0.9

> 1.分享稿件不存在返回10008

#### Version  2.0.8

> 1.ipad高版本电影走pgc

#### Version  2.0.7

> 1.播放器新增清晰度66流量浮层

#### Version  2.0.6

> 1.视频详情页点踩

#### Version  2.0.5

> 1.商单私单接口替换

#### Version  2.0.4

> 1.播放页UP主互选广告

#### Version  2.0.3

> 1.修改addCoin rpc

#### Version  2.0.2

> 1.点赞、取消点赞上报

#### Version  2.0.1

> 1.分享的时候增加未登录用户的上报以及up主id

#### Version  2.0.0

> 1.app-view http bm

#### Version  1.9.6

> 1.上报投币行为

#### Version  1.9.5

> 1.视频详情页曝光上报
> 2.私单接口切location

#### Version  1.9.4

> 1.fix view live

#### Version  1.9.3

> 1.history服务迁移目录，修改对该服务的包引用路径
> 2.使用account-service v7
> 3.删除拜年祭无用代码
> 4.上报分享行为

#### Version  1.9.2

> 1.视频详情页相关推荐广告卡片位置下放
> 2.card_index上报

#### Version  1.9.1

> 1.投币成功后上报给infoc

#### Version  1.9.0

> 1.相关推荐AI游戏卡片实验

#### Version  1.8.9

> 1.安卓概念版处理

#### Version  1.8.8

> 1.弹幕开关

#### Version  1.8.7

> 1.修复详情页电影切inner_PGC接口ipad版本

#### Version  1.8.6

> 1.详情页电影切inner_PGC接口

#### Version  1.8.5

> 1.profile to info
> 2.fix player icon

##### Version  1.8.4

> 1.view接口的playicon增加hash字段

##### Version  1.8.3

> 1.抽奖接口直接使用mid

##### Version  1.8.2

> 1.取消灰度mid
> 2.拜年祭去除协管，用户信息从内存获取
> 3.随机时间戳走配置
> 4.调整beginTime顺序到第一个，方便接口查看

##### Version  1.8.1

> 1.增加灰度mid判断

##### Version  1.8.0

> 1.相关推荐开启广告
> 2.新清晰度流量消耗

##### Version  1.7.10

> 1.优化相关推荐广告逻辑
> 2.分享增加稿件校验

##### Version  1.7.9

> 1.开启私单
> 2.add 增加稿件校验
> 3.1.去掉电影推荐

##### Version  1.7.8

> 1.点赞透传err

##### Version  1.7.7

> 1.调整拜年祭获奖名单

##### Version  1.7.6

> 1.拜年祭没开始时，增加默认中奖文案

##### Version  1.7.5

> 1.视频详情页相关推荐游戏拼接title

##### Version  1.7.4

> 1.客户端bug，服务端兼容
> 2.view iplimit

##### Version  1.7.3

> 1.视频详情页展示用户最近5个投稿

##### Version  1.7.2

> 1.ops rider

##### Version  1.7.1

> 1.视频不喜欢上报
> 2.fix view contributions
> 3.去掉充电无用字段

##### Version  1.7.0

> 1.拜年祭相关接口上线

##### Version  1.6.1

> 1.视频详情页相关视频运营位

##### Version  1.6.0

> 1.稿件详情页增加大会员活动信息
> 2.增加移动端获取token的接口，用于大会员清晰度校验

##### Version  1.5.5

> 1.相关视频缓存修复

##### Version  1.5.4

> 1.视频详情页ip过滤

##### Version  1.5.3

> 1.视频详情页点赞状态

##### Version  1.5.2

> 1.构建镜像

##### Version  1.5.1

> 1.视频详情页相关视频增加goto param uri字段
> 2.like 双写databus

##### Version  1.5.0

> 1.视频详情页点赞切新接口

##### Version  1.4.9

> 1.直播接口切换

##### Version  1.4.8

> 1.去掉旧逻辑代码

##### Version  1.4.7

> 1.player icon
> 2.去掉番剧承包

##### Version  1.4.6

> 1. from bangumi

##### Version  1.4.5

> 1. mid为0不调用card

##### Version  1.4.4

> 1. 视频详情页接movie_aid_info

##### Version  1.4.3

> 1. 视频详情页接inner/pgc
> 2. 点赞双写

##### Version  1.4.2

> 1.删除多余的 ecode.NoLogin

##### Version  1.4.1

> 1.修复专栏投币

##### Version  1.4.0

> 1.投币、收藏提示用户关注
> 2.视频详情页切换稿件pb接口
> 3.换router,去掉idnetify.Access

##### Version  1.3.12

> 1.增加广告透传字段

##### Version  1.3.11

> 1.删除所有http缓存管理接口代码

##### Version  1.3.10

> 1.确认了mc的slab问题，删除多余的缓存代码

##### Version  1.3.9

> 1.arcPBCache永不过期

##### Version  1.3.8

> 1.优化，投稿数>=20的时候才请求稿件详情

##### Version  1.3.7

> 1.修改cache/up接口逻辑，增量更新缓存的同时删除老缓存

##### Version  1.3.6

> 1.修复缓存不存在时，view报-404的问题

##### Version  1.3.5

> 1.去除业务自定义缓存，保留最原始的view缓存

##### Version  1.3.4

> 1.增加调试日志

##### Version  1.3.3

> 1.view缓存更新逻辑

##### Version  1.3.2

> 1.PB 初始化

##### Version  1.3.1

> 1.view接口全量pb

##### Version  1.3.0

> 1.开始双写archive pb缓存

##### Version  1.2.5

> 1.cache stat 接口增加Stat3缓存更新逻辑

##### Version  1.2.4

> 1.remove local cache

##### Version  1.2.3


1.dmRegion全量

##### Version  1.2.2


1.视频详情页增加投稿banner

##### Version  1.2.1


1.影视区二级分区开启充电

##### Version  1.2.0


1.未登录用户获取登录引导贴片

##### Version  1.1.0


> 1.视频详情页长简介长度限制

##### Version  1.0.13


> 1.视频详情页长简介

##### Version  1.0.12


> 1.音频版本过滤

##### Version  1.0.11


> 1.视频详情页音频
> 2.视频详情页增加up主投稿数
> 3.切新的movidID2Aid接口
> 4.族群弹幕

##### Version  1.0.10


> 1.腾讯视频外链禁止下载

##### Version  1.0.9


> 1.切换到view3接口

##### Version  1.0.8


> 1.修复投币

##### Version  1.0.7


> 1.修复用户昵称头像

##### Version  1.0.6


> 1.视频详情页接入音频

##### Version   1.0.5


> 1.修复稿件投诉接口

##### Version   1.0.4


> 1.自制稿件开放充电

##### Version   1.0.3


> 1.视频详情页缓存切新


##### Version   1.0.2


> 1.视频详情页充电

##### Version   1.0.1


> 1.修复app-view panic

##### Version   1.0.0


> 1.初始化项目
