package block

import (
	"time"
)

// DBUser .
type DBUser struct {
	ID     int64
	MID    int64
	Status BlockStatus
	CTime  time.Time
	MTime  time.Time
}

// DBUserDetail .
type DBUserDetail struct {
	ID         int64
	MID        int64
	BlockCount int64
	CTime      time.Time
	MTime      time.Time
}

// DBHistory .
type DBHistory struct {
	ID        int64
	MID       int64
	AdminID   int64
	AdminName string
	Source    BlockSource
	Area      BlockArea
	Reason    string
	Comment   string
	Action    BlockAction
	StartTime time.Time
	Duration  int64
	Notify    bool
	CTime     time.Time
	MTime     time.Time
}

// MCBlockInfo .
type MCBlockInfo struct {
	BlockStatus BlockStatus `json:"bs"`
	StartTime   int64       `json:"st"`
	EndTime     int64       `json:"et"`
}

// MCUserDetail .
type MCUserDetail struct {
	BlockCount int64 `json:"block_count"`
}
