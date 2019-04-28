package util

import (
	"fmt"
	"time"
)

const (
	// TimeFormat .
	TimeFormat = "2006-01-02 15:04:05"
)

// JSONTime .
type JSONTime time.Time

// MarshalJSON .
func (jt JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("%q", time.Time(jt).Format(TimeFormat))
	return []byte(stamp), nil
}

// Before .
func (jt JSONTime) Before(t time.Time) bool {
	return time.Time(jt).Before(t)
}
