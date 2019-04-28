package dao

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Categories(t *testing.T) {
	Convey("should get data", t, func() {
		res, err := d.Categories(context.TODO())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func Test_ArticlesStats(t *testing.T) {
	Convey("get data", t, func() {
		res, err := d.ArticlesStats(context.TODO(), []int64{1})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})

	Convey("no data", t, func() {
		res, err := d.ArticlesStats(context.TODO(), []int64{100000})
		So(err, ShouldBeNil)
		So(res, ShouldBeNil)
	})
}
func Test_AddComplaint(t *testing.T) {
	Convey("add data", t, func() {
		err := d.AddComplaint(context.TODO(), 1, 2, 3, "reason", "http://1.ipg")
		So(err, ShouldBeNil)
	})
}

func Test_Notices(t *testing.T) {
	Convey("get data", t, func() {
		t := time.Unix(1513322993, 0)
		res, err := d.Notices(context.TODO(), t)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func Test_NoticeState(t *testing.T) {
	Convey("add data", t, func() {
		mid := int64(100)
		state := int64(1)
		err := d.UpdateNoticeState(context.TODO(), mid, state)
		So(err, ShouldBeNil)
		Convey("get data", func() {
			res, err := d.NoticeState(context.TODO(), mid)
			So(err, ShouldBeNil)
			So(res, ShouldEqual, state)
		})
	})
}

func Test_Hotspot(t *testing.T) {
	Convey("should get data", t, func() {
		res, err := d.Hotspots(context.TODO())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func Test_SearchArts(t *testing.T) {
	Convey("should get data", t, func() {
		_searchInterval = 24 * 3600 * 365
		res, err := d.SearchArts(context.TODO(), 0)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func Test_CheatFilter(t *testing.T) {
	Convey("add data", t, func() {
		err := d.AddCheatFilter(context.TODO(), 100, 2)
		So(err, ShouldBeNil)
		err = d.DelCheatFilter(context.TODO(), 100)
		So(err, ShouldBeNil)
	})
}
