package model

import (
	"encoding/json"
	"go-common/library/time"
)

// 接口Action定义
const (
	ActionRecommend = iota
	ActionPlay
	ActionLike
	ActionCancelLike
	ActionFollow
	ActionCancelFollow
	ActionCommentAdd
	ActionCommentLike
	ActionCommentReport
	ActionFeedList
	ActionShare
	ActionDanmaku
	ActionPlayPause
	ActionPushRegister
	ActionPushSucced
	ActionPushCallback
	ActionBlack
	ActionCancelBlack
	ActionVideoSearch
	ActionUserSearch
	ActionUserUnLike
)

// App platform
const (
	PlatAndroid = iota + 1
	PlatIOS
)

const (
	// FeedListLen 为feed list中返回的数量
	FeedListLen = 10
	// SpaceListLen 空间长度
	SpaceListLen = 20
	// MaxInt64 用于最大int64
	MaxInt64 = int64(^uint64(0) >> 1)
	// BatchUserLen 批量请求用户信息时最大数量
	BatchUserLen = 50
)

const (
	//FromBILI video.from bilibili
	FromBILI = 0
	//FromBBQ video.from bbq
	FromBBQ = 1
	//FromCMS video.from cms
	FromCMS = 2
)

// FeedMark record the struct which returned to app in feed api
type FeedMark struct {
	LastSvID    int64     `json:"last_svid"`
	LastPubtime time.Time `json:"last_pubtime"`
	IsRec       bool      `json:"is_rec"`
}

// CursorValue 用于cursor的定位，这里可以当做通用结构使用，使用者自己根据需求定义cursor_id的含义
type CursorValue struct {
	CursorID   int64     `json:"cursor_id"`
	CursorTime time.Time `json:"cursor_time"`
}

//HTTPRpcRes ..
type HTTPRpcRes struct {
	Code int             `json:"code"`
	Msg  string          `json:"message"`
	Data json.RawMessage `json:"data"`
}
