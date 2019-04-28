tag-service tag平台服务       

#### V1.2.13
> 1.新增点赞、点踩、批量点赞点踩次数、upbind、adminbind、default bind grpc支持.

#### V1.2.12
> 1.去除grpc/v1目录

#### V1.2.11
> 1.频道国际版，频道分类新增attr字段传送.

#### V1.2.10
> 1.修复点赞、点踩缓存联动.

#### V1.2.9
> 1.修复订阅和自定义排序逻辑冲突.

#### V1.2.8
> 1.优化ResOidsByTid SQL 性能问题，添加IGNORE INDEX(ix_ctime)

#### V1.2.7
> 1.聚合资源-tag关系缓存.
> 2.修MutiResTagMap nil panic.
> 3.增加ut.

#### V1.2.6
> 1.删除tag-service频道处理代码.

#### V1.2.5
> 1.添加相似频道grpc api.
> 2.添加单个tag查询grpc api.

#### V1.2.4
> 1.修复空slice memcache getmuti问题.

#### V1.2.3
> 1.修复创建tag参数少的问题.

#### V1.2.2
> 1.频道头图&短评.

#### V1.2.1
> 1.提供频道GRPC api.
> 2.提供resource-tag 关系查询 grpc api.
> 3.提供tag信息查询 grpc api.

#### V1.2.0
> 1.删除tagResource
> 2.删除rpc resourceLog  
 
#### V1.1.46
> 1.举报接入grpc.

#### V1.1.45
> 1.举报切回gorpc，修改grpc命名.

#### V1.1.44
> 1.增加tag log .

#### V1.1.43
> 1.重整举报逻辑，修复第一举报者、已举报、管理员.

#### V1.1.42
> 1.修复频道规则拆离异常case.

#### V1.1.41
> 1.使用grpc，迁移举报新增接口到grpc.

#### V1.1.40
> 1.修复举报稿件-tag单用户重复举报.

#### V1.1.39
> 1.使用标准ut.

#### V1.1.38
> 1.稿件绑定tag增加日志 & 已删除的tag 不回源. 

#### V1.1.37
> 1.删掉相似tag http请求 

#### V1.1.36
> 1.迁移封禁和小黑屋接口到新接口

#### V1.1.35
> 1.修复频道计算规则，清除缓存问题.

#### V1.1.34
> 1.频道下沉tag-service.
> 2.显示频道命中规则.

#### V1.1.33
> 1.增加up主默认绑定一二级分区和管理员审核绑定一二级分区.

#### V1.1.32
> 1.删除不必要的log输出.   

#### V1.1.31
> 1.优化回源tag mc数据.   

#### V1.1.30
> 1.增加举报有效分和举报次数计数功能.

#### V1.1.29
> 1.删除service代码.

#### V1.1.28
> 1.使用公共配置.

### V1.1.27
> 1.迁移BM.

### V1.1.26
> 1.增加获取频道下线状态数据.

### V1.1.25
> 1.修复订阅关注数400限制.
> 2.增加获取稿件下tag数据（包含默认绑定分区tag和默认频道）.

### V1.1.24
> 1.rank_result          

### V1.1.23
> 1.ranking    

### V1.1.22
> 1.ResOidsByTid   

### V1.1.21
> 1.tag group 

### V1.1.20
> 1.WhiteUser & limitResource

### V1.1.19
> 1.add register

### V1.1.18第一举报人
> 1.fix report log action .  

### V1.1.18第一举报人
> 1.fix report log action .  

### V1.1.18第一举报人
> 1.增加举报rpc .  

### V1.1.17
> 1.增加举报rpc .  

### V1.1.16
> 1.fix insert into report data.     

### V1.1.15
> 1.fix bindtag log action.       

### V1.1.13
> 1.reportAction rpc          

### V1.1.12
> 1.user bind tag      

### V1.1.11
> 1.reource log tname       

### V1.1.10
> 1.user bind  & admin bind 

### V1.2.0
> 1.增加频道业务   

### V1.1.9
> 1.增加tag create rpc hide rpc         

### V1.1.8
> 1.RPC tag & resource bind   

### V1.1.7
> 1.fix like&hate num    

### V1.1.6
> 1.fix like&hate    

### V1.1.5
> 1.RPC like&hate    
### V1.1.4
> 1.RPC resAction    

### V1.1.3
> 1.fix mc panic    

### V1.1.2
> 1.合并model层 proto 文件支持bazel自动化构建      

### V1.1.1
> 1.使用account-service v7   
 
### V1.1.0
> 1.迁移到main 目录          

### V1.0.13
> 1.增加action 查询rpc    

### V1.0.12
> 1.增加resource log rpc api    

### V1.0.11
> 1.修改tag-service rpc AddCustomSubTag method. 

### V1.0.10
> 1.增加自定义排序功能    

### V1.0.9
> 1.增加举报日志处理状态
> 2.举报日志屏蔽admin用户操作记录    
> 3.tagname lenght 32 limit     

### V1.0.8
> 1.fix context        

### V1.0.7
> 1.fix card to info     

### V1.0.6
> 1.fix sub sort  

### V1.0.5
> 1.fix tag count,rpc report log & sub sort    

### V1.0.4
> 1.fix user sub tag    

### V1.0.3
> 1.增加rpc count cache  

### V1.0.2
> 1.增加rpc count接口   

### V1.0.1
> 1.修改rpc sub接口   

### V1.0.0
> 1.重构tag,实现arcTag HTTP&RPC 接口   
