package like

import (
	xtime "go-common/library/time"
)

// action type
const (
	LIKESCORE = 1
)

// Action def.
type Action struct {
	ID     int64      `json:"id"`
	Lid    int64      `json:"lid"`
	Mid    int64      `json:"mid"`
	Action int64      `json:"action"`
	Ctime  xtime.Time `json:"ctime"`
	Mtime  xtime.Time `json:"mtime"`
	Sid    int64      `json:"sid"`
	IP     int64      `json:"ip"`
	IPv6   []byte     `json:"ipv6"`
}

// LidLikeSum def .
type LidLikeSum struct {
	Likes int64
	Lid   int64
}
