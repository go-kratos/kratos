# credit 用户账号会员系统的服务层。
# v1.25.0
> 1.接入 account-service gRPC服务

# v1.24.3
> 1.增加勋章发放重试

# v1.24.2
> 1.增加封禁类型（恶意冒充他人）

# v1.24.1
> 1.fix block list panic

# v1.24.0
> 1.使用member-service的block的接口 

# v1.23.1
> 1.fix jury cache redelete   

# v1.23.0
> 1.使用filter-service的grpc协议   

# v1.22.1
> 1.fix 挂件的post协议   

# v1.22.0
> 1.替换老发放勋章接口  
> 2.优化代码  

# v1.21.0
> 1.改为 metadata RemoteIP

# v1.20.2
> 1.fix async var change  

# v1.20.1
> 1.去掉cache压缩  

# v1.20.0
> 1.风纪委缓存优化 
> 2.实名认证使用memrpc

# v1.19.1 
> 1.劳改是否答题增加缓存 

# v1.19.0
> 1.use new auth 

# v1.18.3
> 1.fix block time bug 

# v1.18.2
> 1.强制永久封禁的 blocked day的为0   

# v1.18.1
> 1.新增动态、相簿、小视频类封禁

# v1.18.0
> 1.答题吐稿件id + title   

# v1.17.2
> 1.fix api /blocked/info/batch/add

# v1.17.1
> 1.原来的content的check没有作用

# v1.17.0
> 1.old bind to new

# v1.16.12
> 1.产品变动小黑屋封禁理由  

# v1.16.11
> 1.add reason code (28、29)  

# v1.16.10
> 1.add reason code (27)  

# v1.16.9
> 1. 匿名去昵称    

# v1.16.8
> 1. debug 观点显示  

# v1.16.7
> 1. 观点换key  

# v1.16.6
> 1.观点列表直接取新表数据  

# v1.16.5
> 1.小黑屋未同步到封禁数据，兼容前端强制封禁次数为1   

# v1.16.4
> 1.rollback answer code  

# v1.16.3
> 1.punish_time to now  

# v1.16.2
> 1.clear old code  
> 2.新增批量封禁接口  
> 3.众议院匿名安全漏洞  

# v1.16.1
> 1.风纪委封禁信息判断  

# v1.16.0
> 1.众裁举报理由变更  

# v1.15.6
> 1.使用account-service v7  

# v1.15.5
> 1.fix async context to timeout bug  

# v1.15.4
> 1.删除 “发布色情低俗信息”  
> 2.新增 “破坏网络安全”、“发布虚假误导信息”  

# v1.15.3 
> 1. 更新“发布怂恿教唆信息”为“发布不适宜内容”  
> 2. 新增“发布青少年不良内容”  

# v1.15.2  
> 1.delete statsd  

# v1.15.1  
> 1.config fix case_give_hours  

# v1.15.0  
> 1.风纪委投票计数分级  

# v1.14.1  
> 1.修复appeal status  

# v1.14.0  
> 1.add case list punishTitle  
> 2.clear old code  

# v1.13.9  
> 1.update cache logic for blocked info  

# v1.13.8  
> 1.update cache key for blocked info and log  

# v1.13.7  
> 1.fix block info status check  

# v1.13.6  
> 1.fix tagID map nil  

# v1.13.5 
> 1.workflow tagid 从配置文件获取  

# v1.13.4 
> 1.修复 blocked sql  

# v1.13.3 
> 1.修改workfolw tagID  

# v1.13.2  
> 1.增加封禁申诉  

# v1.13.1
> 1. 新增批量获取case的info接口  

# v1.13.0
> 1. http使用BM框架  

# v1.12.4
> 1. 案件移交加入举报时间字段  

# v1.12.3
> 1. add goconvey test  

# v1.12.2
> 1. fix_bug  

# v1.12.1
> 1. 案件移交过滤已经移交过的案件  

# v1.12.0
> 1. 弱鸡的观点防刷机制  
> 2. code fix  

# v1.11.15
> 1. expland get manager user ps.  

# v1.11.14
> 1. load manager admin users  
> 2. notice order by id desc  

# v1.11.13   
> 1. 操作者封禁做后台操作者转操作ID转换  

# v1.11.12  
> 1. 优化 KPI list

# v1.11.11  
> 1. fix 排序字段bug  

# v1.11.10   
> 1. 修改blocked/list排序字段  
> 2. 格式优化部分代码   

# v1.11.9  
> 1. 封禁接口新增业务字段和逻辑判断  

# v1.11.8  
> 1. 新增专栏理由  

# v1.11.7  
> 1.fix blockedlist expire bug

# v1.11.6  
> 1.代码优化  

# v1.11.5  
> 1.增加专栏与头图描述  

# v1.11.4
> 1.老萧产品需求，封禁历史记录一年三次以上可永久封禁  

# v1.11.3
> 1.劳改题目is_del字段 

# v1.11.2
> 1.del opinion join  

# v1.11.1
> 1.reason_type fix bug  

# v1.11.0
> 1.blocked_opinion 冗余mid和vote字段  

# v1.10.10
> 1.KPI优化新增被赞数  

# v1.10.9
> 1.支持publish的status字段  
> 2.支持info的status字段  
> 2.支持publish的show_time字段  

# v1.10.8
> 1.修改publish/infos接口返回data为map  

# v1.10.7
> 1.整理业务代码  
> 2.新增查询用户封禁次数接口  
> 3.新增业务封禁  
> 4.批量查询封禁信息  
> 5.批量查询封禁历史  
> 6.批量查询公告信息  
> 7.批量查询封禁委员信息  

# v1.10.6
> 1.修改reason_type  

# v1.10.5
> 1.新增大众总裁逻辑  

# v1.10.4
> 1.修改reason_type  

# v1.10.3
> 1.修改系统消息和徽章发放为co域名接口  
> 2.整理域名  

# v1.10.2
> 1.小黑屋中获得风纪委员的资格的bug  

# v1.10.0
> 1.风纪委获取case逻辑变更  

# v1.9.2
> 1.fix blockedinfo log  

# v1.9.1
> 1.credit 3.0  

# v1.9.0
> 1.credit dao add prom  

# v1.8.9
> 1.fix caseobtain bug  

# v1.8.8
> 1.fix caseobtain bug  

# v1.8.7
> 1.fix panic nil  

# v1.8.6
> 1.fix panic nil  

# v1.8.5
> 1.风纪委 发放勋章  

# v1.8.4
> 1.删除无用类型的参数绑定的默认default  
> 2.加入答题防刷  

# v1.8.3
> 1.fix bug

# v1.8.2
> 1.fix filter add log    

# v1.8.1
> 1.fix filter bug  

# v1.8.0
> 1.迁出member-interface的风纪委的对外接口  
> 2.接入新版httpclient  
> 3.项目修改为credit

# v1.7.3
> 1.添加originType音频  

# v1.7.2
> 1.批量移交风纪委案件接口  

# v1.7.0
> 1.jury2.2 观点关键词屏蔽 观点赞-踩<=-5不显示观点 投票后展示赞踩数  

# v1.6.0
> 1.小黑板  

# v1.5.0
> 1.风纪委公开裁决  
> 1.风纪委公开裁决  

# v1.4.2
> 1.劳改 增加内网软删除接口  
> 2.劳改 增加获取题目错误码  

# v1.4.1 
> 1.添加rpc method  

# v1.4.0 
> 1.合并=进大仓库  

# v1.3.0   
> 1.风纪委众议   

# v1.2.1
> 1.fix 重复提交返回空数据  
> 2.增加永久封禁不能答题  

# v1.2.0
> 1.劳改项目  

# v1.1.1
> 1.修复sql id字段缺失  

# v1.1.0
> 1.风纪委重构  

# v1.0.0
> 1.基础api

