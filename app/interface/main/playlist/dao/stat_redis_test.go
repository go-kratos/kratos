package dao

import (
	"context"
	"testing"

	"go-common/app/interface/main/playlist/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyStat(t *testing.T) {
	var (
		mid = int64(2)
	)
	convey.Convey("keyStat", t, func(ctx convey.C) {
		p1 := keyStat(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaokeyPl(t *testing.T) {
	var (
		pid = int64(1)
	)
	convey.Convey("keyPl", t, func(ctx convey.C) {
		p1 := keyPl(pid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoPlStatCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(2)
		pid = int64(1)
	)
	convey.Convey("PlStatCache", t, func(ctx convey.C) {
		_, err := d.PlStatCache(c, mid, pid)
		ctx.Convey("Then err should be nil.stat should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSetPlStatCache(t *testing.T) {
	var (
		c    = context.Background()
		mid  = int64(2)
		pid  = int64(1)
		stat = &model.PlStat{}
	)
	convey.Convey("SetPlStatCache", t, func(ctx convey.C) {
		err := d.SetPlStatCache(c, mid, pid, stat)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSetStatsCache(t *testing.T) {
	var (
		c       = context.Background()
		mid     = int64(2)
		plStats = []*model.PlStat{}
	)
	plStats = append(plStats, &model.PlStat{ID: 1, Share: 1})
	convey.Convey("SetStatsCache", t, func(ctx convey.C) {
		err := d.SetStatsCache(c, mid, plStats)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoPlsCache(t *testing.T) {
	var (
		c    = context.Background()
		pids = []int64{1, 2, 3}
	)
	convey.Convey("PlsCache", t, func(ctx convey.C) {
		res, err := d.PlsCache(c, pids)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSetPlCache(t *testing.T) {
	var (
		c       = context.Background()
		plStats = []*model.PlStat{}
	)
	plStats = append(plStats, &model.PlStat{ID: 1, View: 100})
	convey.Convey("SetPlCache", t, func(ctx convey.C) {
		err := d.SetPlCache(c, plStats)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelPlCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(2)
		pid = int64(1)
	)
	convey.Convey("DelPlCache", t, func(ctx convey.C) {
		err := d.DelPlCache(c, mid, pid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoStatsCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(2)
	)
	convey.Convey("StatsCache", t, func(ctx convey.C) {
		_, err := d.StatsCache(c, mid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
