package model

// DmHistoryParams .
type DmHistoryParams struct {
	Bsp       *BasicSearchParams
	Oid       int64   `form:"oid" params:"oid" default:"-1"`
	States    []int64 `form:"states,split" params:"states"`
	CtimeFrom string  `form:"ctime_from" params:"ctime_from"`
	CtimeTo   string  `form:"ctime_to" params:"ctime_to"`
}
