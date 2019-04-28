package model

import "time"

// Beat ap_user_beats_info
type Beat struct {
	ID         int64     `json:"id"`
	UID        int64     `json:"uid"`
	Content    string    `json:"content"`
	Status     int       `json:"status"`
	Operator   string    `json:"operator"`
	UpdateTime time.Time `json:"update_time"`
	Ctime      time.Time `json:"ctime"`
	Mtime      time.Time `json:"mtime"`
}
