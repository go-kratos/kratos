package http

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	_outFormat = "2006-01-02 15:04:05"
)

func TestInit(t *testing.T) {
	Convey("range when from str equal to to str", t, func() {
		fromStr := "2017-12-25"
		toStr := "2017-12-25"
		from, to, ok := rangeDate(fromStr, toStr)
		So(ok, ShouldEqual, true)
		So(from.Format(_outFormat), ShouldEqual, "2017-12-25 00:00:00")
		So(to.Format(_outFormat), ShouldEqual, "2017-12-25 23:59:59")
		So(to.Unix()-from.Unix(), ShouldEqual, 86399)
	})
}
