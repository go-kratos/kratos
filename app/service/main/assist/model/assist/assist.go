package assist

import (
	accv1 "go-common/app/service/main/account/api"
	accmdl "go-common/app/service/main/account/model"
	"go-common/library/ecode"
	"go-common/library/time"
)

var (
	// ActEnum action enum
	ActEnum = map[int64]string{
		1: "delete",        // 删除
		2: "shield/hide",   // 屏蔽或隐藏
		3: "protect",       // 保护
		4: "disUser",       // 拉黑用户
		5: "dmPoolMove",    // 移动弹幕到字幕池
		6: "dmPoolIgnore",  // 忽略字幕池的弹幕
		7: "cancelDisUser", // 取消拉黑用户 reverse of disUser
		8: "silence",       // 直播禁言用户
		9: "cancelSilence", // 直播取消禁言用户 reverse of silence
	}
	// TypeEnum type enum
	TypeEnum = map[int64]string{
		1: "arc_com", // 稿件的评论
		2: "arc_dm",  // 稿件的弹幕
		3: "live",    // 直播的禁言
	}
	// IdentifyEnum map
	IdentifyEnum = map[int]error{
		1: ecode.UserIDCheckInvalidCard,
		2: ecode.UserIDCheckInvalidPhone,
	}
)

const (
	//Act Enum
	ActDelete        = 1 // 删除
	ActShieldOrHide  = 2 // 屏蔽或隐藏
	ActProtect       = 3 // 保护
	ActDisUser       = 4 // 拉黑用户
	ActDmPoolMove    = 5 // 移动弹幕到字幕池
	ActDmPoolIgnore  = 6 // 移动弹幕到字幕池
	ActCancelDisUser = 7 // 取消拉黑用户 reverse of ActDisUser
	ActSilence       = 8 // 直播禁言用户
	ActCancelSilence = 9 // 直播取消禁言用户 reverse of ActCancelSilence
	//TypeEnum
	TypeComment = 1
	TypeDm      = 2
	TypeLive    = 3
)

// Assist is Assists model.
type Assist struct {
	Mid       int64                 `json:"mid"`
	AssistMid int64                 `json:"assist_mid"`
	State     int8                  `json:"state"`
	CTime     time.Time             `json:"ctime"`
	MTime     time.Time             `json:"mtime"`
	Total     map[int8]map[int8]int `json:"total"`
}

// Log is single record for assist done
type Log struct {
	ID        int64     `json:"id"`
	Mid       int64     `json:"mid"`
	AssistMid int64     `json:"assist_mid"`
	Type      int8      `json:"type"`
	Action    int8      `json:"action"`
	SubjectID int64     `json:"subject_id"`
	ObjectID  string    `json:"object_id"`
	Detail    string    `json:"detail"`
	State     int8      `json:"state"`
	CTime     time.Time `json:"ctime"`
	MTime     time.Time `json:"-"`
}

// AssistRes is Assists model.
type AssistRes struct {
	Allow  int64 `json:"allow"`
	Assist int64 `json:"assist"`
	Count  int64 `json:"count"`
}

// Up is AssitUp model for space
type Up struct {
	Mid   int64     `json:"mid"`
	CTime time.Time `json:"-"`
}

// AssistUp is AssitUp model for space
type AssistUp struct {
	Mid            int64              `json:"mid"`
	Name           string             `json:"uname"`
	Sign           string             `json:"sign"`
	Avatar         string             `json:"face"`
	OfficialVerify accv1.OfficialInfo `json:"official_verify"`
	CTime          time.Time          `json:"-"`
	Vip            accmdl.VipInfo     `json:"vip"`
}

// Pager struct
type Pager struct {
	Pn    int64 `json:"current"`
	Ps    int64 `json:"size"`
	Total int64 `json:"total"`
}

type AssistUpsPager struct {
	Data  []*AssistUp `json:"data"`
	Pager Pager       `json:"pager"`
}

// SortUpsByCtime .
type SortUpsByCtime []*AssistUp

func (as SortUpsByCtime) Len() int { return len(as) }
func (as SortUpsByCtime) Less(i, j int) bool {
	return as[i].CTime > as[j].CTime
}
func (as SortUpsByCtime) Swap(i, j int) { as[i], as[j] = as[j], as[i] }
