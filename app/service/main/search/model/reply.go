package model

import "fmt"

// ReplyRecordParams search params.
type ReplyRecordParams struct {
	Bsp       *BasicSearchParams
	Mid       int64   `form:"mid" params:"mid"`
	Types     []int64 `form:"types,split" params:"types"`
	States    []int64 `form:"states,split" params:"states"`
	CTimeFrom string  `form:"ctime_from" params:"ctime_from"`
	CTimeTo   string  `form:"ctime_to" params:"ctime_to"`
}

// ReplyRecordUpdateParams search params.
type ReplyRecordUpdateParams struct {
	ID    int64 `json:"id"`
	OID   int64 `json:"oid"`
	MID   int64 `json:"mid"`
	State int   `json:"state"`
}

// IndexName .
func (m *ReplyRecordUpdateParams) IndexName() string {
	return fmt.Sprintf("replyrecord_%d", m.MID%100)
}

// IndexType .
func (m *ReplyRecordUpdateParams) IndexType() string {
	return "base"
}

// IndexID .
func (m *ReplyRecordUpdateParams) IndexID() string {
	return fmt.Sprintf("%d_%d", m.ID, m.OID)
}
