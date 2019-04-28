package model

import (
	xtime "go-common/library/time"
)

// Case is jury case info.
type Case struct {
	ID     int64 `json:"id"`
	Mid    int64 `json:"mid"`
	Status int64 `json:"status"`
	Origin
	JudgeType    int64  `json:"judge_type"`
	PunishResult int64  `json:"punish_result"`
	Agree        int64  `json:"vote_rule"`
	Against      int64  `json:"vote_break"`
	VoteDelete   int64  `json:"vote_delete"`
	PunishStatus int64  `json:"pubish_status"`
	BlockedDay   int64  `json:"blocked_days"`
	RelationID   string `json:"relation_id"`
	Operator     string `json:"operator"`
	PutTotal     int64  `json:"put_total"`
	Stime        string `json:"start_time"`
	Etime        string `json:"end_time"`
	Ctime        string `json:"ctime"`
	CaseType     int8   `json:"case_type"`
	OPID         int64  `json:"oper_id"`
	BusinessTime string `json:"business_time"`
	BCtime       xtime.Time
}

// Publish is blocked_publish info.
type Publish struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Subtitle string `json:"sub_title"`
	PStatus  int8   `json:"publish_status"`
	SStatus  int8   `json:"stick_status"`
	Content  string `json:"content"`
	URL      string `json:"url"`
	PType    int8   `json:"ptype"`
	STime    string `json:"show_time"`
}

// SimCase is simple case struct info.
type SimCase struct {
	ID         int64      `json:"id"`
	Mid        int64      `json:"mid"`
	VoteRule   int64      `json:"vote_rule"`
	VoteBreak  int64      `json:"vote_break"`
	VoteDelete int64      `json:"vote_delete"`
	CaseType   int8       `json:"case_type"`
	Stime      xtime.Time `json:"start_time"`
	Etime      xtime.Time `json:"end_time"`
}

// Jury is simple jury struct info.
type Jury struct {
	ID     int64 `json:"id"`
	Mid    int64 `json:"mid"`
	Status int8  `json:"status"`
}

// CaseVote is simple caseVote struct info.
type CaseVote struct {
	ID      int64      `json:"id"`
	CID     int64      `json:"cid"`
	MID     int64      `json:"mid"`
	Vote    int8       `json:"vote"`
	Expired xtime.Time `json:"expired"`
}

// BLogCaseVote is simple blogCaseVote struct info.
type BLogCaseVote struct {
	ID      int64  `json:"-"`
	CID     int64  `json:"cid"`
	MID     int64  `json:"mid"`
	Vote    int8   `json:"-"`
	Expired string `json:"-"`
	Ctime   string `json:"-"`
	Mtime   string `json:"-"`
}

// Opinion is simple opinion struct info.
type Opinion struct {
	Cid     int64  `json:"cid"`
	Vid     int64  `json:"vid"`
	Content string `json:"content"`
	State   int8   `json:"state"`
}

// BlockDays get user blocked days.
func (c *Case) BlockDays() (forever int8, days int64) {
	switch c.PunishResult {
	case Punish3Days:
		days = 3
	case Punish7Days:
		days = 7
	case Punish15Days:
		days = 15
	case PunishForever:
		forever = int8(1)
	case PunishCustom:
		days = c.BlockedDay
	}
	return
}

// Origin is origin info of blocked.
type Origin struct {
	OriginTitle         string `json:"origin_title"`
	OriginURL           string `json:"origin_url"`
	OriginContent       string `json:"origin_content"`
	OriginContentModify string `json:"origin_content_modify"`
	OriginType          int64  `json:"origin_type"`
	ReasonType          int64  `json:"reason_type"`
}

// BlockedInfo user block info.
type BlockedInfo struct {
	ID  int64 `json:"id"`
	UID int64 `json:"uid"`
	Origin
	BlockedRemark  string `json:"blocked_remark"`
	PunishTime     string `json:"punish_time"`
	PunishType     int64  `json:"punish_type"`
	MoralNum       int64  `json:"moral_num"`
	BlockedDays    int64  `json:"blocked_days"`
	PublishStatus  int64  `json:"publish_status"`
	BlockedType    int64  `json:"blocked_type"`
	BlockedForever int8   `json:"blocked_forever"`
	OperatorName   string `json:"operator_name"`
	CaseID         int64  `json:"case_id"`
	OPID           int64  `json:"oper_id"`
	Status         int64  `json:"status"`
	MTime          string `json:"mtime"`
}

// Kpi is jury kpi info.
type Kpi struct {
	ID            int64 `json:"id"`
	Mid           int64 `json:"mid"`
	Rate          int8  `json:"rate"`
	Rank          int64 `json:"rank"`
	RankPer       int64 `json:"rank_per"`
	RankTotal     int64 `json:"rankTotal"`
	HandlerStatus int64 `json:"handlerStatus"`
}

// PunishResultDays punish res days.
func PunishResultDays(blockedTimes int64) (punishResult, blockedDay int64) {
	switch {
	case blockedTimes == 0:
		punishResult = int64(Block7Days)
		blockedDay = BlockTimeSeven
	case blockedTimes == 1:
		punishResult = int64(Block15Days)
		blockedDay = BlockTimeFifteen
	case blockedTimes > 1:
		punishResult = int64(BlockForever)
		blockedDay = BlockTimeForever
	}
	return
}
