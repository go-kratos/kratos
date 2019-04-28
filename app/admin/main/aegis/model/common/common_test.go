package common

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestFormatTime_Scan(t *testing.T) {
	convey.Convey("FormatTime_Scan", t, func(ctx convey.C) {
		a := FormatTime("2018-01-01 01:00:00")
		b := time.Time{}
		err := a.Scan(b)
		ctx.Convey("FormatTime_Scan", func(ctx convey.C) {
			convey.So(err, convey.ShouldBeNil)
		})
	})
}

func TestWaitTime(t *testing.T) {
	convey.Convey("WaitTime", t, func(ctx convey.C) {
		b := time.Time{}
		WaitTime(b)
	})
}

func TestParseWaitTime(t *testing.T) {
	convey.Convey("ParseWaitTime", t, func(ctx convey.C) {
		b := time.Time{}.Unix()
		ParseWaitTime(b)
	})
}

var s = IntTime(time.Now().Unix())

func TestIntTime_Scan(t *testing.T) {
	convey.Convey("IntTime_Scan", t, func(ctx convey.C) {
		b := time.Time{}
		err := s.Scan(&b)
		ctx.Convey("IntTime_Scan", func(ctx convey.C) {
			convey.So(err, convey.ShouldBeNil)
		})
	})
}

func TestIntTime_Value(t *testing.T) {
	convey.Convey("IntTime_Value", t, func(ctx convey.C) {
		s.Value()
	})
}

func TestIntTime_UnmarshalJSON(t *testing.T) {
	convey.Convey("IntTime_UnmarshalJSON", t, func(ctx convey.C) {
		err := s.UnmarshalJSON(nil)
		ctx.Convey("IntTime_UnmarshalJSON", func(ctx convey.C) {
			convey.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFilterName(t *testing.T) {
	convey.Convey("FilterName", t, func(ctx convey.C) {
		a := FilterName("hahah一个")
		convey.So(a, convey.ShouldEqual, "hahah")
	})
}

func TestFilterChname(t *testing.T) {
	convey.Convey("FilterChname", t, func(ctx convey.C) {
		a := FilterChname("hahah一个")
		convey.So(a, convey.ShouldEqual, "一个")
	})
}

func TestFilterBusinessName(t *testing.T) {
	convey.Convey("FilterBusinessName", t, func(ctx convey.C) {
		a := FilterBusinessName("hahah一个")
		convey.So(a, convey.ShouldEqual, "hahah一个")
	})
}

func TestUnique(t *testing.T) {
	convey.Convey("Unique", t, func(ctx convey.C) {
		b := []int64{1, -1, 0, 1, 1}
		a := Unique(b, true)
		convey.So(len(a), convey.ShouldEqual, 1)
		convey.So(a[0], convey.ShouldEqual, 1)
	})
}

func TestCopyMap(t *testing.T) {
	convey.Convey("CopyMap", t, func(ctx convey.C) {
		b := map[int64][]int64{}
		a := CopyMap(map[int64][]int64{1: []int64{1}}, b, true)
		convey.So(len(a), convey.ShouldEqual, 1)
	})
}
