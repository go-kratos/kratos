package util

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

var periodTimeTest = [][]string{
	// time, period, now, expected
	{"11:00:00", "1h", "10:00:00", "0000-01-01 11:00:00"},
	{"11:00:00", "1m", "10:00:00", "0000-01-01 10:01:00"},
	{"11:00:00", "1s", "10:00:00", "0000-01-01 10:00:01"},

	{"11:00:00", "8h", "10:00:00", "0000-01-01 11:00:00"},
	{"11:00:00", "12h", "10:00:00", "0000-01-01 11:00:00"},
	{"11:00:00", "23h", "10:00:00", "0000-01-01 11:00:00"},

	{"09:00:10", "1h", "10:00:00", "0000-01-01 10:00:10"},
	{"09:00:10", "1m", "10:00:00", "0000-01-01 10:00:10"},
	{"09:00:10", "1s", "10:00:00", "0000-01-01 10:00:00"},

	{"09:00:10", "8h", "10:00:00", "0000-01-01 17:00:10"},
	{"09:00:10", "12h", "10:00:00", "0000-01-01 21:00:10"},
	{"09:00:10", "4h", "10:00:00", "0000-01-01 13:00:10"},

	{"09:00:10", "24h", "10:00:00", "0000-01-02 09:00:10"},
	{"09:00:10", "8h", "10:00:00", "0000-01-01 17:00:10"},
}

func TestGetNextPeriodTime(t *testing.T) {
	Convey("test get period", t, func() {
		for i, v := range periodTimeTest {
			//needTime, _ := time.Parse("15:04:05", v[0])
			period, _ := time.ParseDuration(v[1])
			now, _ := time.Parse("15:04:05", v[2])
			expected, _ := time.Parse("2006-01-02 15:04:05", v[3])
			actual, err := GetNextPeriodTime(v[0], period, now)
			So(err, ShouldEqual, nil)
			t.Logf("[%d]actual:+%v, expected:+%v", i, actual, expected)
			So(actual.Equal(expected), ShouldEqual, true)
		}
	})
}

func TestUnSetBit64(t *testing.T) {
	var (
		testCase = [][]int64{
			// attr, bit, result
			{1, 0, 0},
			{2, 0, 2},
			{2, 1, 0},
			{3, 1, 1},
			{3, 64, 3},
		}
	)

	Convey("test for unset bit 64", t, func() {
		for _, v := range testCase {
			So(UnSetBit64(v[0], uint(v[1])), ShouldEqual, v[2])
		}
	})
}
