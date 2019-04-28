package util

import (
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/up/util/timerqueue"
	"go-common/library/log"
	"go-common/library/net/http/blademaster"
	xtime "go-common/library/time"
)

var (
	//GlobalTimer timer queue
	GlobalTimer = timerqueue.New()
)

//ParseCommonTime parse to library common time
func ParseCommonTime(layout string, value string) (t xtime.Time, err error) {
	date, e := time.ParseInLocation(layout, value, time.Local)
	err = e
	if err == nil {
		t = xtime.Time(date.Unix())
	}
	return
}

//Unique unique the slice
func Unique(intSlice []int64) []int64 {
	keys := make(map[int64]bool)
	var list []int64
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

//GetContextValueInt64 get context int64
func GetContextValueInt64(c *blademaster.Context, key string) (v int64, ok bool) {
	var vtemp, o = c.Get(key)
	ok = o
	if ok {
		v, _ = vtemp.(int64)
	}
	return
}

//GetContextValueString get context string
func GetContextValueString(c *blademaster.Context, key string) (v string, ok bool) {
	var vtemp, o = c.Get(key)
	ok = o
	if ok {
		v, _ = vtemp.(string)
	}
	return
}

const (
	trimSet = "\r\n "
)

//ExplodeInt64 explode string to slice
func ExplodeInt64(str string, seperator string) (result []int64) {
	var strMids = strings.Split(str, seperator)
	for _, v := range strMids {
		mid, e := strconv.ParseInt(strings.Trim(v, trimSet), 10, 64)
		if e != nil {
			continue
		}
		result = append(result, mid)
	}
	return
}

//ExplodeUint32 explode string to slice
func ExplodeUint32(str string, seperator string) (result []uint32) {
	var strMids = strings.Split(str, seperator)
	for _, v := range strMids {
		mid, e := strconv.ParseInt(strings.Trim(v, trimSet), 10, 64)
		if e != nil {
			continue
		}
		result = append(result, uint32(mid))
	}
	return
}

//GetNextPeriodTime get next period time from current time
// clock, like "03:05:00"
// period, like 24h
// currentTime, like now
// return the next alarm time for this clock
func GetNextPeriodTime(clock string, period time.Duration, currentTime time.Time) (next time.Time, err error) {
	var now = currentTime
	var startTime, e = time.Parse("15:04:05", clock)
	err = e
	if err != nil {
		log.Error("clock is not right, config=%s, should like '12:00:00'")
		startTime = time.Date(2000, 1, 1, 3, 0, 0, 0, now.Location())
	}
	next = time.Date(now.Year(), now.Month(), now.Day(), startTime.Hour(), startTime.Minute(), startTime.Second(), 0, now.Location())
	d := next.Sub(now)
	for d < 0 {
		next = next.Add(period)
		d = next.Sub(now)
	}
	for d > period {
		next = next.Add(-period)
		d = next.Sub(now)
	}
	return
}

//TruncateDate 截取到整天，舍去时分秒
func TruncateDate(tm time.Time) time.Time {
	var y, m, d = tm.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, tm.Location())
}
