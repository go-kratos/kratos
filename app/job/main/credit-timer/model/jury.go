package model

import "time"

// case status
const (
	CaseStatusGranting  = 1 // 发放中
	CaseStatusGrantStop = 2 // 停止发放
	CaseStatusDealing   = 3 // 结案中
	CaseStatusDealed    = 4 // 已裁决
	CaseStatusRestart   = 5 // 待重启
	CaseStatusUndealed  = 6 // 未裁决

	JudgeTypeUndeal  = 0 // 未裁决
	JudgeTypeViolate = 1 // 违规
	JudgeTypeLegal   = 2 // 未违规

	VoteTypeUndo    = 0 // 未投票
	VoteTypeViolate = 1 // 违规-封禁
	VoteTypeDelete  = 4 // 违规-删除
	VoteTypeLegal   = 2 // 不违规
	VoteTypeGiveUp  = 3 // 放弃投票

	// JuryInvalid
	JuryBlocked = 1
	JuryExpire  = 2
	JuryAdmin   = 3
)

// Case is jury case info.
type Case struct {
	ID  int64 `json:"id"`
	Mid int64 `json:"mid"`

	Agree        int64 `json:"agree"`
	Against      int64 `json:"against"`
	AdminAgree   int64 `json:"ad_agree"`
	AdminAgainst int64 `json:"ad_against"`
	PunishStatus int8  `json:"pubish_status"`
	PunishResult int8  `json:"pubish_result"`
	BlockDays    int64 `json:"block_days"`
}

// Kpi is jury kpi info.
type Kpi struct {
	ID        int64     `json:"id"`
	Mid       int64     `json:"mid"`
	Point     int64     `json:"point"`
	Day       time.Time `json:"day"`
	Rate      int64     `json:"rate"`
	Rank      int64     `json:"rank"`
	RankPer   int64     `json:"rank_per"`
	RankTotal int64     `json:"rankTotal"`
	Expired   time.Time `json:"expired"`
	PreCount  int64     `json:"-"`
}

// KpiPoint is jury kpi point info.
type KpiPoint struct {
	ID           int64     `json:"id"`
	Mid          int64     `json:"mid"`
	Day          time.Time `json:"day"`
	Point        int64     `json:"point"`
	ActiveDays   int64     `json:"activeDays"`
	VoteTotal    int64     `json:"voteTotal"`
	VoteRadio    int64     `json:"voteRadio"`
	BlockedTotal int64     `json:"blockedTotal"`
	Expired      time.Time `json:"expired"`
	OpinionNum   int64     `json:"opinion_num"`
	OpinionLikes int64     `json:"opinion_likes"`
	OpinionHates int64     `json:"opinion_hates"`
}

// KpiData is jury kpi data info.
type KpiData struct {
	KpiPoint
	VoteRealTotal int64 `json:"vote_real_total"`
}
