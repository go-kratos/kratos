package model

import "time"

// ShieldKeyWork  ShieldKeyWorks
type ShieldKeyWork struct {
	ID              int64     `json:"id"`
	UID             int64     `json:"uid"`
	OriginalKeyword string    `json:"origin_keyword"`
	KeyWord         string    `json:"keyword"`
	Ctime           time.Time `json:"ctime"`
}
