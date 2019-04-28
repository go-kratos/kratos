package model

import (
	"time"
)

const (
	// ProtectApplyLimit protect apply limit
	ProtectApplyLimit = 20
)

// Pager comment
type Pager struct {
	Total      int `json:"total"`
	Current    int `json:"current"`
	Size       int `json:"size"`
	TotalCount int `json:"total_count"`
}

// Pa 保护弹幕
type Pa struct {
	ID       int64
	CID      int64
	UID      int64
	ApplyUID int64
	AID      int64
	Playtime float32
	DMID     int64
	Msg      string
	Status   int
	Ctime    time.Time
	Mtime    time.Time
}

// Apply apply protect dm
type Apply struct {
	ID       int64   `json:"id"`
	AID      int64   `json:"aid"`
	CID      int64   `json:"cid"`
	Title    string  `json:"title"`
	ApplyUID int64   `json:"-"`
	Pic      string  `json:"pic"`
	Uname    string  `json:"uname"`
	Msg      string  `json:"msg"`
	Playtime float32 `json:"playtime"`
	Ctime    string  `json:"ctime"`
}

// ApplySortPlaytime what
type ApplySortPlaytime []*Apply

func (c ApplySortPlaytime) Len() int {
	return len(c)
}

func (c ApplySortPlaytime) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c ApplySortPlaytime) Less(i, j int) bool {
	if c[i].CID == c[j].CID {
		return c[i].Playtime < c[j].Playtime
	}
	return c[i].CID > c[j].CID
}

// ApplySortID what
type ApplySortID []*Apply

// Len get len
func (c ApplySortID) Len() int {
	return len(c)
}

// Swap change dm
func (c ApplySortID) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// Less count
func (c ApplySortID) Less(i, j int) bool {
	return c[i].ID > c[j].ID
}

// ApplyListResult get
type ApplyListResult struct {
	Pager *Pager
	List  []*Apply
}

// Video video info
type Video struct {
	Aid   int64  `json:"aid"`
	Title string `json:"title"`
}

// ApplyUserStat user stat
type ApplyUserStat struct {
	Aid    int64
	UID    int64
	Status int
	Ctime  time.Time
}

// ApplyUserNotify user notify
type ApplyUserNotify struct {
	Title     string
	Aid       int64
	Protect   int
	Unprotect int
}
