package operator

import (
	"time"

	"go-common/library/log"
	xtime "go-common/library/time"
)

type Reddot struct {
	StartTime xtime.Time `json:"start_time,omitempty"`
	EndTime   xtime.Time `json:"end_time,omitempty"`
}

// ReddotChange  Reddot change
func (r *Reddot) ReddotChange(startStr, endStr string) {
	if startStr != "" && endStr != "" {
		r.StartTime = timeStrToInt(startStr)
		r.EndTime = timeStrToInt(endStr)
	}
}

// timeStrToInt string to int
func timeStrToInt(timeStr string) (timeInt xtime.Time) {
	var err error
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, timeStr, loc)
	if err = timeInt.Scan(theTime); err != nil {
		log.Error("timeInt.Scan error(%v)", err)
	}
	return
}
