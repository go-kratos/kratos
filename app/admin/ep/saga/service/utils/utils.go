package utils

import (
	"bytes"
	"sort"
	"strconv"
	"time"

	"go-common/app/admin/ep/saga/conf"
)

// CalAverageTime 计算时间的平均/*分位 时间
func CalAverageTime(timeType int, timeArray []float64) (result float64) {
	sort.Float64s(timeArray)
	if timeType == 0 {
		// 判断数量为零
		if len(timeArray) == 0 {
			result = 0
		} else {
			var sum float64
			for _, t := range timeArray {
				sum += t
			}
			result = sum / float64(len(timeArray))
		}
	} else if timeType > 0 && timeType < 11 {
		if len(timeArray) == 0 {
			result = 0
		} else {
			index := len(timeArray) * timeType / 10
			result = timeArray[index]
		}
	}
	return
}

// CalSizeTime ...
func CalSizeTime(time, max, min float64) (float64, float64) {
	if max == 0 {
		max = time
		return max, min
	}
	if min == 0 {
		min = time
		return max, min
	}

	if time > max {
		max = time
		return max, min
	}

	if time != 0 && time < min {
		min = time
		return max, min
	}

	return max, min
}

// CalSyncTime ...
func CalSyncTime() (since, until *time.Time) {
	syncDays := conf.Conf.Property.SyncData.DefaultSyncDays

	year, month, day := time.Now().Date()

	untilTime := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	sinceTime := time.Date(year, month, day-syncDays, 0, 0, 0, 0, time.Local)

	since = &sinceTime
	until = &untilTime

	return
}

// CombineSlice ...
func CombineSlice(s1, s2 []float64) []float64 {
	slice := make([]float64, len(s1)+len(s2))
	copy(slice, s1)
	copy(slice[len(s1):], s2)
	return slice
}

// InSlice ...
func InSlice(key interface{}, list []string) bool {
	for _, item := range list {
		if key == item {
			return true
		}
	}
	return false
}

// Unicode2Chinese ...
func Unicode2Chinese(str string) string {
	buf := bytes.NewBuffer(nil)
	i, j := 0, len(str)
	for i < j {
		x := i + 6
		if x > j {
			buf.WriteString(str[i:])
			break
		}
		if str[i] == '\\' && str[i+1] == 'u' {
			hex := str[i+2 : x]
			r, err := strconv.ParseUint(hex, 16, 64)
			if err == nil {
				buf.WriteRune(rune(r))
			} else {
				buf.WriteString(str[i:x])
			}
			i = x
		} else {
			buf.WriteByte(str[i])
			i++
		}
	}
	return buf.String()
}
