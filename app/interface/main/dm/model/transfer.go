package model

import (
	"time"

	xtime "go-common/library/time"
)

// all variable used in dm transfer
const (
	TransferJobStatInit     = int8(0)
	TransferJobStatFinished = int8(1)
	TransferJobStatFailed   = int8(2)
)

// TransferJob dm transfer
type TransferJob struct {
	ID      int64
	FromCID int64
	ToCID   int64
	MID     int64
	Offset  float64
	State   int8
	Ctime   time.Time
	Mtime   time.Time
}

// TransferHistory transfer list item
type TransferHistory struct {
	ID     int64      `json:"id"`
	PartID int32      `json:"part_id"`
	CID    int64      `json:"cid"`
	Title  string     `json:"title"`
	CTime  xtime.Time `json:"ctime"`
	State  int8       `json:"state"`
}

// CidInfo is archive_video model.
type CidInfo struct {
	Aid        int64      `json:"aid"`
	Title      string     `json:"title"`
	Desc       string     `json:"desc"`
	Filename   string     `json:"filename"`
	Index      int        `json:"index"`
	Status     int16      `json:"status"`
	StatusDesc string     `json:"status_desc"`
	FailCode   int8       `json:"fail_code"`
	FailDesc   string     `json:"fail_desc"`
	CTime      xtime.Time `json:"ctime"`
}
