package dao

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoBlacklistAdd(t *testing.T) {
	convey.Convey("BlacklistAdd", t, func(ctx convey.C) {
		var (
			mids = []int64{10000, 10001}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.BlacklistAdd(mids, mids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoBlacklistIn(t *testing.T) {
	convey.Convey("BlacklistIn", t, func(ctx convey.C) {
		var (
			mids = []int64{1, 2, 3, 4, 5, 6, 7, 8}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			blacklist, err := d.BlacklistIn(mids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(blacklist, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBlacklistUp(t *testing.T) {
	convey.Convey("BlacklistUp", t, func(ctx convey.C) {
		var (
			id     = int64(1)
			status = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.BlacklistUp(id, status)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoBlacklistIndex(t *testing.T) {
	convey.Convey("BlacklistIndex", t, func(ctx convey.C) {
		var (
			mid = int64(1)
			pn  = int(1)
			ps  = int(10)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			pager, err := d.BlacklistIndex(mid, pn, ps)
			ctx.Convey("Then err should be nil.pager should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(pager, convey.ShouldNotBeNil)
			})
		})
	})
}
