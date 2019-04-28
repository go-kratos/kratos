package tag

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestTagInsertActivityInfo(t *testing.T) {
	convey.Convey("InsertActivityInfo", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			vals = "(1,2,3,4,5)"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.InsertActivityInfo(c, vals)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestTagListActivityInfo(t *testing.T) {
	convey.Convey("ListActivityInfo", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			avs, err := d.ListActivityInfo(c)
			ctx.Convey("Then err should be nil.avs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(avs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestTagGetVideoActivityInfo(t *testing.T) {
	convey.Convey("GetVideoActivityInfo", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			activities = []int64{}
			pn         = int(0)
			ps         = int(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			info, err := d.GetVideoActivityInfo(c, activities, pn, ps)
			ctx.Convey("Then err should not be nil.info should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(info, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestTagGetCmActivityInfo(t *testing.T) {
	convey.Convey("GetCmActivityInfo", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			activityID = int64(0)
			pn         = int(0)
			ps         = int(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			info, err := d.GetCmActivityInfo(c, activityID, pn, ps)
			ctx.Convey("Then err should not be nil.info should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(info, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestTagGetVideoTypes(t *testing.T) {
	convey.Convey("GetVideoTypes", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rmap, err := d.GetVideoTypes(c)
			ctx.Convey("Then err should not be nil.rmap should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rmap, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestTagGetColumnTypes(t *testing.T) {
	convey.Convey("GetColumnTypes", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rmap, err := d.GetColumnTypes(c)
			ctx.Convey("Then err should be nil.rmap should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rmap, convey.ShouldNotBeNil)
			})
		})
	})
}
