package upcrmservice

import (
	"time"
)

//GetDateStamp get date from time stamp
func GetDateStamp(timeStamp time.Time) time.Time {
	var y, m, d = timeStamp.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, timeStamp.Location())
}
