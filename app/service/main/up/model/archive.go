package model

import (
	"go-common/library/time"
)

//Archive for db.
type Archive struct {
	ID    int64 `json:"id"`
	MID   int64 `json:"mid"`
	State int   `json:"state"`
}

// AidPubTime aid's pubdate and copyright
type AidPubTime struct {
	Aid       int64     `json:"aid"`
	PubDate   time.Time `json:"pubdate"`
	Copyright int8      `json:"copyright"`
}
