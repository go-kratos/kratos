package model

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"
	"time"

	xtime "go-common/library/time"
)

const (
	// BusinessArchive .
	BusinessArchive = "archive"
	// RankOrderByDesc .
	RankOrderByDesc = "desc"
	// RankOrderByAsc .
	RankOrderByAsc = "asc"
	// SyncInsert .
	SyncInsert = "insert"
	// SyncUpdate .
	SyncUpdate = "update"
	// SyncDelete .
	SyncDelete = "delete"
	// TimeFormat .
	TimeFormat = "2006-01-02 15:04:05"
	// FlagExist .
	FlagExist = true
)

// ArchiveMeta .
type ArchiveMeta struct {
	ID      int64 `json:"id"`
	Aid     int64 `json:"aid"`
	Typeid  int64 `json:"typeid"`
	Pubtime Stime `json:"pubtime"`
	*ArchiveType
	*ArchiveStat
	*ArchiveTv
}

// ArchiveType .
type ArchiveType struct {
	ID  int64 `json:"id"`
	Pid int64 `json:"pid"`
}

// ArchiveStat .
type ArchiveStat struct {
	ID    int64 `json:"id"`
	Aid   int64 `json:"aid"`
	Click int64 `json:"click"`
}

// ArchiveTv .
type ArchiveTv struct {
	ID      int64 `json:"id"`
	Aid     int64 `json:"aid"`
	Result  int8  `json:"result"`
	Deleted int8  `json:"deleted"`
	Valid   int8  `json:"valid"`
}

// StatViewMsg .
type StatViewMsg struct {
	Type      string `json:"type"`
	ID        int64  `json:"id"`
	Count     int    `json:"count"`
	Timestamp int64  `json:"timestamp"`
}

// CanalMsg .
type CanalMsg struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// SetPubtime .
func (a *ArchiveMeta) SetPubtime() xtime.Time {
	return xtime.Time(a.Pubtime)
}

// SetPid .
func (a *ArchiveType) SetPid() int16 {
	return int16(a.Pid)
}

// SetClick .
func (a *ArchiveStat) SetClick() int {
	return int(a.Click)
}

// DoReq .
type DoReq struct {
	Business  string `form:"business" validate:"required"`
	Action    string `form:"action" validate:"required"`
	MinID     int64  `form:"minid"`
	MaxID     int64  `form:"maxid"`
	BeginTime string `form:"begintime"`
	EndTime   string `form:"endtime"`
}

// MgetReq .
type MgetReq struct {
	Business string  `form:"business" validate:"required"`
	Oids     []int64 `form:"oids,split" validate:"required"`
}

// MgetResp resp of mget
type MgetResp struct {
	List map[int64]*Field `json:"list"`
}

// SortReq .
type SortReq struct {
	Business string            `form:"business" validate:"required"`
	Field    string            `form:"field" validate:"required"`
	Order    string            `form:"order" validate:"required"`
	Filters  map[string]string `form:"filters" validate:"required"`
	Oids     []int64           `form:"oids,split" validate:"required"`
	Pn       int               `form:"pn"`
	Ps       int               `form:"ps"`
}

// SortResp .
type SortResp struct {
	Result []int64 `json:"result"`
	Page   *Page   `json:"page"`
}

// GroupReq .
type GroupReq struct {
	Business string  `form:"business" validate:"required"`
	Field    string  `form:"field" validate:"required"`
	Oids     []int64 `form:"Oids,split" validate:"required"`
}

// GroupResp .
type GroupResp struct {
	List []*Group `json:"list"`
}

// Group .
type Group struct {
	Key   string `json:"key"`
	Count int    `json:"count"`
}

// Page Pager
type Page struct {
	Pn    int `json:"pn"`
	Ps    int `json:"ps"`
	Total int `json:"total"`
}

// Stime .
type Stime int64

// Scan scan time.
func (st *Stime) Scan(src interface{}) (err error) {
	switch sc := src.(type) {
	case time.Time:
		*st = Stime(sc.Unix())
	case string:
		var i int64
		i, err = strconv.ParseInt(sc, 10, 64)
		*st = Stime(i)
	}
	return
}

// Value get time value.
func (st Stime) Value() (driver.Value, error) {
	return time.Unix(int64(st), 0), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (st *Stime) UnmarshalJSON(data []byte) error {
	timestamp, err := strconv.ParseInt(string(data), 10, 64)
	if err == nil {
		*st = Stime(timestamp)
		return nil
	}
	t, err := time.ParseInLocation(`"2006-01-02 15:04:05"`, string(data), time.Local)
	if err == nil {
		*st = Stime(t.Unix())
	}
	return nil
}
