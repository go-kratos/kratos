package model

import (
	"time"
)

// dm transfer state
const (
	StatInit      = 0
	StatFinished  = 1
	StatFailed    = 2
	StatTransfing = 3
)

// Transfer dm transfer task
type Transfer struct {
	ID      int64
	FromCid int64
	ToCid   int64
	Mid     int64
	Offset  float64
	State   int8
	Dmid    int64
	Ctime   time.Time
	Mtime   time.Time
}
