package model

import (
	"encoding/json"
	"time"

	sml "go-common/app/service/main/relation/model"
)

// Message define binlog databus message.
type Message struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// Relation user_relation_fid_0~user_relation_fid_49,user_relation_mid_0~user_relation_mid_49
type Relation struct {
	Mid       int64  `json:"mid,omitempty"`
	Fid       int64  `json:"fid,omitempty"`
	Attribute uint32 `json:"attribute"`
	Status    int    `json:"status"`
	MTime     string `json:"mtime"`
	CTime     string `json:"ctime"`
}

// Stat user_relation_stat
type Stat struct {
	Mid       int64 `json:"mid,omitempty"`
	Following int64 `json:"following"`
	Whisper   int64 `json:"whisper"`
	Black     int64 `json:"black"`
	Follower  int64 `json:"follower"`
}

// LastChangeAt is.
func (r *Relation) LastChangeAt() (at time.Time, err error) {
	// FIXME(zhoujiahui): ctime and mtime should not be used here
	return time.ParseInLocation("2006-01-02 15:04:05", r.MTime, time.Local)
}

// Attr is.
func (r *Relation) Attr() uint32 {
	return sml.Attr(r.Attribute)
}

// IsRecent is.
func (r *Relation) IsRecent(at time.Time, trange time.Duration) bool {
	lastAt, err := r.LastChangeAt()
	if err != nil {
		return false
	}
	if lastAt.Sub(at) > trange {
		return true
	}
	return false
}
