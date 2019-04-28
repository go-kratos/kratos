package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/archive/api"
	feed "go-common/app/service/main/feed/model"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_ArchivesCache(t *testing.T) {
	arc := api.Arc{Aid: 1, PubDate: xtime.Time(100), Title: "title"}
	c := context.TODO()
	Convey("add cache", t, func() {
		err := d.AddArchivesCacheMap(c, map[int64]*api.Arc{1: &arc})
		So(err, ShouldBeNil)
		Convey("get cache return cached data", func() {
			cached, missed, err := d.ArchivesCache(c, []int64{1})
			So(err, ShouldBeNil)
			So(missed, ShouldBeEmpty)
			So(cached, ShouldResemble, map[int64]*api.Arc{1: &arc})
		})

		Convey("del cache return null", func() {
			err := d.DelArchiveCache(c, 1)
			So(err, ShouldBeNil)
			cached, missed, err := d.ArchivesCache(c, []int64{1})
			So(err, ShouldBeNil)
			So(cached, ShouldBeEmpty)
			So(missed, ShouldResemble, []int64{1})
		})
	})
}

func Test_BangumiCache(t *testing.T) {
	bangumi := feed.Bangumi{SeasonID: 1, Title: "t"}
	c := context.TODO()
	Convey("add cache", t, func() {
		err := d.AddBangumisCacheMap(c, map[int64]*feed.Bangumi{1: &bangumi})
		So(err, ShouldBeNil)
		Convey("get cache return cached data", func() {
			cached, missed, err := d.BangumisCache(c, []int64{1})
			So(err, ShouldBeNil)
			So(missed, ShouldBeEmpty)
			So(cached, ShouldResemble, map[int64]*feed.Bangumi{1: &bangumi})
		})

		Convey("return missed", func() {
			miss := int64(2000)
			cached, missed, err := d.BangumisCache(c, []int64{miss})
			So(err, ShouldBeNil)
			So(cached, ShouldBeEmpty)
			So(missed, ShouldResemble, []int64{miss})
		})
	})
}
