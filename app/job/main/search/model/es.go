package model

import (
	"time"

	xtime "go-common/library/time"
)

// ES .
type ES struct {
	Addr string
}

// Offset .
type Offset struct {
	OffID      int64
	OffTime    xtime.Time
	ReviewID   int64
	ReviewTime int64
}

// OffsetID .
func (o *Offset) OffsetID() int64 {
	return o.OffID - o.ReviewID
}

// OffsetTime .
func (o *Offset) OffsetTime() string {
	return time.Unix(int64(o.OffTime), 0).Format("2006-01-02 15:04:05")
}

// var .
var (
	ExistsAction = map[string]bool{"repair": true, "all": true}
)
