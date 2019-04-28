package model

import (
	xtime "go-common/library/time"
	"time"
)

// Pagination is
type Pagination struct {
	Sort      string `form:"sort"`
	Order     string `form:"order"`
	PS        int    `form:"ps" validate:"min=0,max=50"`
	PN        int    `form:"pn" validate:"min=0"`
	MTimeFrom string `form:"mtime_from"`
	MTimeTo   string `form:"mtime_to"`
}

// FollowersParam is
type FollowersParam struct {
	Pagination
	Fid int64 `form:"fid" validate:"min=1,required"`
	Mid int64 `form:"mid" validate:"min=0"`
}

// FollowingsParam is
type FollowingsParam struct {
	Pagination
	Mid int64 `form:"mid" validate:"min=1,required"`
	Fid int64 `form:"fid" validate:"min=0"`
}

// LogsParam is
type LogsParam struct {
	Mid int64 `form:"mid" validate:"min=1,required"`
	Fid int64 `form:"fid" validate:"min=1,required"`
}

// ParseTime is
func ParseTime(ts string) (xt xtime.Time, err error) {
	var (
		t time.Time
	)
	if t, err = time.Parse("2006-01-02 15:04:05", ts); err != nil {
		return
	}
	xt.Scan(t)
	return
}

// Desc is
func (p Pagination) Desc() bool {
	return p.Sort == "desc"
}

// ArgMid is
type ArgMid struct {
	Mid int64 `form:"mid" validate:"min=1,required"`
}

// ArgMids is
type ArgMids struct {
	Mids []int64 `form:"mids,split" validate:"dive,gt=0"`
}
