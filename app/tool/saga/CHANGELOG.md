# v5.23.3
1. 解决企业微信发送失败问题，即部门ID号更改导致wechat接口获取用户为空引起。

# v5.23.2
1. 解决某些情况下saga提示合并成功，但实际未合并成功的问题。
2. 解决pipeline hook时的panic错误。

# v5.23.1
1. 修复saga panic错误。

# v5.23.0
1. 增加合并时生成的commit信息中，包括MR的title和描述。
2. 优化权限信息获取方式。

# v5.22.9
1. 解决pipeline状态改变发送通知时偶发的panic错误。

# v5.22.8
1. 增加超级权限用户的定制功能。
2. 增加等待合并或者正在执行合并的MR数量提示。

# v5.22.7
1. 解决saga发送通知时偶发的panic错误。

# v5.22.6
1. 解决 saga 偶尔发生panic导致caster实例重启问题。

# v5.22.5
1. 增加权限文件删除、移动时，对应的权限信息变更情况。

# v5.22.4
1. 临时恢复目标分支和配置的正则表达式分支判断失误的问题。

# v5.22.3
1. 解决目标分支和配置的正则表达式分支判断失误的问题。
2. 解决使用了目标分支直接作为权限分支的问题。

# v5.22.2
1. 增加SAGA的UT代码。

# v5.22.1
1. 解决delay merge功能出现的问题


# v5.22.0
1，增加配置权限仅限于当前目录的定制（即权限约束不再向下递归）
2，增加label准入的定制（即发版阶段设置了label的MR才允许合入）
3，增加sven平台配置更改立即生效以及配置更改后的权限信息自动同步
4，增加delay合并功能的定制（即+mr后，等pipeline跑过后自动合入，并且不需要retry pipeline）
5，调整了部分代码结构

# v5.21.1
1. 增加hbase存储时的容错

# v5.21.0
1. 增加hbase存储
2. 增加权限关联分支
3. 增加pipeline是否关联saga流程的定制
4. 增加最少review人数的定制
5. 优化代码结构和流程

# v5.20.7
1. 解决role info重复显示的问题

# v5.20.6
1. 优化role info的显示，去除all的todo

# v5.20.5
1. 解决role info显示的格式问题

# v5.20.4
1. 增加role info中owner的显示
2. 增加target分支不在白名单中的提示

# v5.20.3
1. 更改路由从V1到V2

# v5.20.2
1. +mr支持读取+1
2. 优化更新权限接口
3. 去掉多余的更新权限方法
4. 执行失败的时候立即释放锁

# v5.20.1
1. 过滤重复+merge
2. 修正note的显示问题
3. 修正owners检查的逻辑
4. 修改灰度指令+ok/+mr

# v5.20.0
1. 去掉build+lint
2. 增加retry机制
3. saga支持多实例
4. 去saga本地git操作，全改为api操作
5. saga报告改到pipeline里执行
6. 统一为redis(之前微信使用到的mc暂时未去掉，后期saga-admin封装好接口后再去掉)
7. assiged通知暂时屏蔽，后续找到好的技术方案再加进去
8. 增加友好提示信息

# v5.19.9
1. retry机制改为webhook实现
2. changelog、Swagger改到pipeline执行

# v5.19.8
1. 增加retry机制
2. 去掉build+lint

# v5.19.7
1. +merge之前判断pipeline是否通过
2. 通知根据pipeline状态是否改变来发

# v5.19.6
1. skip to audit the non-exist repo and print error log

# v5.19.5
1. add pipeline notification for all repository

# v5.19.4
1. fix swagger check bug again

# v5.19.3
1. fix taskchain
2. fix swagger check bug

# v5.19.2
1. 支持Pipeline失败通知

# v5.19.1
1. 支持自动同步企业微信名单

# v5.19.0
1. 增加 swagger 规则检查

# v5.18.9
1. 获取需要添加的企业微信名单并保存在缓存中，然后定期将名单发送给指定的人

# v5.18.8
1. 修复热更update和handle mr时，死锁问题

# v5.18.7
1. 修复reload repo时，空指针异常

# v5.18.6
1. saga对webhook自主管理，在文件中配置需要audit的webhook
2. 优化config load错误时，错误日志

# v5.18.5
1. repo 配置支持ignore filelist

# v5.18.4
1. merge成功和失败时，发送企业微信通知

# v5.18.3
1. 测试pipeline里pre-merge功能

# v5.18.2
1. 修复android-v4在热更时，误判为变化

# v5.18.1
1. 监听文件改动，支持热更

# v5.18.0
1. 支持热更

# v5.17.2
1. 修复contributor.go被识别为CONTRIBUTORS.md的情况

# v5.17.1
1. 修复Reviewers和Owners为空

# v5.17.0
1. 使用gorm代替mysql
2. 新增企业微信通知接口

# v5.16.3
1. 关闭saga触发pipeline功能

# v5.16.2
1. 修复 golint 不生效

# v5.16.1
1. 修复 eslint

# v5.16.0
1. 增加 daemonSimple 防止gitlab 邮箱轰炸
2. 增加 gitlab reward emoji 作为review标志
3. 切换gitlab接口到v4

# v5.15.3 
1. 修复go build ui err 

# v5.15.2
1. 修复reset error

# v5.15.1
1. 修复go build constraints 

# v5.15.0
1. 重构鉴权系统
2. 支持repos默认配置

# v5.14.2
1. 去掉path check
2. 优化大mr ut策略

# v5.14.1
1. 支持target branches正则表达式

# v5.14.0
1. 新增MR定制化target branches功能

# v5.13.1
1. 更新 path check
2. 修复 gitlab 适配

# v5.13.0
1. http router 切换到 bm

# v5.12.1
1. 修复新创建trigger后，空指针的panic

# v5.12.0
1. 增加path检查新部门ep
2. 将 linter 拆分为二进制版本，供gitlab ci使用

# v5.11.1
1. 修复MR未能触发gitlab pipeline的问题

# v5.11.0
1. 加入MR触发gitlab pipeline

# v5.10.2
1. 优化eslint流程

# v5.10.1
1. 优化eslint输出
2. 优化staticcheck

# v5.10.0
1. 加入php静态检查
2. 加入eslint静态检查
3. 加入 assign 通知 , review 双向通知
4. 加入path检查、解析
5. 加入changelog解析，appid、version版本
6. go vet 所有规则开放

# v5.9.3
1. 修复go build重名问题

# v5.9.2
1. 删除statsd依赖

# v5.9.1
1. 优化启动环境变量

# v5.9.0
1. 支持任意类型repo接入
2. 支持合并时，最小review数检查

# v5.8.13
1. 修复go build 作用域

# v5.8.12
1. 重构check工具
2. 改进分值计算
 
# v5.8.11
1. 增加 accpet ut 检查

# v5.8.10
1. 提升 go build 速度
2. 覆盖单元测试 build 检查 
3. 优化 task 运行日志显示

# v5.8.9
1. 修复ut selector call 

# v5.8.8
1. 放过revert分支

# v5.8.7
1. 优化ut算法

# v5.8.6
1. 修复panic

# v5.8.5
1. 修复ut在go build失败后仍然工作的bug
2. 修复覆盖率显示问题

# v5.8.4
1. 修复ut assign nil panic

# v5.8.3
1. 修复ut assign bug.

# v5.8.2
1. 修复ut assign bug

# v5.8.1
1. 修复ut检查错误

# v5.8.0
1. 增加静态单元测试覆盖率检查
2. 修复兼容xxx_test的pkg命名的单元测试 
3. 更详细和友好的提示
4. health检查定时任务

# v5.7.3
1. 修复merge没有检查unittest的错误
2. unittest兼容xxx_test的pkg命名

# v5.7.2
1. 修复conf

# v5.7.1
1. unittest纳入merge规范检查项
2. clean up code
3. 修改report

# v5.7.0
1. 对接rider retag
2. 支持rider构建时retag（+rider v1.0.0）

# v5.6.1   
1. 改进代码风格
2. 增加若干注释

# v5.6.0
1. 增加unittest检查

# v5.5.1
1. cleanup code
2. 替换merge命令

# v5.5.0
1. 重构merge鉴权，支持contributors解析的方式
2. 支持多人合作merge

# v5.4.0
1. 添加自动发布功能(+deploy [env])

# v5.3.0
1. 增加若干静态检查工具：simple,unused,gofmt,cyclo

# v5.2.0
1. 重构rider自动构建流程(+rider)
2. 接入发布api 
3. 修复若干bug

# v5.1.1
1. 修复saga diff pkg 检测算法

# v5.1.0
1. 增加任务过程展示
2. 并行化go check工具执行

# v5.0.0
1. 重构任务系统 

# v4.3.1
1. 增加目录权限白名单
2. 更新Accept MR接口

# v4.3.0
1. 支持gitlab comment hook
2. 升级gitlab新版API
3. report加入折叠功能

# v4.2.0
1. 增加CHANGELOG检查

# v4.1.0
1. 增加分支名检查，不合规的直接关闭MR

# v4.0.0
1. # vendor纳入build测试
2. 加入staticcheck
3. 定期健康检查，如发现问题邮件通知
4. 去掉 go test（未来在rider中跑测试）
5. 接入服务树

# v3.0.0
1. 项目文件变更后邮件发送
2. CONTRIBUTORS owner解析

# v2.0.0
1. 加入更多代码检查工具：go vet , golint , go test -cover
2. 更精准的Affected PKG
3. 报告内容优化
4. DAG优化，bug修复
5. DAG通过事件、周期重构
6. 增加若干log
7. 错误处理依赖github.com/pkg/errors

# v1.0.0
1. 初始化项目
