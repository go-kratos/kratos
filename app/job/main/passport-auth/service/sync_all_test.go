package service

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_compareToken(t *testing.T) {
	once.Do(startService)
	Convey("Test compareToken", t, func() {
		//s.compareToken("201806", 0, 10)
	})
}

func Test_compareCookie(t *testing.T) {
	once.Do(startService)
	Convey("Test compareCookie", t, func() {
		//s.compareCookie("201806", 0)
	})
}
