package block

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

// Contain .
func (b BlockSource) Contain() bool {
	switch b {
	case BlockSourceBlackHouse, BlockSourceSys, BlockSourceManager, BlockSourceBplus:
		return true
	default:
		return false
	}
}

const (
	// BlockSourceBlackHouse 小黑屋封禁
	BlockSourceBlackHouse BlockSource = iota + 1
	// BlockSourceSys 系统封禁
	BlockSourceSys
	// BlockSourceManager 管理后台
	BlockSourceManager
	// BlockSourceBplus B+相关(动态、im、小视频)
	BlockSourceBplus
)

// String .
func (b BlockSource) String() string {
	switch b {
	case BlockSourceBlackHouse:
		return "小黑屋封禁"
	case BlockSourceSys:
		return "系统封禁"
	default:
		return ""
	}
}

const (
	// BlockLogBizID 用户审核日志
	BlockLogBizID int = 122
)

// BlockArea 封禁业务
type BlockArea uint8

// Contain .
func (b BlockArea) Contain() bool {
	switch b {
	case BlockAreaNone, BlockAreaReply, BlockAreaDanmaku, BlockAreaMessage, BlockAreaTag, BlockAreaProfile, BlockAreaArchive, BlockAreaMusic, BlockAreaArticle, BlockAreaSpaceBanner, BlockAreaDynamic, BlockAreaAlbum, BlockAreaQuickVideo:
		return true
	default:
		return false
	}
}

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

// BlockAction .
type BlockAction uint8

const (
	// BlockActionLimit 限时封禁
	BlockActionLimit BlockAction = iota + 1
	// BlockActionForever 永久封禁
	BlockActionForever
	// BlockActionAdminRemove 后台解封
	BlockActionAdminRemove
	// BlockActionSelfRemove 自动解封
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
	Nickname    string      `json:"nickname"`
	Username    string      `json:"username"` // 注册生成时不可更改的username
	Tel         string      `json:"tel"`
	TelStatus   int32       `json:"tel_status"`
	Mail        string      `json:"mail"` // 绑定的邮箱
	Level       int32       `json:"level"`
	SpyScore    int8        `json:"spy_score"`
	FigureRank  int8        `json:"figure_rank"`
	RegTime     int64       `json:"reg_time"`
	BlockStatus BlockStatus `json:"block_status"` // blockStatus 封禁状态 0. 未封禁 1. 永久封禁 2. 限时封禁
	BlockCount  int         `json:"block_count"`
}

// ParseStatus .
func (b *BlockInfo) ParseStatus(db *DBUser) {
	switch db.Status {
	case BlockStatusCredit:
		b.BlockStatus = BlockStatusLimit
	default:
		b.BlockStatus = db.Status
	}
}

// BlockHistory 封禁历史
type BlockHistory struct {
	Type       BlockMgrType `json:"type"`
	Operator   string       `json:"operator"`    // 操作人
	Reason     string       `json:"reason"`      // 封禁原因
	Action     BlockAction  `json:"action"`      // 操作类型
	ActionTime int64        `json:"action_time"` // 操作时间
	RemoveTime int64        `json:"remove_time"` // 解封时间
	Comment    string       `json:"comment"`
}

// BlockDetail blockDetail.
type BlockDetail struct {
	Status  BlockStatus     `json:"status"`
	Total   int             `json:"total"`
	History []*BlockHistory `json:"history"`
}

// ParseDB .
func (b *BlockHistory) ParseDB(data *DBHistory) {
	switch data.Action {
	case BlockActionForever, BlockActionLimit:
		switch data.Source {
		case BlockSourceSys, BlockSourceManager:
			b.Type = BlockMgrTypeSys
		default:
			b.Type = BlockMgrTypeCredit
		}
	case BlockActionSelfRemove, BlockActionAdminRemove:
		b.Type = BlockMgrTypeRemove
	}
	b.Operator = data.AdminName
	if data.Area.String() == "" {
		b.Reason = data.Reason
	} else {
		b.Reason = data.Area.String() + " - " + data.Reason
	}
	b.Action = data.Action
	b.ActionTime = data.StartTime.Unix()
	if b.Action == BlockActionLimit {
		b.RemoveTime = data.StartTime.Add(time.Second * time.Duration(data.Duration)).Unix()
	}
	b.Comment = data.Comment
}

// BlockMessage 通知消息体
type BlockMessage struct {
	MID    int64       `json:"mid"`    // 用户mid
	Area   BlockArea   `json:"area"`   // BlockArea 封禁类型 1. 小黑屋(小黑屋和manager后台封禁) 2. 系统封禁(反作弊及监控系统上报) 3.解封 (所有后台，用户前台自助的解封)
	Status BlockStatus `json:"status"` // blockStatus 封禁状态 0. 未封禁 1. 永久封禁 2. 限时封禁
}

// BlockMgrType mgr后台用
type BlockMgrType uint8

// BlockType enum
const (
	BlockMgrTypeCredit = iota + 1
	BlockMgrTypeSys
	BlockMgrTypeRemove
)

// BlockMgrSource mgr后台用
type BlockMgrSource uint8

// BlockMgrSource enum
const (
	BlockMgrSourceSys = iota + 1
	BlockMgrSourceCredit
)
