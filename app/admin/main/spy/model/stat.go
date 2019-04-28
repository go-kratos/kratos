package model

import "time"

// type def.
const (
	AccountType = iota
	ArchiveType
	ActivityType
)

// Statistics def.
type Statistics struct {
	ID        int64     `json:"id"`
	TargetMid int64     `json:"target_mid"`
	TargetID  int64     `json:"target_id"`
	EventID   int64     `json:"event_id"`
	State     int8      `json:"state"`
	Type      int8      `json:"type"`
	Isdel     int8      `json:"is_del"`
	Quantity  int64     `json:"quantity"`
	EventName string    `json:"event_name"`
	Ctime     time.Time `json:"ctime"`
	Mtime     time.Time `json:"mtime"`
	CtimeUnix int64     `json:"ctimeunix"`
	MtimeUnix int64     `json:"mtimeunix"`
}

// StatPage def.
type StatPage struct {
	TotalCount int64         `json:"total_count"`
	Pn         int           `json:"pn"`
	Ps         int           `json:"ps"`
	Items      []*Statistics `json:"items"`
}
