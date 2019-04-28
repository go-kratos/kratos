package model

import (
	"time"
)

// DBElecUPRank .
type DBElecUPRank struct {
	ID        int64
	UPMID     int64
	PayMID    int64
	PayAmount int64
	Ver       int64
	Hidden    bool
	CTime     time.Time
	MTime     time.Time
}

// DBElecAVRank .
type DBElecAVRank struct {
	ID        int64
	AVID      int64
	UPMID     int64
	PayMID    int64
	PayAmount int64
	Ver       int64
	Hidden    bool
	CTime     time.Time
	MTime     time.Time
}

// DBElecMessage .
type DBElecMessage struct {
	ID      int64
	Ver     int64
	AVID    int64
	UPMID   int64
	PayMID  int64
	Message string
	Replied bool
	Hidden  bool
	CTime   time.Time
	MTime   time.Time
}

// DBElecUserSetting .
type DBElecUserSetting struct {
	ID    int64
	MID   int64
	Value int
	CTime time.Time
	MTime time.Time
}
