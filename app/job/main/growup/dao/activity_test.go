package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetCActivities(t *testing.T) {
	convey.Convey("GetCActivities", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			acs, err := d.GetCActivities(c)
			ctx.Convey("Then err should be nil.acs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(acs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoListUpActivity(t *testing.T) {
	convey.Convey("ListUpActivity", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ups, err := d.ListUpActivity(c, id)
			ctx.Convey("Then err should be nil.ups should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ups, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetActivityBonus(t *testing.T) {
	convey.Convey("GetActivityBonus", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			actBonus, err := d.GetActivityBonus(c, id)
			ctx.Convey("Then err should be nil.actBonus should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(actBonus, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetAvUploadByMID(t *testing.T) {
	convey.Convey("GetAvUploadByMID", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(10)
			limit = int(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			avs, err := d.GetAvUploadByMID(c, id, limit)
			ctx.Convey("Then err should be nil.avs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(avs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetArchiveInfo(t *testing.T) {
	convey.Convey("GetArchiveInfo", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			activityID = int64(10)
			id         = int64(100)
			limit      = int(200)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			avs, err := d.GetArchiveInfo(c, activityID, id, limit)
			ctx.Convey("Then err should be nil.avs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(avs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateUpActivityState(t *testing.T) {
	convey.Convey("UpdateUpActivityState", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			id       = int64(100)
			oldState = int(10)
			newState = int(200)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.UpdateUpActivityState(c, id, oldState, newState)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertUpActivityBatch(t *testing.T) {
	convey.Convey("InsertUpActivityBatch", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			vals = "(10, 'test', 100, 90, 100, 1, 1000, 3, '2018-06-23')"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.InsertUpActivityBatch(c, vals)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}
