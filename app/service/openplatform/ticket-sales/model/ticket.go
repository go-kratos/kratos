package model

import (
	xtime "go-common/library/time"
)

// Ticket ticket 表结构
type Ticket struct {
	ID              int64      `json:"id"`
	UID             int64      `json:"uid"`
	OID             int64      `json:"oid"`
	SID             int64      `json:"sid"`
	Price           int64      `json:"price"`
	Src             int16      `json:"src"`
	Type            int16      `json:"type"`
	Status          int16      `json:"status"`
	Qr              string     `json:"qr"`
	RefID           int64      `json:"ref_id"`
	SkuID           int64      `json:"sku_id"`
	SeatID          int64      `json:"seat_id"`
	Seat            string     `json:"seat"`
	RefundApplyTime xtime.Time `json:"refund_apply_time"`
	ETime           xtime.Time `json:"etime"`
	CTime           xtime.Time `json:"ctime"`
	MTime           xtime.Time `json:"mtime"`
}

// TicketSend ticket_send 表结构
type TicketSend struct {
	ID      int64      `json:"id"`
	SID     int64      `son:"sid"`
	SendTID int64      `json:"send_tid"`
	RecvTID int64      `json:"recv_tid"`
	SendUID int64      `json:"send_uid"`
	RecvUID int64      `json:"recv_uid"`
	RecvTel string     `json:"recv_tel"`
	Status  int16      `json:"status"`
	CTime   xtime.Time `json:"ctime"`
	MTime   xtime.Time `json:"mtime"`
	OID     int64      `json:"oid"`
}
