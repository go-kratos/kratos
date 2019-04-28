package model

import xtime "go-common/library/time"

// ArgPage define page param.
type ArgPage struct {
	PN int64 `form:"pn" default:"1"`
	PS int64 `form:"ps" validate:"min=0,max=30" default:"30"`
}

// ArgBlockedNumUser user blocked number param.
type ArgBlockedNumUser struct {
	MID int64 `form:"mid" validate:"required"`
}

// ResBlockedNumUser  user blocked number result.
type ResBlockedNumUser struct {
	BlockedSum int `json:"blockedSum"`
}

// ArgIDs ids param.
type ArgIDs struct {
	IDs []int64 `form:"ids,split" validate:"min=0,max=100"`
}

// ArgMIDs mids param.
type ArgMIDs struct {
	MIDs []int64 `form:"mids,split" validate:"min=0,max=100"`
}

// ResJuryerStatus blocked juryer status result.
type ResJuryerStatus struct {
	Expired xtime.Time `json:"expired"`
	Mid     int64      `json:"mid"`
	Status  int8       `json:"status"`
}

// ArgJudgeBlocked judge blocked param.
type ArgJudgeBlocked struct {
	MID      int64  `form:"mid" validate:"required"`
	OID      int64  `form:"oper_id" default:"0"`
	BDays    int    `form:"blocked_days" default:"0"`
	BForever int8   `form:"blocked_forever" default:"0"`
	BRemark  string `form:"blocked_remark" default:""`
	MoralNum int    `form:"moral_num" default:"0"`
	OContent string `form:"origin_content" default:""`
	OTitle   string `form:"origin_title" default:""`
	OType    int8   `form:"origin_type"  validate:"min=1,max=20"`
	OURL     string `form:"origin_url" default:""`
	PTime    int64  `form:"punish_time" validate:"required"`
	PType    int8   `form:"punish_type"  validate:"min=1,max=10"`
	RType    int8   `form:"reason_type" validate:"min=1,max=40"`
	OPName   string `form:"operator_name" default:""`
}

// ArgJudgeBatchBlocked judge batch blocked param.
type ArgJudgeBatchBlocked struct {
	MID      []int64 `form:"mids,split"   validate:"min=1,max=200"`
	OID      int64   `form:"oper_id" default:"0"`
	BDays    int     `form:"blocked_days" default:"0"`
	BForever int8    `form:"blocked_forever" default:"0"`
	BRemark  string  `form:"blocked_remark" default:""`
	MoralNum int     `form:"moral_num" default:"0"`
	OContent string  `form:"origin_content" default:""`
	OTitle   string  `form:"origin_title" default:""`
	OType    int8    `form:"origin_type" validate:"min=1,max=20"`
	OURL     string  `form:"origin_url" default:""`
	PTime    int64   `form:"punish_time"  validate:"required"`
	PType    int8    `form:"punish_type" validate:"min=1,max=10"`
	RType    int8    `form:"reason_type" validate:"min=1,max=40"`
	OPName   string  `form:"operator_name" default:""`
}

// ArgHistory blocked historys param.
type ArgHistory struct {
	MID   int64 `form:"mid"  validate:"required"`
	STime int64 `form:"start"  validate:"required"`
	PN    int   `form:"pn" default:"1"`
	PS    int   `form:"ps" validate:"min=0,max=100" default:"100"`
}

// ResBLKHistorys blocked historys result.
type ResBLKHistorys struct {
	TotalCount int64          `json:"total_count"`
	PN         int            `json:"pn"`
	PS         int            `json:"ps"`
	Items      []*BlockedInfo `json:"items"`
}

// ArgJudgeCase judge case param.
type ArgJudgeCase struct {
	AID          int64      `json:"aid"`
	MID          int64      `json:"mid"`
	Operator     string     `json:"operator"`
	OperID       int64      `json:"oper_id"`
	OContent     string     `json:"origin_content"`
	OTitle       string     `json:"origin_title"`
	OType        int64      `json:"origin_type"`
	OURL         string     `json:"origin_url"`
	ReasonType   int64      `json:"reason_type"`
	OID          int64      `json:"oid"`
	RPID         int64      `json:"rp_id"`
	TagID        int64      `json:"tag_id"`
	Type         int64      `json:"type"`
	Page         int64      `json:"page"`
	BCTime       xtime.Time `json:"business_time"`
	RelationID   string     `json:"-"`
	PunishResult int8       `json:"-"`
	BlockedDays  int32      `json:"-"`
}

// ArgDElQS labour question del param.
type ArgDElQS struct {
	ID    int64 `form:"id"  validate:"required"`
	IsDel int64 `form:"is_del"  validate:"min=1,max=3"`
}

// ArgBlockedList  blocked list param.
type ArgBlockedList struct {
	OType int8 `form:"otype" default:"0"`
	BType int8 `form:"btype" default:"-1"`
	PN    int  `form:"pn" validate:"min=1" default:"1"`
	PS    int  `form:"ps" default:"20"`
}
