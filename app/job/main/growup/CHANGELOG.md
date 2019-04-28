#### growup-job

##### Version 0.0.1
> 1. 初始化growup-job
> 2. 每天14点, 更新视频稿件补贴系数表;
> 3. 凌晨3点计算用户分区收入;
> 4. 凌晨4点更新音频账号表;

##### Version 0.0.2
> 1. 规范了HTTP接口
> 2. 增加每天凌晨4点更新视频账号表
> 3. 删除音频账号表每天更新操作
> 4. 修改了查询av_charge_statis的SQL语句
> 5. 分批插入数据表av_charge_ratio
> 6. 增加查询av_charge_statis的日志
> 7. 优化查询sql

##### Version 0.0.3
> 1. 暂时移除定时更新收益帐号信息任务
> 2. 修复UP主每日计算昨日分区总收入逻辑问题

##### Version 0.0.4
> 1. 更改计算mid分区收入的执行时间,由3点改为22点
> 2. 修改计算mid分区收入的HTTP接口, 可以将要计算的日期传给接口
> 3. 将过期检查任务转移至job中

##### Version 0.0.5
> 1. 更改计算UP主分区昨日收入时间,由每天22点计算昨日收入改为每天0点计算前日收入

##### Version 0.0.6
> 1. 每日17点查询私单和绿洲商单，将最新的记录添加到黑名单
> 2. 修改获取私单的地址, appkey, secret
> 3. 更改统计稿件补贴系数查询数据库表, 由av_charge_statis 改为av_daily_charge_xx
> 4. 更新1-31计算数据
> 5. 添加数据库事务支持
> 6. 暂时删除绿洲商单加入稿件黑名单
> 7. 修改up_account

##### Version 0.0.7
> 1. 添加每日发送昨日收入邮件功能

##### Version 0.0.8
> 1. 修改邮件正文中昨日收入总金额计算

##### Version 0.0.9
> 1. 将每日稿件补贴系数计算和每日UP主分区收入计算的参数写到配置文件中
> 2. 每日17点查询绿洲商单，并添加到黑名单

##### Version 0.1.0
> 1. 将计算稿件补贴系数的冗余info log去掉

##### Version 0.1.1
> 1. 将发送邮件操作由同步改为异步
> 2. 修改up_info_video signed_at

##### Version 0.1.2
> 1. 修改每日执行稿件补贴系数每次查找数量, 将参数写在配置文件中
> 2. 优化黑名单，黑名单添加是否签约和昵称字段
> 3. 修改up_info_video mid = 365517 signed_at
> 4. 添加计算稿件每月收入接口
> 5. 修改商单表名

##### Version 0.3.2
> 1. 修复up_account老数据

##### Version 0.3.3
> 1. 添加运营标签活动ID计算

##### Version 0.3.4
> 1. 添加每月up主和稿件统计接口

##### Version 0.3.5
> 1. 优化up主和稿件统计接口

##### Version 0.3.6
> 1. 添加up主和稿件统计接口日志

##### Version 0.3.7
> 1. 修复up主和稿件统计接口bug

##### Version 0.3.8
> 1. 修改根据活动ID取稿件信息分页

##### Version 0.3.9
> 1. 修复up_account数据

##### Version 0.4.0
> 1. 修复每天删除activity_info表
> 2. 添加每日结算数据自动邮件通知

##### Version 0.4.1
> 1. 优化黑名单，黑名单添加mid
> 2. 优化定时器

##### Version 0.4.2
> 1. 修复获取稿件黑名单mid的bug

##### Version 0.4.3
> 1. 修改标签每日收入计算, 存储到不同的表中

##### Version 0.4.4
> 1. 修改获取up_tag_income 稿件累计收入事务Tx为context

##### Version 0.4.5
> 1. 删除每日定时计算标签收入功能

##### Version 0.4.6
> 1. 重构标签功能

##### Version 0.4.7
> 1. 优化标签代码和逻辑

##### Version 0.4.8
> 1. 优化av_daily_charge sql语句

##### Version 0.4.9
> 1. 标签添加重要log
> 2. 添加标签up_tag_income删除接口

##### Version 0.4.10
> 1. 标签添加重要log

##### Version 0.4.11
> 1. 标签修复获取av_daily_charge bug
> 2. 添加标签收入初始化接口

##### Version 0.4.12
> 1. 标签优化数据库查询

##### Version 0.4.13
> 1. 标签修复up_tag_income查询的bug

##### Version 0.4.14
> 1. 标签优化up_tag_income数据

##### Version 0.4.15
> 1. 修复标签计算错误数据

##### Version 0.4.16
> 1. 修复标签计算错误数据

##### Version 0.4.17
> 1. 修复标签计算错误数据

##### Version 0.4.18
> 1. 修复标签计算错误数据

##### Version 0.4.19
> 1. 添加标签单个计算接口

##### Version 0.4.20
> 1. 添加标签log

##### Version 0.4.21
> 1. 修复标签更新没有返回错误的bug

##### Version 0.4.22
> 1. 修复标签bug
> 2. 删除无用代码
> 3. 重写反作弊代码

##### Version 0.4.23
> 1. 增加up主补贴统计
> 2. 添加每日插入activity_info

##### Version 0.4.24
> 1. 修复每日插入activity_info tagid错误的bug

##### Version 0.4.25
> 1. 计算每月up收入添加提现后计算逻辑

##### Version 0.4.26
> 1. 修复每月结算后部分up主date_version没有更新到最新月的数据

##### Version 0.4.28
> 1. 为基础收入准备数据
> 2. 添加tag计算收入的定时任务
> 3. 删除计算ratio, 删除activity_info表操作

##### Version 0.4.30
> 1. 优化基础收入更新
> 2. 修改基础收入更新
> 3. 标签升级

##### Version 0.4.31
> 1. 回放历史基础收入
> 2. 标签升级

##### Version 0.4.32
> 1. 添加新标签定时任务
> 2. 删除service每日13点运行稿件标签ratio计算
> 3. 修复反作弊历史问题

##### Version 0.4.33
> 1. 修改每日统计标签收入邮件, 标签收入计算逻辑, 每日20: 00定时发送邮件
> 2. 标签收入定时任务改为每天19.30

##### Version 0.4.34
> 1. 修正基础收入数据
> 2. 添加测试代码

##### Version 0.4.35
> 1. 修改router
> 2. 添加测试代码
> 3. 合并每日12: 00发送的4封数据统计邮件为一封邮件

##### Version 0.4.39
> 1. 修复反作弊数据同步bug
> 2. 修正邮件内容中统计日期

##### Version 0.4.40
> 1. 已签约UP统计邮件添加已申请UP, 已申请UP日增长率字段

##### Version 0.4.41
> 1. 修复标签数据bug

##### Version 0.4.42
> 1. 修复计算数据bug

##### Version 0.4.43
> 1. 删除up主统计旧逻辑

##### Version 0.5.1
> 1. 创作激励新计算删除无用的代码

##### Version 0.5.2
> 1. 添加创作激励计算定时任务

##### Version 0.5.3
> 1. 清理无用代码

##### Version 0.5.4
> 1. 添加creativelog

##### Version 0.5.5
> 1. 添加creative log

##### Version 0.5.6
> 1. 添加安全验证，修复定时任务导致的bug

##### Version 0.5.7
> 1. 删除无用的日志

##### Version 0.5.8
> 1. up_tag_income 添加up主维度统计

##### Version 0.5.9
> 1. 添加创作激励补贴计算

##### Version 0.6.0
> 1.  计算拆分成补贴和收入，17点计算补贴数据，18点计算收入数据
> 2. 修复up_info_video数据

##### Version 0.6.1
> 1. 增加up主信息每日更新

##### Version 0.6.2
> 1. 修复标签数据
> 2. 修复标签上传时间为null的数据

##### Version 0.6.4
> 1. 修改标签调节

##### Version 0.6.5
> 1. 修改标签调节

##### Version 0.6.6
> 1. 创作激励半年账单接口
> 2. 创作激励预算管理接口
> 3. 修复反作弊用户更新隐患
> 4. 添加权限更新接口

##### Version 0.6.7
> 1. 标签每日计算添加防御

##### Version 0.6.8
> 1. 添加活动模版每日计算
> 2. 修复错误数据

##### Version 0.6.9
> 1. 添加活动模版日志

##### Verison 0.7.0
> 1. 修复帐号类型错误

##### Version 0.7.1
> 1. 修改活动模版获取稿件播放量

##### Version 0.7.2
> 1.  添加活动模版定时任务

##### Version 0.7.3
> 1.  修复活动模版不需要报名时没有昵称的bug

##### Version 0.7.4
> 1.  活动模版添加log

##### Version 0.7.5
> 1. 修复活动模版bug

##### Version 0.7.6
> 1. 活动模版删除无用的log

##### Version 0.7.7
> 1. 修复up_info_video少数据的bug

##### Version 0.7.8
> 1. 修复up_info_video数据的bug

##### Version 0.7.9
> 1. 添加活动模版日志

##### Version 0.8.0
> 1. 优化活动模版数据表插入

##### Version 0.8.1
> 1. 优化活动模版数据表插入去除事务

##### Version 0.8.2
> 1. 修复活动模版数据表插入错误

##### Version 0.8.3
> 1. 修复up_account数据

##### Version 0.8.4
> 1. 修复up_income数据

##### Version 0.8.6
> 1. 修复base_income数据

##### Version 0.8.7
> 1. 优化数据库批量插入
> 2. 优化创作激励预算计算
> 3. pgc用户从video表导入到column表

##### Version 0.8.8
> 1. 同步业务基础收入

##### Version 0.8.9
> 1. 增加专栏计算

##### Version 0.9.0
> 1. 增加专栏定期检查
> 2. 重新计算up_income

##### Version 0.9.1
> 1. 优化计算

##### Version 0.9.2
> 1. 修复违规扣除多扣数据

##### Version 0.9.3
> 1. 优化up_account

##### Version 0.9.4
> 1. 修复up_income total_income数据

##### Version 0.9.5
> 1. 修复up_income total_income数据

##### Version 0.9.6
> 1. 修复up_income total_income数据

##### Version 0.9.7
> 1. 修复up_income total_income数据

##### Version 0.9.8
> 1. 优化统计计算时间

##### Version 0.9.9
> 1. 重跑统计历史数据

##### Version 0.10.0
> 1. 统计添加日志

##### Version 0.10.1
> 1. 删除无用的代码

##### Version 0.10.2
> 1. 优化数据库查询

##### Version 0.10.3
> 1. 删除无用的代码

##### Version 0.10.4
> 1. 重跑统计历史数据

##### Version 0.10.5
> 1. 优化数据库查询

##### Version 0.10.6
> 1. 添加up_income weekly monthly

##### Version 0.10.7
> 1. 修复up_income weekly monthly av_base_income

##### Version 0.10.8
> 1. 素材计算接入

##### Version 0.10.9
> 1. 同步素材库

##### Version 0.11.0
> 1. 优化同步素材库

##### Version 0.11.1
> 1. 同步素材pgc

##### Version 0.11.2
> 1. 修复素材计算问题

##### Version 0.11.3
> 1. 去掉素材的安全检查

##### Version 0.11.4
> 1. 优化获取活动接口

##### Version 0.11.5
> 1. 修复素材统计

##### Version 0.11.6
> 1. 增加数据监控log

##### Version 0.11.7
> 1. 修复素材基础收入

##### Version 0.12.0
> 1. 修复素材计算bug

##### Version 0.12.1
> 1. 添加素材收入统计

##### Version 0.12.3
> 1. 优化统计计算

##### Version 0.12.4
> 1. up主视频半年账单计算

##### Version 0.12.5
> 1. up主视频半年账单数据按最后一次加入时间计算

##### Version 0.12.6
> 1. up主视频半年账单优化

##### Version 0.12.7
> 1. 素材专栏补贴管理

##### Version 0.12.8
> 1. 修复特殊字符无法插入数据的bug

##### Version 0.12.9
> 1. 补贴管理增加安全检测
> 2. 添加专栏和素材预算计算

##### Version 0.13.0
> 1. 修复标签投稿时间最后一天未计算的bug

##### Version 0.13.1
> 1. 修复标签投稿时间最后一天未计算的数据

##### Version 0.13.2
> 1.修复因标题引起的插入错误

##### Version 0.13.3
> 1.修改半年账单文案

##### Version 0.13.4
> 1. 修复素材计算错误

##### Version 0.13.5
> 1. 删除预算错误数据

##### Version 0.13.6
> 1. 修复预算未计算只有固定调节收入的bug

##### Version 0.13.7
> 1. 修复预算未计算只有固定调节收入的bug

##### Version 0.15.2
> 1. 计算调节改版
> 2. 增加task_status控制
> 3. 标签添加专栏和素材

##### Version 0.15.3
> 1. up weekly monthly添加业务total_income

##### Version 0.15.4
> 1. 优化sql查询

##### Version 0.15.5
> 1. 优化补贴计算

##### Version 0.16.0
> 1. 优化标签计算

##### Version 0.16.1
> 1. 添加自动扣除、自动处罚、自动过审

##### Version 0.16.2
> 1. 优化日志

##### Version 0.16.3
> 1. 优化自动扣除日志

##### Version 0.16.4
> 1. 优化自动过审
> 2. 优化专栏补贴计算

##### Version 0.16.5
> 1. 删除无用代码
> 2. 优化自动扣除邮件输出

##### Version 0.16.6
> 1. 优化每日获取素材接口

##### Version 0.16.7
> 1. 自动扣除，惩罚，过审接口收件人可配置

##### Version 0.16.8
> 1. 自动惩罚添加白名单，下次惩罚时不处理

##### Version 0.16.9
> 1. 修复标签调节bug
> 2. 补充unit test

##### Version 0.17.1
> 1. up主半年账单批量删除

##### Version 0.17.2
> 1. 修正配置

##### Version 0.17.3
> 1. 自动惩罚去除白名单

##### Version 0.17.4
> 1. 优化自动扣除: 当天状态改为原创后不扣除

##### Version 0.18.1
> 1. 修复自动扣除无法第二天无法恢复的bug

##### Version 0.18.2
> 1. 添加动态抽奖

##### Version 0.18.3
> 1. 修改动态抽奖host
> 2. 添加bgm白名单
> 3. 修改更新up_account逻辑减去兑换收入

##### Version 0.18.4
> 1. 优化sync.WaitGroup
> 2. 优化porder配置

##### Version 0.19.1
> 1. 优化host配置

##### Version 0.19.2
> 1. 优化邮件告警

##### Version 0.19.3
> 1. 优化计算任务安全体系

##### Version 0.19.4
> 1. 修改动态转发抽奖api

##### Version 0.19.5
> 1. wrap background task execution

##### Version 0.19.6
> 1. 修改收入去泡沫逻辑

##### Version 0.19.8
> 1. 添加一周年标签计算接口

##### Version 0.19.9
> 1. 优化一周年标签计算接口

##### Version 0.20.0
> 1. 收入计算时检查数据平台正确性

##### Version 0.20.1
> 1. 一周年同步up_account数据

##### Version 0.20.2
> 1. 统计数据写降速

##### Version 0.20.3
> 1. 统计前置任务失败添加邮件

