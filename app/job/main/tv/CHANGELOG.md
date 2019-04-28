#### tv-job

#### Version 1.7.6
> 1. tv数据上报多表联合查询上报改为单表查询上报

#### Version 1.7.5
> 1. seasonCMS结构体新增pay_status字段，写入缓存中，用于tv-interface的卡片贴标
> 2. 消费databus更新ep缓存逻辑中新增pay_status和subtitle字段

#### Version 1.7.4
> 1. Season的Area字段由int改为string型，兼容多个地区合拍的情况；目前兼容老数据，当area仍为int时，走老逻辑，否则走新逻辑
> 2. 修复空season提审问题

#### Version 1.7.3
> 1. 同步pgc付费信息给牌照

#### Version 1.7.2
> 1. 修改sync.Waitgroup的写法，将s.waiter.Add()移出
> 2. 修改cache.New的写法，使用新的Fanout

#### Version 1.7.1
> 1. 补全dao/app下的ut

#### Version 1.7.0
> 1. 判断UGC稿件为付费稿件时，不导入tv后台
> 2. 对稿件的判断，如ugc付费，原创等全部收敛到model中，方便以后迭代

#### Version 1.6.9
> 1. 稿件自动cms上架逻辑，如果当前稿件为过审且cms下架状态，且稿件新增

#### Version 1.6.8
> 1. 补dao层中archive和ftp的ut

#### Version 1.6.7
> 1. 优化tv数据上报逻辑

#### Version 1.6.6
> 1. 修复一个cid出现在两个稿件时，其中一个cid发送到牌照方后，另外一个cid虽未发送但也显示已发送

#### Version 1.6.5
> 1. 添加详情页风格标签及跳转

#### Version 1.6.4
> 1. UGC提审加速：
- videoDatabus中判断是否为转码完成并且需要提交的，如是则加入channel中
- fullRefresh中的aid遍历逻辑也进由channel处理
> 2. 写法优化：原创逻辑、可提审逻辑等迁入model中，方便复用

#### Version 1.6.3
> 1. 为准备tv会员，临时对所有无过审或免费的season做下架处理，反之做上架处理；上下架操作结果发送到企业微信
> 2. 在配置中设置开关，tv会员上下后，只对无过审的season做下架，不再考虑付费状态

#### Version 1.6.2
> 1. 修复数据上报bug

#### Version 1.6.1
> 1. databus监听逻辑优化，比较媒资字段是否修改后再更新
> 2. 优化PGC过审保持逻辑
- 去掉所有带7的逻辑
- 改为由databus更新放入内存中，另一侧进行消费
- 提审由ep改为season，异步改同步
> 3. PGC提审优化：
- 第一次提审失败后，即进入resub队列，最多重试5次
- 去掉audit_time相关的操作
> 4. PGC提审新增字段：出品方和剧集类型

#### Version 1.6.0
> 1. tv数据上报

#### Version 1.5.9
> 1. tv-job做grpc改造
> 2. 一些log.error改成log.warn

#### Version 1.5.8
> 1. 筛选archive-Notify T消息时，需要大量查询DB中upper表，优化为查询内存中的map
> 2. pgc的全量缓存刷新由异步改为同步，使用type func统一season和ep的刷新方式，减少代码量
> 3. call视频云改为异步批量操作，提升效率，降低DB依赖
> 4. UGC的视频全量扫描，进行以下任务：
- (原有) 稿件缓存刷新
- (原有) 视频缓存刷新
- (新增) 处理删除稿件，如其下仍有未删除的分p，则进行删除
- (新增) 处理cid > 12780000 并且 转码失败的分p，进行分p及稿件（如稿件下无有效分p）的删除
- (新增) 获取稿件下所有可提审的video，分片后拼接提审xml（取消原有扫描DB提审的方式，避免慢查询）
- (新增) 提审时如果命中视频云指定错误码，如10005，则分p删除和稿件（如无其他分p）删除
> 5. 修复seasonCMS的model中的json tag，保证芒果媒资中的origin_name字段可以正常从MC中获取
> 6. dao层所有rows.Next补上rows.Err

#### Version 1.5.7
> 1. 芒果媒资同步配套

#### Version 1.5.6
> 1. 转载视频从manual、auto和init三个途径禁止进入
> 2. cid <= 12780000的分p视频无需提交视频云

#### Version 1.5.5
> 1. 导入up主的每页视频中增加间隙，防止cpu、db query等过快增长触发报警
> 2. 消费databus的消息时，通过配置可限制goroutines的数量
> 3. ugc全量缓存所有配置和pgc配置分开、ugc全量缓存生产时增加间隙、取消job启时刷新全量ugc缓存
> 4. 优化选取ugc提审video和archive的sql，以避免扫描全表
> 5. 修改稿件重新提审不再依赖DB，转而由databus获取消息后放入channel中进行通知

#### Version 1.5.4
> 1. 配合tv-interface的历史记录需求，增加pgc中ep的cms信息的缓存逻辑
> 2. 简化pgc索引页逻辑

#### Version 1.5.3
> 1. service、dao、model层整理，全部改为文件夹封装
> 2. 在redis中维护ugc索引页数据：生活，科技，游戏，时尚，音乐

#### Version 1.5.2
> 1. ugc稿件提审条件中不必要视频转码完成，若其cid <= 12780000，认为已转码完成，无需提审
> 2. 完善tv-job的UT，提升ut通过率至100%

#### Version 1.5.1
> 1. 修复up主刷新时，异步回写缓存导致缓存中数据不全的问题
> 2. 修复up主face因为bfs节点导致误认为是头像更新的问题

#### Version 1.5.0
> 1. [up主管理及同步] 新增dao/upper
- 将手动提审的aid落库时，如发现其up主未在白名单内，加入到白名单内
- 同步ugc分p或者稿件时，携带其稿件的up主信息给牌照（取原名）
- 定时全量刷新tv端up主白名单信息，检测up主头像和昵称的变动，触发牌照同步
> 2. 增加UT覆盖率

#### Version 1.4.7
> 1. 新增逻辑：新增ugc稿件时，如命中pgc分区则忽略
> 2. 新增逻辑：ugc稿件分区修改同步，如命中pgc分区则删除该稿件

#### Version 1.4.6
> 1. 针对archive-notify T中的archive的稿件的封面字段进行补全

#### Version 1.4.5 
> 1. 修复UGC搜索数据源上传ftp时的字段，由id改为aid
> 2. 修改全量刷新缓存的策略，采用串行的方式，降低并发度，避免过多timeout

#### Version 1.4.4
> 1. UGC和PGC搜索数据源上传ftp
> 2. 修改ftp上传的方式，由cron改为sleep

#### Version 1.4.3
> 1. ugc提审处判断视频云转码状态，如果转码未完成，不提审

#### Version 1.4.2
> 1. 修复copyright导致ugc不提审的问题
> 2. 修复pgc删除接口通知牌照时，sign和prefix倒反的问题
> 3. 对于删除接口，如果返回为-404视频未找到，不认为是错误，更新通知状态为成功

#### Version 1.4.1
> 1. 针对当贝和芒果的最新更新第X集的需求，使用过滤方式进行计算，加入缓存中
> 2. 修复问题：在最新ep过审后，不会更新season的缓存，导致前端拿到的newestEP不准确

#### Version 1.4.0
> 1. 在自身表的变更databus消息的消费逻辑中增加，如果archive为过审状态，则缓存archive的arc和view的rpc结果，便于ugc详情页调用archive的数据
> 2. 在fullRefresh中添加ugc的video的逻辑，将ugc的video全量缓存铺满mc

#### Version 1.3.8
> 1. 修复缓存中newestOrder字段的逻辑

#### Version 1.3.7
> 1.  针对ugc稿件，定时刷新全量的cms和鉴权信息；databus刷新增量的cms和鉴权信息

#### Version 1.3.6
> 1. 去掉ping方法中的内容，避免被踢节点
> 2. 向搜索ftp上传文件时增加重试
> 3. 去掉一些重复的log.Error
> 4. 维护全量和增量的pgc和ugc的媒资信息（用于在tv-interface中吐出给当贝）
   > 扩展在mc中的pgc的媒资信息：定时全量刷新，databus增量刷新
   > 新增在mc中的ugc的媒资信息：定时全量刷新，databus增量刷新

#### Version 1.3.5
> 1.  不同步UGC转载的稿件到牌照方

#### Version 1.3.4
> 1.  在mc中的ep鉴权信息中增加watermark信息

#### Version 1.3.3
> 1. 同步新增视频至视频云

#### Version 1.3.2
> 1. 修复牌照字段错误

#### Version 1.3.1
> 1. 修复重启时调用service.Close报错的问题

#### Version 1.3.0
> 1. 对接UGC视频数据
> 2. 手动提审视频落库
> 3. 导入up主的全量历史数据
> 4. 新增数据提审！！
> 5. 数据修改，删除，通知牌照方（分p全量diff判断修改）

#### Version 1.2.10
> 1.fix xml错误的 panic

#### Version 1.2.9
> 1.添加pgc下架过审状态保持的配套逻辑，当season/ep状态为7时，通知牌照方已删除，之后恢复到过审状态。

#### Version 1.2.8
> 1. 修改配置

#### Version 1.2.7
> 1.修复文件写入不截断的问题，加入TRUNC模式

#### Version 1.2.6
> 1.定期向搜索的ftp上传过审season的title的文件
> 2. Update Bazel

#### Version 1.2.5
> 1.同步修改过的season时，检查其下的ep数量，为0时不提交

#### Version 1.2.4
> 1.修改牌照提审前缀为可配置，xds

#### Version 1.2.3
> 1.添加redis相关逻辑，在redis中为每个分区维护一个过审season id的列表
> 2.增加databus消息同步redis逻辑

#### Version 1.2.2
> 1.playurl接口添加参数qn=16, fix视频质量为16

#### Version 1.2.1
> 1.修复文件名命名，去除大小和驼峰
> 2.对接playurl接口，获取playurl后提供给牌照方

#### Version 1.2.0
> 1.增加重新提审逻辑：如果原来已过审，season的check字段，ep的state字段会在pgc更新后变为7。tv-job对于7的进行单独重新提审，而后立刻恢复上线状态。

##### Version 1.1.0
> 1. 添加全量数据MC同步 - 每日一次定时任务将DB中数据刷新到MC中
> 2, 添加异步数据MC同步 - 监听Databus数据更新MC中 - 审核状态+干预数据

##### Version 1.0.1
> 1. 添加season的zone的翻译，除了1中国，2日本之外的数字全部翻译为“其他”

##### Version 1.0.0
> 1. 初始版本：对比pgc同步过来的ep、season表和tv_content_ep表的数据差值，插入到content表，state为1（待审核）
> 2. 选取所有待审核的ep信息，调用牌照方接口提审（包含视频云临时视频url）







