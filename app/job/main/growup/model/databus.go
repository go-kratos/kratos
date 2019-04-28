package model

import "encoding/json"

// ArchiveMsg archive-T databus msg.
type ArchiveMsg struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// ArchiveSub archive
type ArchiveSub struct {
	ID        int64  `json:"id"`
	MID       int64  `json:"mid"`
	Copyright int8   `json:"copyright"`
	State     int    `json:"state"`
	CTime     string `json:"ctime"`
	MTime     string `json:"mtime"`
	PTime     string `json:"ptime"`
}

// BgmSub bgm sub
type BgmSub struct {
	MID   int64 `json:"mid"`
	State int   `json:"state"`
}
