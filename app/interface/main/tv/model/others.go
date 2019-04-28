package model

import "go-common/library/time"

// Channel defines the structure of the channel & splash data
type Channel struct {
	ID      int
	Title   string
	Desc    string
	Splash  string
	Deleted int
	Ctime   time.Time
	Mtime   time.Time
}

//ReqTransode is the request structure for the transcode api
type ReqTransode struct {
	ContType string `form:"cont_type" validate:"required"` // content type: pgc/ugc
	CID      int64  `form:"cid" validate:"required"`
	Action   int64  `form:"action" validate:"min=0,max=2"` // 1 = finished, others = failed
}

// Hotword item def.
type Hotword struct {
	Keyword string `json:"keyword"`
	Status  string `json:"status"`
}

// ReqApply is request for apply time storage
type ReqApply struct {
	CID       int64 `form:"cid" validate:"required"`
	ApplyTime int64 `form:"apply_time" validate:"required"`
}
