package dao

import (
	"context"
	"go-common/app/service/main/relation/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoloadStat(t *testing.T) {
	var (
		ctx  = context.Background()
		mid  = int64(1)
		stat = &model.Stat{
			Mid:       1,
			Follower:  1,
			Following: 1,
			Black:     1,
			Whisper:   1,
		}
	)
	d.SetStatCache(ctx, mid, stat)
	d.storeStat(mid, stat)
	convey.Convey("loadStat", t, func(cv convey.C) {
		d.storeStat(mid, stat)
		cv.Convey("No return values", func(cv convey.C) {
		})

		s1, err := d.loadStat(ctx, mid)
		cv.Convey("loadStat", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(s1, convey.ShouldNotBeNil)
		})

		p1, p2, err := d.StatsCache(ctx, []int64{1})
		cv.Convey("StatsCache", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(p2, convey.ShouldNotBeNil)
			cv.So(p1, convey.ShouldNotBeNil)
		})
	})
}
