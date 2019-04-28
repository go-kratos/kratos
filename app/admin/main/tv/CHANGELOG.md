### tv端视频审核

##### Version 1.7.1
> 1. UGC审核查询和内容库（result=1即过审）迁移到ES，对接ES的SDK，删除原有的DB逻辑
> 2. UP主管理列表增加"干预昵称"字段、搜索昵称改为搜索UP主原名
> 3. UGC内容库增加可按UP主原名或者MID搜索稿件

##### Version 1.7.0
> 1. 新增tv admin tv会员 订单列表查询接口
> 2. 新增tv admin tv会员 价格面板创建接口
> 3. 新增tv admin tv会员 价格面板列表查询接口
> 4. 新增tv admin tv会员 价格面板根据id删除接口
> 5. 新增tv admin tv会员 价格面板根据id查询接口
> 6. 新增tv admin tv会员 用户根据mid查询接口

##### Version 1.6.2
> 1. 合拍影视的地区为多个如"中国,美国",所以将area字段由int型改为string型用于接收此类信息

##### Version 1.6.1
> 1. 允许season的total_num为0（未完结的番剧），修改其参数int型为string型

##### Version 1.6.0
> 1. 允许ep的length为0，去掉form中该字段的validate

##### Version 1.5.9
> 1. ugc付费禁止手动添加进tv库

##### Version 1.5.8
> 1. 新增索引标签排序接口，修改position字段
> 2. 索引标签列表接口排序改为使用position字段排序

##### Version 1.5.7
> 1. 修改模块发布bug

##### Version 1.5.6
> 1. 添加动态标签接口

##### Version 1.5.5
> 1. 修复新页面添加新模块的问题

##### Version 1.5.4
> 1. 修改modules的添加和编辑接口，使得其支持传入类型moretype：2（新pgc索引），3（新ugc索引）；morepage传入pgc或者ugc的一级分区id

##### Version 1.5.3
> 1. 修改模块配置页面 err为nil的bug

##### Version 1.5.2
> 1. 渠道闪屏改版需求，增加筛选条件

##### Version 1.5.1
> 1. PGC的season列表和UGC的archive列表增加主站发布时间字段
> 2. 干预列表增加主站发布时间字段（判断PGC或UGC，读取相应字段展示）
> 3. pgc的category翻译成中文的逻辑改为走配置，不再写死在代码中

##### Version 1.5.0
> 1. 新增索引页干预类型（list和publish接口），发布时验证ugc稿件或pgc剧集属于所提交的分区，否则认为失败
> 2. 对于模块页和精选页的干预进行优化，收敛差异化操作至model中。修改路由，使得三种干预操作放入同一group中（清真的味道~）

##### Version 1.4.9
> 1. 动态分区配置

##### Version 1.4.8
> 1. 添加tv版1.13的索引标签管理模块
> 2. 添加索引标签自动同步pgc cond接口及稿件的分区的机制，自动添加"全部"标签

##### Version 1.4.7
> 1. 修复ep导入时state为0的问题

##### Version 1.4.6
> 1. createSeason接口修改：
- 新增参数：version、producer
- 改用bind的required进行参数校验
- season的update使用反射进行更新字段检测，只更新有变更的字段，如无变更则不更新
> 2. online和hidden冗余代码优化
> 3. pgc过审保持逻辑优化，不再更新state=7进行过渡

##### Version 1.4.5
> 1. 修复ugc审核查询的rows未close的问题

##### Version 1.4.4
> 1. 增加ugc审核查询相关接口（由于ES暂不支持按照title搜索，所以title搜索暂时还是走DB，其他搜索走ES）
> 2. 补齐tv-admin的UT，并修正无法通过的UT
> 3. pgc审核查询接口改版：合并审核中的2种状态、增加ep和season的结果中增加ctime字段
> 4. 增加批量删除媒资接口，用于进行PGC/UGC资料的软删除
> 5. 增加ugc异常cid（已提交转码一定时间，仍未返回值的）的列表导出功能

##### Version 1.4.3
> 1. 添加pgc转码查询后台接口

##### Version 1.4.2
> 1. archive-service和account-service改为grpc
> 2. 配合archive-service的grpc，将typeid从int16改为int32

##### Version 1.4.1
> 1. 芒果推荐位接口：新增、删除、展示、编辑、发布

##### Version 1.4.0
> 1. 手动提审视频时，转载的视为无效，不允许添加

##### Version 1.3.13
> 1. 干预发布支持老的无分类的数据，无分类默认为PGC干预

##### Version 1.3.12
> 1. 两套干预接口统一为一套逻辑，新增支持UGC干预逻辑
> 2. 模块数据源支持ugc（5个pgc的category+5个ugc的一级分区及其下属的二级分区）
> 3. 增加接口，用于吐出模块化所支持的数据源

##### Version 1.3.11
> 1. 按照视频云要求，ugc playurl接口新增参数："platform" = "tvproj"

##### Version 1.3.10
> 1. UGC内容库

##### Version 1.3.9
> 1. cms中新增up主管理（干预、上下架等）

##### Version 1.3.8
> 1. ugc 视频分区过滤

##### Version 1.3.7
> 1. 模块化干预默认取100条数据

##### Version 1.3.6
> 1. 修复分页参数错误

##### Version 1.3.5
> 1. 修复权限错误

##### Version 1.3.4
> 1. 模块化

##### Version 1.3.3
> 1. 修复bm参数映射错误问题

##### Version 1.3.2
> 1. 水印管理
> 2. fix bazel file

##### Version 1.3.1
> 1. 修复pgc数据bind的问题

##### Version 1.3.0
> 1. TV四期需求——UGC相关功能：up主管理接口（添加、删除、导入历史数据等等）
> 2. 批量手动添加稿件接口

##### Version 1.2.12
> 1. 修改MC 找不到key 报500错误
> 2. 修复方法错误问题

##### Version 1.2.11
> 1. 搜索热词干预

##### Version 1.2.10
> 1. 干预列表和干预发布增加对"最新更新"，即category=3的支持
> 2. pgc下架ep或者season时，添加过审保持逻辑，保证鉴权报错准确性
> 3. 避免pgc接口更新时多次打更新，将原先过审的ep/season更新成未过审

##### Version 1.2.9
> 2. 修改默认配置

##### Version 1.2.8
> 1. Bazel Update

##### Version 1.2.7
> 1. 迁移bm框架
> 2. 审核查询接口支持manager权限点控制

##### Version 1.2.6
> 1, 修复pgc数据无法上架问题

##### Version 1.2.5
> 1, 修复season添加时的报错问题，分开searepo和season/create使用的结构体

##### Version 1.2.4
> 1.修正添加干预的时候时间排序问题
> 2.修正ep添加的时候title可以为空
> 3.修正升级管理添加版本的时候灰度更新失效问题

##### Version 1.2.3
> 1.双写内容库tv_content/tv_ep_content表, 原因interface读表和admin写表不一致有bug

##### Version 1.2.2
> 1.修正升级管理是否推送/强制推送勾选问题
> 2.修正dao层dbshow日志记录问题
> 3.修正剧集标题模糊搜索

##### Version 1.2.1
> 1.预览功能对接playurl接口，增加qn=16逻辑
> 2.修复剧集无法显示全部的问题
> 3.封装bfs接口

##### Version 1.2.0
> 1.修改route为verify route

##### Version 1.1.1
> 1.预览功能对接playurl

##### Version 1.0.0
> 1.新增审核通过列表接口
> 2.新增上下线接口

##### Version 1.1.0
> 1.新增干预列表展示，可搜索筛选。在展示时会检验所有干预的有效性，对于无效干预进行接口吐出展示，并且进行删除操作。
> 2.新增干预列表发布接口，发布后删除之前所有干预，添加新干预。
> 3. 新增审核结果查询接口：ep接口和season接口
> 4. 新增内容库/升级管理
