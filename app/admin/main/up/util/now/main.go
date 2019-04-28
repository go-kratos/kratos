// Package now is a time toolkit for golang.
//
// More details README here: https://github.com/jinzhu/now
//
//  import "github.com/jinzhu/now"
//
//  now.BeginningOfMinute() // 2013-11-18 17:51:00 Mon
//  now.BeginningOfDay()    // 2013-11-18 00:00:00 Mon
//  now.EndOfDay()          // 2013-11-18 23:59:59.999999999 Mon
package now

import "time"

// WeekStartDay set week start day, default is sunday
var WeekStartDay = time.Sunday

// TimeFormats default time formats will be parsed as
var TimeFormats = []string{"1/2/2006", "1/2/2006 15:4:5", "2006", "2006-1", "2006-1-2", "2006-1-2 15", "2006-1-2 15:4", "2006-1-2 15:4:5", "1-2", "15:4:5", "15:4", "15", "15:4:5 Jan 2, 2006 MST", "2006-01-02 15:04:05.999999999 -0700 MST"}

// Now now struct
type Now struct {
	time.Time
}

// New initialize Now with time
func New(t time.Time) *Now {
	return &Now{t}
}

// BeginningOfMinute beginning of minute
func BeginningOfMinute() time.Time {
	return New(time.Now()).BeginningOfMinute()
}

// BeginningOfHour beginning of hour
func BeginningOfHour() time.Time {
	return New(time.Now()).BeginningOfHour()
}

// BeginningOfDay beginning of day
func BeginningOfDay() time.Time {
	return New(time.Now()).BeginningOfDay()
}

// BeginningOfWeek beginning of week
func BeginningOfWeek() time.Time {
	return New(time.Now()).BeginningOfWeek()
}

// BeginningOfMonth beginning of month
func BeginningOfMonth() time.Time {
	return New(time.Now()).BeginningOfMonth()
}

// BeginningOfQuarter beginning of quarter
func BeginningOfQuarter() time.Time {
	return New(time.Now()).BeginningOfQuarter()
}

// BeginningOfYear beginning of year
func BeginningOfYear() time.Time {
	return New(time.Now()).BeginningOfYear()
}

// EndOfMinute end of minute
func EndOfMinute() time.Time {
	return New(time.Now()).EndOfMinute()
}

// EndOfHour end of hour
func EndOfHour() time.Time {
	return New(time.Now()).EndOfHour()
}

// EndOfDay end of day
func EndOfDay() time.Time {
	return New(time.Now()).EndOfDay()
}

// EndOfWeek end of week
func EndOfWeek() time.Time {
	return New(time.Now()).EndOfWeek()
}

// EndOfMonth end of month
func EndOfMonth() time.Time {
	return New(time.Now()).EndOfMonth()
}

// EndOfQuarter end of quarter
func EndOfQuarter() time.Time {
	return New(time.Now()).EndOfQuarter()
}

// EndOfYear end of year
func EndOfYear() time.Time {
	return New(time.Now()).EndOfYear()
}

// Monday monday
func Monday() time.Time {
	return New(time.Now()).Monday()
}

// Sunday sunday
func Sunday() time.Time {
	return New(time.Now()).Sunday()
}

// EndOfSunday end of sunday
func EndOfSunday() time.Time {
	return New(time.Now()).EndOfSunday()
}

// Parse parse string to time
func Parse(strs ...string) (time.Time, error) {
	return New(time.Now()).Parse(strs...)
}

// ParseInLocation parse string to time in location
func ParseInLocation(loc *time.Location, strs ...string) (time.Time, error) {
	return New(time.Now().In(loc)).Parse(strs...)
}

// MustParse must parse string to time or will panic
func MustParse(strs ...string) time.Time {
	return New(time.Now()).MustParse(strs...)
}

// MustParseInLocation must parse string to time in location or will panic
func MustParseInLocation(loc *time.Location, strs ...string) time.Time {
	return New(time.Now().In(loc)).MustParse(strs...)
}

// Between check now between the begin, end time or not
func Between(time1, time2 string) bool {
	return New(time.Now()).Between(time1, time2)
}
