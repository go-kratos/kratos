package model

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	xtime "time"
)

// WaitTime 计算等待时长
func WaitTime(ctime xtime.Time) string {
	wt := xtime.Since(ctime)
	h := int(wt.Hours())
	m := int(wt.Minutes()) % 60
	s := int(wt.Seconds()) % 60
	return fmt.Sprintf("%.2d:%.2d:%.2d", h, m, s)
}

//IntTime .
type IntTime int64

// Scan scan time.
func (jt *IntTime) Scan(src interface{}) (err error) {
	switch sc := src.(type) {
	case xtime.Time:
		*jt = IntTime(sc.Unix())
	case string:
		var i int64
		i, err = strconv.ParseInt(sc, 10, 64)
		*jt = IntTime(i)
	}
	return
}

// Value get time value.
func (jt IntTime) Value() (driver.Value, error) {
	return xtime.Unix(int64(jt), 0), nil
}

// Time get time.
func (jt IntTime) Time() xtime.Time {
	return xtime.Unix(int64(jt), 0)
}

// UnmarshalJSON implement Unmarshaler
func (jt *IntTime) UnmarshalJSON(data []byte) error {
	if data == nil || len(data) <= 1 {
		*jt = 0
		return nil
	}
	if data[0] != '"' {
		// 1.直接判断数字
		sti, err := strconv.Atoi(string(data))
		if err == nil {
			*jt = IntTime(sti)
		}
		return nil
	}

	str := string(data[1 : len(data)-1])

	// 2.标准格式判断
	st, err := xtime.ParseInLocation("2006-01-02 15:04:05", str, xtime.Local)
	if err == nil {
		*jt = IntTime(st.Unix())
		return nil
	}

	*jt = IntTime(0)

	return nil
}

//BaseResponse .
type BaseResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}
