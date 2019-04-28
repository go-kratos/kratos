package common

import (
	"database/sql/driver"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// TimeFormat time format
var TimeFormat = "2006-01-02 15:04:05"

//GrayField gray config for each business
type GrayField struct {
	Name  string
	Value string
}

// Pager .
type Pager struct {
	Total int `json:"total" reflect:"ignore"`
	Pn    int `form:"pn" default:"1" json:"pn" reflect:"ignore"`
	Ps    int `form:"ps" default:"20" json:"ps" reflect:"ignore"`
}

// BaseOptions 公共参数
type BaseOptions struct {
	BusinessID int64  `form:"business_id" json:"business_id"`
	NetID      int64  `form:"net_id" json:"net_id"`
	FlowID     int64  `form:"flow_id" json:"flow_id"`
	UID        int64  `form:"uid" json:"uid" submit:"int"`
	OID        string `form:"oid" json:"oid" submit:"string"`
	RID        int64  `form:"rid" json:"rid"`
	Role       int8   `form:"role" json:"role"`
	Debug      int8   `form:"debug" json:"debug"`
	Uname      string `form:"uname" json:"uname" submit:"string"`
}

// FormatTime .
type FormatTime string

// Scan .
func (f *FormatTime) Scan(src interface{}) (err error) {
	switch sc := src.(type) {
	case time.Time:
		*f = FormatTime(sc.Format("2006-01-02 15:04:05"))
	case string:
		*f = FormatTime(sc)
	}
	return
}

// WaitTime 计算等待时长
func WaitTime(ctime time.Time) string {
	wt := time.Since(ctime)
	h := int(wt.Hours())
	m := int(wt.Minutes()) % 60
	s := int(wt.Seconds()) % 60
	return fmt.Sprintf("%.2d:%.2d:%.2d", h, m, s)
}

//ParseWaitTime 。
func ParseWaitTime(ut int64) string {
	h := ut / 3600
	m := ut % 3600 / 60
	s := ut % 60
	return fmt.Sprintf("%.2d:%.2d:%.2d", h, m, s)
}

// Group .
type Group struct {
	ID        int64  `json:"group_id"`
	Name      string `json:"group_name"`
	Note      string `json:"group_note"`
	Tag       string `json:"group_tag"`
	FontColor string `json:"font_color"`
	BgColor   string `json:"bg_color"`
}

// IntTime .
type IntTime int64

// Scan scan time.
func (jt *IntTime) Scan(src interface{}) (err error) {
	switch sc := src.(type) {
	case time.Time:
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
	return time.Unix(int64(jt), 0), nil
}

// Time get time.
func (jt IntTime) Time() time.Time {
	return time.Unix(int64(jt), 0)
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
	st, err := time.ParseInLocation("2006-01-02 15:04:05", str, time.Local)
	if err == nil {
		*jt = IntTime(st.Unix())
		return nil
	}

	*jt = IntTime(0)

	return nil
}

// FilterName .
func FilterName(s string) (res string) {
	exp := "[^a-zA-Z0-9_]+"
	reg, err := regexp.Compile(exp)
	if err != nil {
		res = s
		return
	}

	res = reg.ReplaceAllString(s, "")
	return
}

// FilterChname .
func FilterChname(s string) (res string) {
	exp := "[^0-9_\u4e00-\u9fa5]+"
	reg, err := regexp.Compile(exp)
	if err != nil {
		res = s
		return
	}

	res = reg.ReplaceAllString(s, "")
	return
}

// FilterBusinessName .
func FilterBusinessName(s string) (res string) {
	exp := "[^a-zA-Z\u4e00-\u9fa5]+"
	reg, err := regexp.Compile(exp)
	if err != nil {
		res = s
		return
	}

	res = reg.ReplaceAllString(s, "")
	return
}

//Unique remove duplicated value from slice
func Unique(ids []int64, gthan0 bool) (res []int64) {
	res = []int64{}
	mm := map[int64]int64{}
	for _, id := range ids {
		if mm[id] == id || (gthan0 && id <= 0) {
			continue
		}

		res = append(res, id)
		mm[id] = id
	}
	return
}

//CopyMap copy src to dest
func CopyMap(src, dest map[int64][]int64, gthan0 bool) (res map[int64][]int64) {
	if dest == nil {
		dest = map[int64][]int64{}
	}

	for k, v := range src {
		dest[k] = append(dest[k], v...)
		dest[k] = Unique(dest[k], gthan0)
	}

	res = dest
	return
}
