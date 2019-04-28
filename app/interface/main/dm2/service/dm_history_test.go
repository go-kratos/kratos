package service

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSearchDMHisIndex(t *testing.T) {
	Convey("test dm history date index", t, func() {
		res, err := svr.SearchDMHisIndex(context.TODO(), 1, 10109227, "2018-04")
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
		t.Log(res)
	})
}

func TestSearchDMHistory(t *testing.T) {
	Convey("", t, func() {
		date, _ := time.Parse("2006-01-02", "2018-04-24")
		// convert 2006-01-02-->2016-01-02 23:59:59
		tm := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, time.Local)
		xml, err := svr.SearchDMHistory(context.TODO(), 1, 10109227, tm.Unix())
		So(err, ShouldBeNil)
		t.Log(xml)
	})
}
