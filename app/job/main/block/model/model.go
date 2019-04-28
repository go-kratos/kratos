package model

import (
	"time"
)

// BlockStatus 封禁状态 0. 未封禁 1. 永久封禁 2. 限时封禁
type BlockStatus uint8

const (
	// BlockStatusFalse 未封禁
	BlockStatusFalse BlockStatus = iota
	// BlockStatusForever 永久封禁
	BlockStatusForever
	// BlockStatusLimit 限时封禁
	BlockStatusLimit
	// BlockStatusCredit 小黑屋封禁
	BlockStatusCredit
)

// BlockSource 封禁来源 1. 小黑屋(小黑屋和manager后台封禁) 2. 系统封禁(反作弊及监控系统上报) 3.解封 (所有后台，用户前台自助的解封)
type BlockSource uint8

const (
	// BlockSourceCredit 小黑屋封禁
	BlockSourceCredit BlockSource = iota + 1
	// BlockSourceSys 系统封禁
	BlockSourceSys
	// BlockSourceRemove 解封
	BlockSourceRemove
)

// String .
func (b BlockSource) String() string {
	switch b {
	case BlockSourceCredit:
		return "小黑屋封禁"
	case BlockSourceSys:
		return "系统封禁"
	case BlockSourceRemove:
		return "解封"
	default:
		return ""
	}
}

const (
	// BlockLogBizID 用户审核日志
	BlockLogBizID = 122
	// BlockJOBManagerID blockjob manager uid
	BlockJOBManagerID = 1493
	// BlockJOBManagerName blockjob manager name
	BlockJOBManagerName = "blockjob"
)

// BlockArea 封禁业务
type BlockArea uint8

// const .
const (
	BlockAreaNone BlockArea = iota
	BlockAreaReply
	BlockAreaDanmaku
	BlockAreaMessage
	BlockAreaTag
	BlockAreaProfile // 个人资料
	BlockAreaArchive
	BlockAreaMusic
	BlockAreaArticle
	BlockAreaSpaceBanner // 空间头图
	BlockAreaDynamic     // 动态
	BlockAreaAlbum       // 相册
	BlockAreaQuickVideo  //小视频
)

// String .
func (b BlockArea) String() string {
	switch b {
	case BlockAreaReply:
		return "评论"
	case BlockAreaDanmaku:
		return "弹幕"
	case BlockAreaMessage:
		return "私信"
	case BlockAreaTag:
		return "标签"
	case BlockAreaProfile:
		return "个人资料"
	case BlockAreaArchive:
		return "投稿"
	case BlockAreaMusic:
		return "音频"
	case BlockAreaArticle:
		return "专栏"
	case BlockAreaSpaceBanner:
		return "空间头图"
	case BlockAreaDynamic:
		return "动态"
	case BlockAreaAlbum:
		return "相册"
	case BlockAreaQuickVideo:
		return "小视频"
	default:
		return ""
	}
}

// BlockAction .
type BlockAction uint8

const (
	// BlockActionLimit 限时封禁
	BlockActionLimit BlockAction = iota + 1
	// BlockActionForever 永久封禁
	BlockActionForever
	// BlockActionAdminRemove 后台解封
	BlockActionAdminRemove
	// BlockActionSelfRemove 自助解封
	BlockActionSelfRemove
)

// String .
func (b BlockAction) String() string {
	switch b {
	case BlockActionLimit:
		return "限时封禁"
	case BlockActionForever:
		return "永久封禁"
	case BlockActionAdminRemove:
		return "后台解封"
	case BlockActionSelfRemove:
		return "自动解封"
	default:
		return ""
	}
}

// BlockInfo 封禁信息
type BlockInfo struct {
	MID         int64       `json:"mid"`
	BlockStatus BlockStatus `json:"block_status"` // blockStatus 封禁状态 0. 未封禁 1. 永久封禁 2. 限时封禁
	// EndTime     int64       `json:"end_time"`     // 结束封禁时间 unix time 永久封禁为0
}

// ParseDB .
func (b *BlockInfo) ParseDB(data *DBUser) {
	b.MID = data.MID
	b.BlockStatus = data.Status
}

// ParseMC .
func (b *BlockInfo) ParseMC(data *MCBlockInfo, mid int64) {
	b.MID = mid
	b.BlockStatus = data.BlockStatus
}

// BlockHistory 封禁历史
type BlockHistory struct {
	Area       BlockArea   `json:"type"`
	Operator   string      `json:"operator"`    // 操作人
	Reason     string      `json:"reason"`      // 封禁原因
	Action     BlockAction `json:"action"`      // 操作类型
	ActionTime int64       `json:"action_time"` // 操作时间
	RemoveTime int64       `json:"remove_time"` // 解封时间
	Comment    string      `json:"comment"`
}

// ParseDB .
func (b *BlockHistory) ParseDB(data *DBHistory) {
	b.Area = data.Area
	b.Operator = data.AdminName
	b.Reason = data.Reason
	b.Action = data.Action
	b.ActionTime = data.StartTime.Unix()
	b.RemoveTime = data.StartTime.Add(time.Second * time.Duration(data.Duration)).Unix()
	b.Comment = data.Comment
}

// BlockMessage 通知消息体
type BlockMessage struct {
	MID    int64       `json:"mid"`    // 用户mid
	Area   BlockArea   `json:"area"`   // BlockArea 封禁类型 1. 小黑屋(小黑屋和manager后台封禁) 2. 系统封禁(反作弊及监控系统上报) 3.解封 (所有后台，用户前台自助的解封)
	Status BlockStatus `json:"status"` // blockStatus 封禁状态 0. 未封禁 1. 永久封禁 2. 限时封禁
}
