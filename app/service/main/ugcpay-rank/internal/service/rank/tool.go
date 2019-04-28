package rank

import (
	"time"
)

// MonthVer return ver from time.Time
func MonthVer(t time.Time) int64 {
	return int64(t.Year()*100 + int(t.Month()))
}
