#### 投稿API

##### Version 2.4.1
>1.重构特殊用户组查询接口,使用grpc并且进行批量查询  

##### Version 2.4.0
>1.重构商单   

##### Version 2.3.13
>1.修复联合投稿中有封禁用户的提示语错误的bug

##### Version 2.3.12
>1.标题增加newline过滤

##### Version 2.3.11
>1.添加联合投稿功能，在稿件web编辑和web新增接口中加了联合投稿验证

##### Version 2.3.10
>1.升级账号API为grpc服务,简化配置  

##### Version 2.3.9
>1.Web端:支持高频投稿的用户组 

##### Version 2.3.8
>1.兼容支持投稿客户端ipad, app.upfrom添加UpFromIpad(11),仅限ios 
>2.完善,MaxAllVsCnt的日志格式,方便查询问题  

##### Version 2.3.7
>1.fix init service    

##### Version 2.3.6
>1.fix infoc的service new和close  
>2.支持app的投票动态带入object replacement character 

##### Version 2.3.5
>1.添加支持投票，app支持传递vote_id和vote_title    
>2.支持pic_count和video_count上报到infocproc   

##### Version 2.3.4
>1.升级并精简化自定义错误码  

##### Version 2.3.3
>1.修复和基础库错误码不兼容的问题  

##### Version 2.3.2
>1.为防止客户端(pc windows端除外)，添加和编辑的时候触发敏感词，而且用户侧已经无法修改老稿件的分P简介，决定禁止添加和修改分P的简介  

##### Version 2.3.1 for APP535
>1.素材上报支持 拍摄稿件合拍+编辑器主题使用  
>2.APP端投稿的时候允许上报POI地理位置信息    
>3.删除数据上报的部分
>4.精简日志打印 

- 编辑器：切割，字段split  
- 编辑器：裁剪，字段cut  
- 编辑器：旋转，字段rotate  
- 拍摄：闪光灯，字段flashlight  
- 拍摄：倒计时，字段countdown  
- 拍摄：美颜，字段beauty  

##### Version 2.3.0
>1.UGC投稿内容付费  
1. 统一校验-是否同意协议
2. 统一校验-创作类型+定价区间+用户灰度资格
3. 统一校验-私单+商单+UGC付费是否冲突
4. 添加稿件step1：注册内容付费（接口调用）
5. 添加稿件step2：同意当前协议（ES日志添加）
6. 编辑时候：如果已经付费的稿件，未开放之前都是可以编辑的
7. 编辑时候：如果已经付费的稿件，开放之后要删除必须在60天之后
8. 编辑时候：为了兼容其他端，付费稿件如果未提交付费信息，就按照原信息覆盖ap.UgcPay  
9. 编辑时候：如果开启调价就无脑调用   

### Version 2.2.21
>1.和服务端下发相关的素材类型兼容app上报数据集合或者单个数据, string/int   

##### Version 2.2.20
>1.兼容editor.bgms为interface，允许上报int和string    

##### Version 2.2.19
>1.由于客户端5.34上报格式问题, 暂时取消bgms上报

##### Version 2.2.18
>1.封面支持第一个第三方cdn存储:acgvideo.com  

##### Version 2.2.17
>1.封面支持第一个第三方cdn存储:clouddn.com   

##### Version 2.2.16
>1.update ut for dao  

##### Version 2.2.15
>1.fixbug 解决bm bindWith 两次读取body为EOF bug

##### Version 2.2.14
>1.fixbug 解决投稿行为日志缺少关键数据bug

##### Version 2.2.13
>1.投稿行为日志记录原始提交数据

##### Version 2.2.12
>1.APP支持投稿完成之后推荐关注创作中心官方号的MID  

##### Version 2.2.11
>1.稿件字幕提交全量开放   

##### Version 2.2.10
>1.es最多存32k,投稿日志就最多记100p吧

##### Version 2.2.9
>1.支持拍摄开关的数据上报 

##### Version 2.2.8
>1.增加编辑操作时候的总分P数目的限制，落在配置文件中，自定义返回提示信息     

##### Version 2.2.7
>1.去掉不必须要的access_key和cookie的传递   

##### Version 2.2.6
>1.修复添加稿件商单检验问题

##### Version 2.2.5
>1.添加字幕提交的业务

##### Version 2.2.4
>1.编辑器使用信息上报支持2d的投稿贴纸信息, videoup_stickers  
>2.打Warn日志报告有人在稿件标题中进行了xss攻击  

##### Version 2.2.3
>1.标题增加xss过滤

##### Version 2.2.2
>1.app端添加稿件的时候支持绑定抽奖活动 
>2.service的context使用TODO，防止cancel  

##### Version 2.2.1
>1.暂时去掉简介格式化的校验信息,remove checkDescForMap 

##### Version 2.2.0 (for app 5.31)
>1.稿件提交增加新视频的素材绑定  
>2.添加新增稿件的业务来源说明 

##### Version 2.1.37
>1.添加不可见字符的unicode正则集合，U+202E(right-to-left override) 

##### Version 2.1.36
>1.按照分区校验活动可用性，只做videoall的校验 

##### Version 2.1.35
>1.投稿日志记录build和buvid

##### Version 2.1.34
>1.添加和提交的时候，所有的分P必须去重，直接提示错误信息 

##### Version 2.1.33
>1.升级ip获取的方式，使用metadata.RemoteIP

##### Version 2.1.32
>1.feature: update acts for videoall

##### Version 2.1.31
>1.feature: 升级bm组件初始化，升级authService和verifyService 

##### Version 2.1.30
>1.feature: 添加和编辑稿件的时候，全平台的如果mission_id大于零且有效，则交叉对比创作类型，转载的稿件禁止参加活动  

##### Version 2.1.29
>1.fixVideoErrorsFrom 

##### Version 2.1.28
>1.添加投稿日志

##### Version 2.1.27
>1.去掉48小时filename的校验 

##### Version 2.1.26
>1.fix bug:tag service:5xx超时错误未捕获，并且做了数据判断 

##### Version 2.1.25
>1.tag同步绑定一级、二级分区名

##### Version 2.1.24
>1.投稿时候的充电开关展示的逻辑交各客户端自己把握,APP端暂时默认开启，等后续APP升级再做自主选择 
>2.Windows客户端的tag添加校验逻辑，直接报错提示予以修改 

##### Version 2.1.23
>1.fix bug:自定义错误兼容ecode的interface Detail

##### Version 2.1.22
>1.添加web端30秒投稿过快的警告提示 
>2.允许sid在投稿的时候提交上来 

##### Version 2.1.21
>1.Edit:完结动画32，连载动画33不允许编辑自制原创类型的稿件 

##### Version 2.1.20
>1.Add:完结动画32，连载动画33不允许添加自制原创类型的稿件 

##### Version 2.1.19
>1.去掉新投稿和编辑稿件的tag同步逻辑

##### Version 2.1.18
>1.去掉APP封面的限制 

##### Version 2.1.17
>1.为APP投稿缺少封面信息增加日志提示，方便定位具体问题

##### Version 2.1.16
>1.APP添加稿件的第一个tag当做校验的tag，并获取对应的活动id，当做参加的活动ID  

##### Version 2.1.15
>1.UGC投稿入口禁止出现ASMR分区,提示:该分区不存在  

##### Version 2.1.14
>1.兼容app，不允许活动数据负数存入db 
>2.商业产品的投稿接口重构，异步化第三方业务 

##### Version 2.1.13
>1.提供UGC商单新增稿件的接口，新增加upfrom:UpFromCM(10),另外强制校验分区必须为广告分区(166) 
>2.顺便去掉支持代码模式的代码

##### Version 2.1.12
>1.为清理配置做准备 

##### Version 2.1.11
>1.缓存用户投稿偏好数据之分区信息集合  

##### Version 2.1.10
>1.全端兼容允许稿件描述为空 

##### Version 2.1.9
>1.patch:前端native对mission_id没做有效性验证，会提交-1，暂时服务端兼容，小于等于0的都当做活动不参与 

##### Version 2.1.8
>1.fix: 检测并过滤掉前端提交的分P信息为null，保护后端业务的正确性  

##### Version 2.1.7
>1.兼容修改稿件的时候只选择了tag活动，去除后面的逗号 

##### Version 2.1.6
>1.兼容活动服务不可用的情况  

##### Version 2.1.5
>1.patch:对活动mission_id和分区做强校验，不允许活动跨分区投稿和编辑 

##### Version 2.1.4
>1.fix: 只选择活动tag，也参加了活动，当活动tag被过滤完之后，应该做是否参加活动的判断（也就是说允许只选择活动tag），后面会对活动tag和活动id不匹配做校验 

##### Version 2.1.3
>1.缓存中的活动tag过滤数据使用split之后的单个tag或者活动的名字作为key 
>2.单个tag的长度限制和tag服务的一致，30个字符 

##### Version 2.1.2
>1.patch:缓存中的活动tag过滤数据都用split之后的单个tag作为key 

##### Version 2.1.1
>1.非投稿第三方服务请求异步调用 
>2.去掉代码模式的判断，web前端已经没有代码模式投稿的入口   

##### Version 2.1.0
>1.添加和编辑稿件的时候，支持dynamic字段的添加和修改，并且要去除app对这两个字段protect的保护操作  
>2.添加稿件和编辑稿件支持参加或者修改活动 
>3.支持watermark水印数据的添加和展现,APP端添加的时候可以修改 

##### Version 2.0.6
>1.fix Bug:未参加活动的情况下不允许添加活动的tag，自定义标签包含不可选的活动tag，请修改后重新提交 

##### Version 2.0.5
>1.Dynamic完善日志输出 
>2.注释去掉WebFilterServer 

##### Version 2.0.4
>1.过滤服务添加操作记录日志,自定义错误之后moni是捕获不到日志的,需要自己手动加日志记录   

##### Version 2.0.3
>1.web投稿对接敏感词过滤服务，新增独立字符串验证的过滤，并且在创建或者编辑稿件的时候，对[转载来源、标题、简介、稿件推荐语]四个字段进行批量校验 

##### Version 2.0.2
>1.为避免审核后台改动稿件分区导致格式化简介错乱，在稿件编辑的时候，直接去掉校验数据的存在与否checkDescForMap   

##### Version 2.0.1
>1.APP修改暂时不允许修改四个字段:dynamic,porder,order,mission,desc_format_id，等待接入  

##### Version 2.0.0
>1.UpSpecial请求迁移到up-service 
>2.迁移到master目录  

##### Version 1.9.9
>1.分P细节校验和提交失败的日志内容再次添加上 

##### Version 1.9.8
>1.使用account-service v7

##### Version 1.9.7
>1.禁止转载，只有在0,8,9三种upfrom值情况下才允许修改 

##### Version 1.9.6
>1.分P错误码的接口有人在master加了，补锅，equal interface补锅   

##### Version 1.9.5
>1.修改文案, "暂无"=>"-"  

##### Version 1.9.4
>1.原始的描述信息在APP端在开启格式化简介的情况下是允许为空的，在checkDescFormat的时候，取消不应该的下限校验 

##### Version 1.9.3
>1.fix bug, errcode(21052)：转载类型的投稿，转载来源添加到描述信息之后再校验会不定概率的出现超出规定限制的长度 

##### Version 1.9.2
>1.fix bug：主app允许在编辑的时候描述信息为空提交 

##### Version 1.9.1
>1.fix bug：转载的稿件在投稿的时候描述不可能是空的，自制的会参与判断 

##### Version 1.9.0
>1.支持主APP投稿 

##### Version 1.8.8
>1.去掉对statsd的依赖 

##### Version 1.8.7
>1.移除对数据库的依赖，upSpecial信息从videoup-service获取

##### Version 1.8.6
>1.修复tagBind的不完整绑定，需要事先判断是否有修改才决定请求tagAPI，防止多次编辑产生的无效请求 

##### Version 1.8.5
>1.update baldemaster for identify sec hotfix     

##### Version 1.8.4
>1.HttpServer组件升级到BladeMaster      

##### Version 1.8.3
>1.去掉语言强制匹配       

##### Version 1.8.2
>1.支持日文投稿      

##### Version 1.8.1
>1.支持忽略等级和答题的白名单     
>2.更新单元测试以符合saga的要求      

##### Version 1.8.0
>1.对接私单项目之游戏广告交易平台     

##### Version 1.7.2
>1.取消代码模式对分P的单次提交数量限制    

##### Version 1.7.1
>1.添加稿件的时候一次性最多100P  
>2.编辑稿件的时候一次性最多100P, 创作姬除外(暂时不支持编辑时候添加分P)  
>3.Tag添加去重步骤 

##### Version 1.7.0
>1.创作姬投稿添加稿件接口上线    

##### Version 1.6.4
>1.支持前端投稿使用新转码方案, 参数为upos    

##### Version 1.6.3
>1.videoup全Post接口全去掉ecode.NoLogin   

##### Version 1.6.2
>1.在查询redis中filename是否过期之前,提前检测filename是否为空 

##### Version 1.6.1
>1.创作姬编辑稿件,如果稿件在-30,-40,-1,1,0环境下只修改tag，那么直接调用tagUpbind和update稿件库的tag字段  
>2.创作姬编辑稿件添加tag过滤  

##### Version 1.6.0
>1.提供给创作姬稿件编辑的接口，暂时支持title,tag,desc,open_elec 

##### Version 1.5.23
>1.替换稿件被锁定的错误提示，当前稿件已锁定，可能处于以下状态之一(审核锁定，网警锁定,用户删除)  

##### Version 1.5.22
>1.禁止在剧集二级分区下进行添加稿件，提示"改分区不存在"  

##### Version 1.5.21
>1.账号实名制使用protobuffer的RPC3来查询用户信息  

##### Version 1.5.20
>1.FixBug: Add的时候不能用AP.Aid,默认为0，应该使用返回的aid  

##### Version 1.5.19
>1.Web和Client添加稿件的时候重新绑定tag的信息，调用tag-service的upbind  

##### Version 1.5.18
>1.Fix:代码模式下经过转换的Videos里没有Filename，但是有Cid，所以不需要检测上传的Filename是否超过了48小时的上传期限限制  

##### Version 1.5.17
>1.完善Redis filename校验的日志，有待观察和校验  

##### Version 1.5.16
>1.fix: 在创建添加稿件的时候，不需要更新Tag,因为这个时候TagService还没有aid关联到  

##### Version 1.5.15
>1.添加支持稿件的动态字段Dynamic输入 

##### Version 1.5.14
>1.Fix: 除了PC端，暂时不允许其他端修改稿件的版权信息标记位，沿用已有的值   

##### Version 1.5.13
>1.feature: 实名认证，APP修改稿件信息  
>2.优化商单Http请求的4个API  

##### Version 1.5.12
> 1.add feature: 添加稿件和修改稿件的时候，tag统一不允许为空字符串  

##### Version 1.5.11
> 1.hot fix: 手机端无法修改分P信息，所以手机端去除对分P信息的长度校验逻辑    

##### Version 1.5.10
> 1.feature: 在添加或者编辑稿件的时候提交空分P稿件，服务端强行提示对应的文案  

##### Version 1.5.9
> 1.优化: 简介格式化的数据从videoup-service的接口获取，不连数据库查询  

##### Version 1.5.8
> 1.Profile RPC接口使用创作中心自定义的账号错误码，方便后续账号查找问题  

##### Version 1.5.7
> 1.已经被打回的稿件在编辑的时候，如果参加了活动，就强制校验活动时间的有效性  

##### Version 1.5.6
> 1.投稿的时候添加特殊限制条件,[电影]分区的特殊二级分区[欧美电影,日本电影,国产电影,其他国家]不允许提交多P稿件, 错误码为:VideoupForbidMultiVideoForTypes    

##### Version 1.5.5
> 1.重构分P上传超出48小时期限的错误提示, VideoupFilenameExpired  

##### Version 1.5.4
> 1.合并大仓库  

##### Version 1.5.3
> 1.投稿编辑时候的投稿充电更新接口重构

##### Version 1.5.2
> 1.windows client 统一接入实名制

##### Version 1.5.1
> 1.充电的ArcShow接口重构,开启充电内部接口HOST:http://elec.bilibili.co

#### Version 1.5.0
> 1.格式化描述字段提交
> 2.冗余desc_format字段用来存储格式化之后的描述信息
> 3.校验对应的desc_format_id的逻辑
    a: desc_format_id的有效性校验
    b: 分区类型和格式化模板的有效性校验
    c: 创作类型和格式化模板的有效性校验

##### Version 1.4.2
> 1.完善protectFieldForEdit，特定情况下才允许修改typeID,Copyright,tag,missionID
> 2.Fix add dealTag for edit app and client
> 3.Fix add dealTag and dealElec for add client

##### Version 1.4.1
> 1.账号的_identifyInfo接口私有化

##### Version 1.4.0
> 1.私单功能上线

#### Version 1.3.3
> 1.特殊emoji字符的过滤，允许输入基本的emoji字符集合，但是高级的由于暂时表级别还不支持，所以必须手动过滤

#### Version 1.3.2
> 1.fix bug: 只有-2:打回的稿件才能允许修改分区TypeID
> 2.add protect: 只有-2:打回的稿件才能允许修改创作类型copyright

#### Version 1.3.1
> 1.实名认证手机或者身份证有一个有效就可以
> 2.pre接口身份认证自定义: 2,"请先绑定手机"; 4,"请先完成实名认证"

##### Version 1.3.0
> 1. 实名认证, 只针对web端

##### Version 1.2.29
> 1. govendor升级一波， go-business  v2.34.0， go-common v6.24.4
> 2. 对接新的bfs上传的错误代码

##### Version 1.2.28
> 1. 添加稿件的时候Must和编辑的时候一样，必须放在第一位，否则会被正则校验三国文字所影响

##### Version 1.2.26
> 1. BFS的配置全部落到配置文件，包括传输超时和传输文件的最大值
> 2. 加上HTTP限速的代码和配置文件里面的两个配置项

##### Version 1.2.25
> 1. HotFix BFS timeout 2 seconds

##### Version 1.2.24
> 1. 添加Prom的LibClient监控

##### Version 1.2.23
> 1. fix bug: 编辑稿件的时候防止活动id为负数参数传递

##### Version 1.2.22
> 1. 升级go-common v6.17.0 和 go-business v2.24.1

##### Version 1.2.21
> 1.修复活动tag剔除

##### Version 1.2.20
> 1.bfs参数微调，大小5mb以内

##### Version 1.2.19
> 1.up投稿时，处理绑定tag

##### Version 1.2.18
> 1.更新govendor
> 2.添加prom monitor代码
> 3.修改redis的配置格式

##### Version 1.2.17
> 1.删除无用代码
> 2.修复错误码透传

##### Version 1.2.16
> 1.代码模式修改为传cid

##### Version 1.2.15
> 1.更新客户端请求方式VerifyPost=>UserPost, 统一使用c.Get("mid")获取用户mid

##### Version 1.2.14
> 1.分区信息脱离dede老库,使用videoup-service的typesAPI获取分区数据

##### Version 1.2.12
> 1.app编辑接口兼容ios

##### Version 1.2.11
> 1.Update GoBusiness to v2.13.2

##### Version 1.2.10
> 1.修复活动取消的逻辑

##### Version 1.2.9
> 1.重构活动在创建和编辑稿件的时候触发的逻辑

##### Version 1.2.8
> 1.添加allowSubmit，限制filename的48h提交

##### Version 1.2.7
> 1.过滤不可见零宽字符
> 2.去掉删除48hfilename逻辑，放job中实现

##### Version 1.2.6
> 1.修正简介回车换行被过滤的bug

##### Version 1.2.5
> 1.取消allowSubmit的限制

##### Version 1.2.4
> 1.更新分P的日志格式, 修正回车换行被过滤的bug

##### Version 1.2.3
> 1.添加allowSubmit，限制filename的超时

##### Version 1.2.2
> 1.Update MonitorPing Left redisPing

##### Version 1.2.1
> 1.添加MonitorPing

##### Version 1.2.0
> 1.接入新配置中心
> 2.过滤ascii码中的不可见字符
> 3.添加编辑和新增稿件时的info日志

##### Version 1.1.1
> 1.支持windows client对接商单ID

##### Version 1.1.0
> 1.videoup对接商单ID的实现

##### Version 1.0.8
> 1.提交的时候判断活动的etime是否有效

##### Version 1.0.7
> 1.上次封面限制改为2mb

##### Version 1.0.6
> 1.修改判断dtime的逻辑

##### Version 1.0.5

> 1.转载来源描述信息拼接在一起

##### Version 1.0.4

> 1.取消对稿件状态为-40的稿件定时发布是否有效的判断

##### Version 1.0.3

> 1.投稿tag数允许12个

##### Version 1.0.2

> 1.用户编辑稿件的时候去掉对(-7)暂缓的不允许上传的限制
> 2.手机端加上在编辑的时候加上充电选项
> 3.删除form里的cover防止打印到日志

##### Version 1.0.1

> 1.先注释48h判断，等全量15天之后再打开

##### Version 1.0.0

> 1.投稿API初版
