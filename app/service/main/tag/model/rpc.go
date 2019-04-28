package model

import (
	"time"
)

// ArgID .
type ArgID struct {
	ID     int64
	Mid    int64
	RealIP string
}

// ArgName .
type ArgName struct {
	Name   string
	Mid    int64
	RealIP string
}

// ArgIDs .
type ArgIDs struct {
	IDs    []int64
	Mid    int64
	RealIP string
}

// ArgNames .
type ArgNames struct {
	Names  []string
	Mid    int64
	RealIP string
}

// ArgAid .
type ArgAid struct {
	Aid    int64
	Mid    int64
	RealIP string
}

// ArgSub .
type ArgSub struct {
	Mid    int64
	Type   int
	Pn     int
	Ps     int
	Order  int
	RealIP string
}

// ArgCustomChannel .
type ArgCustomChannel struct {
	Mid    int64
	Type   int
	Order  int
	RealIP string
}

// ArgAddSub .
type ArgAddSub struct {
	Mid    int64
	Tids   []int64
	RealIP string
}

// ArgCancelSub .
type ArgCancelSub struct {
	Mid    int64
	Tid    int64
	RealIP string
}

// ArgCustomSub .
type ArgCustomSub struct {
	Mid    int64
	Type   int
	Tids   []int64
	RealIP string
}

// ArgResTags .
type ArgResTags struct {
	Oid    int64
	Type   int32
	Mid    int64
	RealIP string
}

// ArgResTag arg res tag.
type ArgResTag struct {
	Oid  int64
	Type int32
}

// ArgMutiResTag arg muti res tags.
type ArgMutiResTag struct {
	Oids []int64
	Type int32
}

// ArgTagLog .
type ArgTagLog struct {
	Oid    int64
	Type   int32
	LogID  int64
	RealIP string
}

// ArgResAction .
type ArgResAction struct {
	Mid    int64
	Tid    int64
	Oid    int64
	Type   int32
	RealIP string
}

// ArgHide .
type ArgHide struct {
	Tid    int64
	State  int32
	RealIP string
}

// ArgResActions .
type ArgResActions struct {
	Mid    int64
	Tids   []int64
	Oid    int64
	Type   int32
	RealIP string
}

// ArgResTagLog .
type ArgResTagLog struct {
	Oid    int64
	Type   int32
	Mid    int64
	Pn     int
	Ps     int
	RealIP string
}

// ResSub .
type ResSub struct {
	Tags  []*Tag `json:"tags"`
	Total int    `json:"total"`
}

// ResSubSort .
type ResSubSort struct {
	Sort  []*Tag `json:"sort"`
	Tags  []*Tag `json:"tags"`
	Total int    `json:"total"`
}

// ArgUserDel .
type ArgUserDel struct {
	Oid    int64
	Tid    int64
	Type   int8
	Mid    int64
	Role   int8
	RealIP string
}

// ArgUserAdd .
type ArgUserAdd struct {
	Oid    int64
	Mid    int64
	Type   int8
	Name   string
	Role   int8
	RealIP string
}

// ArgBind .
type ArgBind struct {
	Resource *Resource
	Action   int32
	Now      time.Time
	RealIP   string
}

// ArgCheckName .
type ArgCheckName struct {
	Name   string
	Type   int8
	Now    time.Time
	RealIP string
}

// ArgCreate .
type ArgCreate struct {
	Tag    *Tag
	Tags   []*Tag
	RealIP string
}

// ArgUPBind .
type ArgUPBind struct {
	Oid    int64
	Mid    int64
	Tids   []int64
	Type   int32
	RealIP string
}

// ArgUserBind .
type ArgUserBind struct {
	Oid    int64
	Mid    int64
	Tid    int64
	Type   int32
	Role   int32
	Action int32
	RealIP string
}

// ArgReportAction .
type ArgReportAction struct {
	Oid     int64
	Mid     int64
	LogID   int64
	Type    int32
	PartID  int32
	Reason  int32
	Score   int32
	Content string
	RealIP  string
}

// ArgRes .
type ArgRes struct {
	Tid    int64
	Type   int32
	Limit  int64
	RealIP string
}

// ArgRankingRegion .
type ArgRankingRegion struct {
	Rid    int64
	RealIP string
}

// ResBangumi .
type ResBangumi struct {
	Sids    []int64
	Bangumi []*RankingBangumi
}

// ArgHots .
type ArgHots struct {
	Rid    int64
	Type   int64
	RealIP string
}

// ArgDefaultBind .
type ArgDefaultBind struct {
	Oid    int64
	Mid    int64
	Type   int32
	Tids   []int64
	RealIP string
}

// ArgChannelCategory arg channel category.
type ArgChannelCategory struct {
	LastID int64
	Size   int32
	State  int32
}

// ArgChannels arg channels.
type ArgChannels struct {
	LastID int64
	Size   int32
}

// ArgChannelRule arg channelRules.
type ArgChannelRule struct {
	LastID int64
	Size   int32
	State  int32
}
