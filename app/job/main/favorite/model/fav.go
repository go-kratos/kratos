package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

var (
	// ErrFavResourceExist error this has been favoured.
	ErrFavResourceExist = errors.New("error this has been favoured")
	// ErrFavResourceAlreadyDel error this has been unfavoured.
	ErrFavResourceAlreadyDel = errors.New("error this has been unfavoured")
)

const (
	// CacheNotFound .
	CacheNotFound = -1
	// SyncInsert binlog action.
	SyncInsert = "insert"
	// SyncUpdate binlog action.
	SyncUpdate = "update"
	// SyncDelete binlog action.
	SyncDelete = "delete"
)

// CanelMessage binlog message.
type CanelMessage struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// OldCount .
type OldCount struct {
	ID    int64 `json:"id"`
	Aid   int64 `json:"aid"`
	Count int64 `json:"count"`
	CTime Stime `json:"ctime"`
	MTime Stime `json:"mtime"`
}

// NewCount .
type NewCount struct {
	ID    int64 `json:"id"`
	Type  int8  `json:"type"`
	Oid   int64 `json:"oid"`
	Count int64 `json:"count"`
	CTime Stime `json:"ctime"`
	MTime Stime `json:"mtime"`
}

// OldFolder .
type OldFolder struct {
	ID       int64  `json:"id"`
	Mid      int64  `json:"mid"`
	Name     string `json:"name"`
	CurCount int    `json:"cur_count"`
	State    int8   `json:"state"`
	CTime    Stime  `json:"ctime"`
	MTime    Stime  `json:"mtime"`
}

// NewFolder .
type NewFolder struct {
	ID    int64  `json:"id"`
	Type  int8   `json:"type"`
	Mid   int64  `json:"mid"`
	Name  string `json:"name"`
	Count int    `json:"count"`
	Attr  int8   `json:"attr"`
	State int8   `json:"state"`
	CTime Stime  `json:"ctime"`
	MTime Stime  `json:"mtime"`
}

// VideoFolder .
type VideoFolder struct {
	ID       int64 `json:"id"`
	Mid      int64 `json:"mid"`
	Fid      int64 `json:"fid"`
	VideoFid int64 `json:"video_fid"`
	CTime    Stime `json:"ctime"`
	MTime    Stime `json:"mtime"`
}

// OldVideo .
type OldVideo struct {
	ID    int64 `json:"id"`
	Mid   int64 `json:"mid"`
	Fid   int64 `json:"fid"`
	Aid   int64 `json:"aid"`
	CTime Stime `json:"ctime"`
	MTime Stime `json:"mtime"`
}

// NewRelation .
type NewRelation struct {
	ID    int64 `json:"id"`
	Type  int8  `json:"type"`
	Mid   int64 `json:"mid"`
	Fid   int64 `json:"fid"`
	Oid   int64 `json:"oid"`
	State int8  `json:"state"`
	CTime Stime `json:"ctime"`
	MTime Stime `json:"mtime"`
}

// OldFolderSort .
type OldFolderSort struct {
	ID    int64  `json:"id"`
	Mid   int64  `json:"mid"`
	Sort  string `json:"sort"`
	CTime Stime  `json:"ctime"`
	MTime Stime  `json:"mtime"`
}

// NewFolderSort .
type NewFolderSort struct {
	ID    int64  `json:"id"`
	Type  int8   `json:"type"`
	Mid   int64  `json:"mid"`
	Sort  []byte `json:"sort"`
	CTime Stime  `json:"ctime"`
	MTime Stime  `json:"mtime"`
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

// StatMsg .
type StatMsg struct {
	Play  *int64 `json:"play"`
	Fav   *int64 `json:"fav"`
	Share *int64 `json:"share"`
	Oid   int64  `json:"oid"`
}

// StatCount .
type StatCount struct {
	Type      string `json:"type"`
	ID        int64  `json:"id"`
	Count     int64  `json:"count"`
	DisLike   int64  `json:"dislike_count"`
	TimeStamp int64  `json:"timestamp"`
}

// PlayReport .
type PlayReport struct {
	ID       int64  `json:"id"`
	Mid      int64  `json:"mid"`
	LV       string `json:"lv"`
	IP       string `json:"ip"`
	Buvid    string `json:"buvid"`
	DeviceID string `json:"device_id"`
	UA       string `json:"ua"`
	Refer    string `json:"refer"`
	TS       int64  `json:"ts"`
}
