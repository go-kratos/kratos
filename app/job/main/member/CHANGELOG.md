member项目的Job

# Version 2.21.3
> 1. 删除部分无用代码

# Version 2.21.0
> 1. 删除所有 user_detail_ 依赖

# Version 2.20.1
> 1. 删除 bfsDatabus配置

# Version 2.20.0
> 1. 下线 user_verify 表
> 2. 下线 member_verify 表

# Version 2.19.4
> 1. 实名认证超时自动驳回
> 2. UpdateRealnameFromMSG方法改为一个事务内操作realnameInfo和realname_apply表
> 3. 使用 Infov 上报解析身份证后的部分公开数据，修正时间格式

# Version 2.19.1
> 1. 修正经验缓存通知时序

# Version 2.19.0
> 1. 修正视频分享经验

# Version 2.18.4
> 1. 限制单节点的经验数据库写入量

# Version 2.18.3
> 1. 使用字符串来传递登录日志中的 IP 字段

# Version 2.18.2
> 1. 补发芝麻实名认证经验

# Version 2.17.0
> 1. 增加芝麻实名认证

# Version 2.16.2
> 1. fix conf redefined

# Version 2.16.1
> 1. fix remove block source check

# Version 2.16.0
> 1. realname 新表realname_info写入
> 2. realname card_md5 更新算法

# Version 2.15.0
> 1. 所有 base 修改延迟 5 秒再删一次缓存

# Version 2.14.0
> 1. 合并block

# Version 2.13.6
> 1. 经验insert binlog 删除缓存 

# Version 2.13.5
> 1. 在实名主库写入完成后增加clean cache 通知
> 2. 优化account notify

# Version 2.13.4
> 1. 删除 hbase
> 2. 删除 hbase client

# Version 2.13.2
> 1. 修复实名删除缓存策略，防止主从延迟导致脏缓存

# Version 2.13.1
> 1.add notify sender mark.

# Version 2.13.0
> 1. 增加实名删缓存兜底策略，保证最终一致

# Version 2.12.3
> 1. 增加节操变化通知业务方清理缓存

# Version 2.12.2
> 1. 修复经验等级获取

# Version 2.12.1
> 1. 增加等级变化清缓存通知

# Version 2.10.0
> 1. 观看视频补偿登陆加节操

# Version 2.9.5
> 1. 观看视频补偿登陆加节操

# Version 2.9.4
> 1. 更改realname cache key

# Version 2.9.3
> 1. 头像检查增加动态开关

# Version 2.9.2
> 1. 节操日志带 log_id

# Version 2.9.1
> 1. AI 审核头像带备注

# Version 2.9.0
> 1. 经验日志写 report

# Version 2.8.2
> 1. realname 优化

# Version 2.8.1
> 1. 修复名字全量同步bug

# Version 2.8.0
> 1. 删除 头像，签名从老库同步的逻辑
> 2. 新库头像回写老库

# Version 2.7.1
> 1. 头像自动审核添加operator

# Version 2.7.0
> 1. 去除 reload 配置文件

# Version 2.6.0
> 1. 增加头像审核相关逻辑

# Version 2.5.4
> 1. 节操恢复逻辑：由每天登录恢复一点改为恢复从上次恢复到当前时间累加

# Version 2.5.3
> 1. change recoverMoral status.

# Version 2.5.2
> 1. fix moral recover

# Version 2.5.1
> 1. 去掉节操和生日同步等代码
> 2. 增加每日登录恢复节操逻辑

# Version 2.5.0
> 1. migrate to bm

# Version 2.4.3
> 1. reissue login exp

# Version 2.4.2
> 1. del rank sync

# Version 2.4.1
> 1. fix exp init

# Version 2.4.0
> 1. fix exp init

# Version 2.3.19
> 1. fix exp init

# Version 2.3.18
> 1. del sex sync

# Version 2.3.17
> 1. 老 verify 同步到新 official 表

# Version 2.3.16
> 1. job 里不再写 base 缓存

# Version 2.3.15
> 1. 去掉老mc清理

# Version 2.3.14
> 1. 修复经验日志被不同 operator 覆盖的问题

# Version 2.3.13
> 1. 修复实名认证card_data base64解码

# Version 2.3.12
> 1. 修复实名认证card_type & country

# Version 2.3.11
> 1. 修复实名认证dao

# Version 2.3.10
> 1. 修复实名认证同步时间parse

# Version 2.3.9
> 1. 实名认证订阅老表+同步老库数据

# Version 2.3.8
> 1. 将数据库查不到，返回结果改为nil

# Version 2.3.7
> 1. 严格节操一致检查

# Version 2.3.6
> 1. 节操全量同步，如果老的里面有值新的没有就把值初始化为7000

# Version 2.3.5
> 1. 删除缓存集群由mc改为mcTmp

# Version 2.3.4
> 1. 增加删除节操缓存
> 2. 去掉detail 变化通知

# Version 2.3.3
> 1. 增加临时缓存集群清理

# Version 2.3.2
> 1. add binlog&fixer birthday sync to base.

# Version 2.3.1
> 1. super clean code.  
> 2. remove add&set exp check&init logic.  
> 3. configureable exp proc count.  

# Version 2.3.0
>1. fix moral log.

# Version 2.2.9
>1. purge cache on base update.  

# Version 2.2.8
> 1. add incr sync moral and log.  
> 2. add exp databus log.

# Version 2.2.7
> 1. sync moral.  
> 2. enlarge accproc goroutine.

# Version 2.2.6
> 1. support disposable exp.  
> 2. change unmarshal error continue.  

# Version 2.2.5
> 1. fix passportSubproc. 

# Version 2.2.4
> 1. add exp init logic.  
> 2. chenge log format.  
> 3. perfect code.

# Version 2.2.3
> 1. fix face url
> 2. nil base protect
> 3. skip on zero mid
> 4. login log message parse with two struct
> 5. base info with cert
> 6. notify purge cache and logging
> 7. restrict exp change and readable message
> 8. unify exp log field

# Version 2.2.2
> 1. add  goroutine to consume log 

# Version 2.2.1
> 1. fix sync name
> 2. add sync range
> 3. add debug log
> 4. fix sql in


# Version 2.2.0
> 1. exp reconstruct
> 2. acc proc count  
> 3. skip exp log sync  

# Version 2.1.2
> 1. clarify spacesta and silence

# Version 2.1.1
> 1. fix user_addit info

# Version 2.1.0
> 1. fix sync aso source

# Version 2.0.9
> 1. add sync aso source

# Version 2.0.8
> 1. add feature gates

# Version 2.0.7
> 1. check data in job leader
> 2. remove statsd

# Version 2.0.6
> 1. add check exp code.
> 2. change member info sync.

# Version 2.0.5
> 1. add change exp log

# Version 2.0.4
> 1. rm change sleep 

# Version 2.0.3
> 1. change sleep time

# Version 2.0.2
> 1. add time sleep

# Version 2.0.1
> 1. add log

# Version 2.0.0
> 1. check exp db to accdb  

# Version 1.9.5
> 1.unify job model  
> 2.add ctime/mtime field for base and detail

# Version 1.9.4
> 1.simplify code

# Version 1.9.3
> 1.fix len(Name)>0

# Version 1.9.2
> 1.fix map init

# Version 1.9.1
> 1.fix map key int64

# Version 1.9.0
> 1.fix map key type

# Version 1.8.9
> 1.fix un-init map

# Version 1.8.8
> 1.simplify code

# Version 1.8.7
> 1.fix infinity loop

# Version 1.8.6
> 1.fix atomic load int64

# Version 1.8.5
> 1.fix slice init

# Version 1.8.4
> 1.add batch get from db

# Version 1.8.3
> 1.add wocao juran recover le

# Version 1.8.2
> 1.adjust fixer strategy

# Version 1.8.1
> 1.add fixer detail field check

# Version 1.8.0
> 1.randomize checking mids

# Version 1.7.9
> 1.if cannot load baseinfo or detailinfo , rewrite it

# Version 1.7.8
> 1.ignore db timeout 

# Version 1.7.7
> 1.added dat fixer 

# Version 1.7.6
> 1.added missing mid  

# Version 1.7.5
> 1.passthrough sign checking  

# Version 1.7.3
> 1.check all data   

# Version 1.7.2
> 1.passthrough sex face rank checking  

# Version 1.7.1
> check all data

# Version 1.7.0
> 1.old data detail init  

# Version 1.6.0
> 1.add detail init  

# Version 1.5.0
> 1.修改exp 为mc  

# Version 1.4.0
> 1.添加award databus  

# Version 1.3.0 
> 1.增加业务方缓存清理

# Version 1.2.1 
> 1.添加登陆经验开关  

# Version 1.2.0
> 1.经验重构  

# Version 1.1.0 
> 1.经验数据迁移  

# Version 1.0.3
> 1.批量数据校验

# Version 1.0.0
> 1.版本初始化
  2.增量数据同步
  3.缓存清理

