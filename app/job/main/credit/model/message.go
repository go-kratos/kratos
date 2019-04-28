package model

import (
	xtime "go-common/library/time"
)

const (
	// RouteReplyReport report
	RouteReplyReport = "report_add"
)

// Reply param struct
type Reply struct {
	Action  string        `json:"action"`
	MID     int64         `json:"mid"`
	Subject *ReplySubject `json:"subject"`
	Reply   *ReplyMain    `json:"reply"`
	Report  *ReplyReport  `json:"report"`
}

// ReplySubject param struct
type ReplySubject struct {
	OID   int64      `json:"oid"`
	Type  int8       `json:"type"`
	MID   int64      `json:"mid"`
	State int8       `json:"state"`
	CTime xtime.Time `json:"ctime"`
}

// ReplyMain param struct
type ReplyMain struct {
	RPID    int64 `json:"rpid"`
	OID     int64 `json:"oid"`
	Type    int8  `json:"type"`
	MID     int64 `json:"mid"`
	Root    int64 `json:"root"`
	Parent  int64 `json:"parent"`
	Floor   int32 `json:"floor"`
	Count   int32 `json:"count"`
	Rcount  int32 `json:"rcount"`
	Like    int64 `json:"like"`
	Hate    int64 `json:"hate"`
	State   int8  `json:"state"`
	Content *struct {
		Message string `json:"message"`
	} `json:"content"`
	CTime xtime.Time `json:"ctime"`
}

// ReplyReport param struct
type ReplyReport struct {
	ID      int64      `json:"id"`
	OID     int64      `json:"oid"`
	Type    int8       `json:"type"`
	RPID    int64      `json:"rpid"`
	MID     int64      `json:"mid"`
	Reason  int8       `json:"reason"`
	Content string     `json:"content"`
	State   int8       `json:"state"`
	Score   int        `json:"score"`
	Count   int        `json:"count"`
	CTime   xtime.Time `json:"ctime"`
}

// LabourAnswer param struct
type LabourAnswer struct {
	MID   int64      `json:"mid"`
	MTime xtime.Time `json:"mtime"`
}
