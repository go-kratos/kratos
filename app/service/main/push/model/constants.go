package model

const (
	// TempTaskPrefix used to separate from the DB tasks.
	TempTaskPrefix = "t"

	// APPIDBBPhone 哔哩哔哩动画
	APPIDBBPhone = 1

	// HTTPCodeOk http response normally.
	HTTPCodeOk = 0

	// SwitchOff off.
	SwitchOff = 0
	// SwitchOn on.
	SwitchOn = 1

	// DelMiFeedback feedback 删除 (无效token删除方式)
	DelMiFeedback = 1
	// DelMiUninstalled 卸载
	DelMiUninstalled = 2

	// DefaultMessageTitle .
	DefaultMessageTitle = "哔哩哔哩消息"

	// UnknownBuild 未知build号
	UnknownBuild = 0
)

const (
	// MobiAndroid mobi_app android
	MobiAndroid = 1
	// MobiIPhone mobi_app iPhone
	MobiIPhone = 2
	// MobiIPad mobi_app iPad
	MobiIPad = 3
	// MobiAndroidComic
	MobiAndroidComic = 4
)

// task status
const (
	// TaskStatusPending 待审核
	TaskStatusPending = int8(-5)
	// TaskStatusStop 主动停止
	TaskStatusStop = int8(-4)
	// TaskStatusDelay 延期
	TaskStatusDelay = int8(-3)
	// TaskStatusExpired 过期
	TaskStatusExpired = int8(-2)
	// TaskStatusFailed 失败
	TaskStatusFailed = int8(-1)
	// TaskStatusPrepared 未开始
	TaskStatusPrepared = int8(0)
	// TaskStatusDoing 进行中
	TaskStatusDoing = int8(1)
	// TaskStatusDone 已完成
	TaskStatusDone = int8(2)
	// TaskStatusPretreatmentPrepared 等待预处理，处理完后是按平台拆成任务(token形式)
	TaskStatusPretreatmentPrepared = int8(3)
	// TaskStatusPretreatmentDoing 预处理中
	TaskStatusPretreatmentDoing = int8(4)
	// TaskStatusPretreatmentDone 预处理完成
	TaskStatusPretreatmentDone = int8(5)
	// TaskStatusPretreatmentFailed 预处理失败
	TaskStatusPretreatmentFailed = int8(6)
	// TaskStatusWaitDataPlatform 等待从数据平台获取数据
	TaskStatusWaitDataPlatform = int8(7)
)

// data platform
const (
	// DpCondStatusNoFile 没有查询到文件
	DpCondStatusNoFile = -3
	// DpCondStatusPending 待审核
	DpCondStatusPending = -2
	// DpCondStatusFailed 失败的查询
	DpCondStatusFailed = -1
	// DpCondStatusPrepared 准备提交到数据平台的查询
	DpCondStatusPrepared = 0
	// DpCondStatusSubmitting 提交中
	DpCondStatusSubmitting = 1
	// DpCondStatusSubmitted 已经提交的查询
	DpCondStatusSubmitted = 2
	// DpCondStatusPolling 轮询任务看有没有生成文件
	DpCondStatusPolling = 3
	// DpCondStatusDownloading 正在下载文件
	DpCondStatusDownloading = 4
	// DpCondStatusDone 已经完成的查询
	DpCondStatusDone = 5

	// DpTaskTypeMid mid维度查询
	DpTaskTypeMid = 1
	// DptaskTypeToken token维度查询
	DpTaskTypeToken = 2
)

const (
	// TaskTypeAll 后台全量
	TaskTypeAll = 1
	// TaskTypePart 后台批量
	TaskTypePart = 2
	// TaskTypeBusiness 业务推送
	TaskTypeBusiness = 3
	// TaskTypeTokens 批量token推送
	TaskTypeTokens = 4
	// TaskTypeMngMid 后台按mid推送
	TaskTypeMngMid = 5
	// TaskTypeMngToken 后台按token推送
	TaskTypeMngToken = 6
	// TaskTypeStrategyMid 策略层按mid推送
	TaskTypeStrategyMid = 7
	// TaskTypeDataPlatformMid 通过mid维度从数据平台获取token
	TaskTypeDataPlatformMid = 8
	// TaskTypeDataPlatformToken 通过token维度从数据平台获取token
	TaskTypeDataPlatformToken = 9
)

const (
	// LinkTypeBangumi bangumi 协议链接类型
	LinkTypeBangumi = int8(1)
	// LinkTypeVideo 视频
	LinkTypeVideo = int8(2)
	// LinkTypeLive 直播
	LinkTypeLive = int8(3)
	// LinkTypeSplist 专题页
	LinkTypeSplist = int8(4)
	// LinkTypeSearch 搜索
	LinkTypeSearch = int8(5)
	// LinkTypeAuthor 个人空间
	LinkTypeAuthor = int8(6)
	// LinkTypeBrowser 浏览器
	LinkTypeBrowser = int8(7)
	// LinkTypeVipBuy 大会员购买页
	LinkTypeVipBuy = int8(10)
	// LinkTypeCustom 自定义协议内容
	LinkTypeCustom = int8(11)
)

const (
	// 定义参考：http://syncsvn.bilibili.co/app/wiki/blob/master/Android-App-URI.md

	// SchemeBangumiSeasonIOS 番剧详情 iPhone，iPadHD 支持番剧
	SchemeBangumiSeasonIOS = "bilibili://bangumi/season/"
	// SchemeBangumiSeasonAndroid .
	SchemeBangumiSeasonAndroid = "bili:///?type=season&season_id="

	// SchemeVideoIOS 视频详情页 iPhone，iPadHD 支持视频
	SchemeVideoIOS = "bilibili://video/"
	// SchemeVideoAndroid .
	SchemeVideoAndroid = "bili:///?type=bilivideo&avid="

	// SchemeLive 直播详情页, 支持 iOS 和 Android 新协议
	SchemeLive = "bilibili://live/"
	// SchemeLiveAndroid Android 老协议
	SchemeLiveAndroid = "bili:///?type=bililive&roomid="

	// SchemeSplist 专题页 iPhone, iPadHD, Android 支持专题
	SchemeSplist = "bilibili://splist/"

	// SchemeSearchIOS 搜索 iPhone，iPadHD 支持搜索
	SchemeSearchIOS = "bilibili://search/?keyword="
	// SchemeSearchAndroid .
	SchemeSearchAndroid = "bilibili://search/"

	// SchemeAuthorIOS 个人空间 iPhone，iPadHD 支持个人空间
	SchemeAuthorIOS = "bilibili://user/"
	// SchemeAuthorAndroid .
	SchemeAuthorAndroid = "bilibili://author/"

	// SchemeBrowserIOS 指定URL iPhone，iPadHD 支持H5
	SchemeBrowserIOS = "bilibili://browser/?url="
	// SchemeBrowserAndroid .
	SchemeBrowserAndroid = "bili:///?type=weblink&url="

	// SchemeVipBuy 大会员购买页
	SchemeVipBuy = "bilibili://user_center/vip/buy/"
)
