package model

// const reply
const (
	SubTypeArchive     = int8(1)  // 稿件
	SubTypeTopic       = int8(2)  // 专题
	SubTypeDrawyoo     = int8(3)  // 画站
	SubTypeActivity    = int8(4)  // 活动
	SubTypeLive        = int8(5)  // 直播小视频
	SubTypeForbiden    = int8(6)  // 封禁信息
	SubTypeNotice      = int8(7)  // 公告信息
	SubTypeLiveAct     = int8(8)  // 直播活动
	SubTypeActArc      = int8(9)  // 主站活动稿件
	SubTypeLiveNotice  = int8(10) // 直播公告
	SubTypeLivePicture = int8(11) // 文画
	SubTypeArticle     = int8(12) // 文章
	SubTypeTicket      = int8(13) // 票务
	SubTypeMusic       = int8(14) // 音乐
	SubTypeCredit      = int8(15) // 风纪委案件
	SubTypePgcCmt      = int8(16) // pgc点评
	SubTypeDynamic     = int8(17) // 庐山动态
	SubTypePlaylist    = int8(18) // 播单
	SubTypeMusicList   = int8(19) // 音乐播单

	ReportStateNew       = int8(0) // 待审
	ReportStateDelete    = int8(1) // 移除
	ReportStateIgnore    = int8(2) // 忽略
	ReportStateDeleteOne = int8(3) // 一审移除
	ReportStateIgnoreOne = int8(4) // 一审忽略
	ReportStateDeleteTwo = int8(5) // 二审移除
	ReportStateIgnoreTwo = int8(6) // 二审忽略
	ReportStateAddJuge   = int8(8) // 移交仲裁

	AppealBusinessID = int64(13) // 举报工单业务id

	// AutoOPID auto oper_id
	AutoOPID = int64(877)

	// ReplyOriginURL reply origin url.
	ReplyOriginURL = `https://www.bilibili.com/video/av%d/#reply%d`
)
