package mcndao

import (
	"context"
	"testing"
	"time"

	"go-common/app/interface/main/mcn/model/mcnmodel"

	"github.com/smartystreets/goconvey/convey"
)

func TestMcndaoAddCacheMcnSign(t *testing.T) {
	convey.Convey("AddCacheMcnSign", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(0)
			val = &mcnmodel.McnSign{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheMcnSign(c, id, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMcndaoCacheMcnSign(t *testing.T) {
	convey.Convey("CacheMcnSign", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheMcnSign(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMcndaoDelCacheMcnSign(t *testing.T) {
	convey.Convey("DelCacheMcnSign", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelCacheMcnSign(c, id)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMcndaoAddCacheMcnDataSummary(t *testing.T) {
	convey.Convey("AddCacheMcnDataSummary", t, func(ctx convey.C) {
		var (
			c            = context.Background()
			id           = int64(0)
			val          = &mcnmodel.McnGetDataSummaryReply{}
			generateDate = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheMcnDataSummary(c, id, val, generateDate)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMcndaoCacheMcnDataSummary(t *testing.T) {
	convey.Convey("CacheMcnDataSummary", t, func(ctx convey.C) {
		var (
			c            = context.Background()
			id           = int64(0)
			generateDate = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheMcnDataSummary(c, id, generateDate)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMcndaoDelMcnDataSummary(t *testing.T) {
	convey.Convey("DelMcnDataSummary", t, func(ctx convey.C) {
		var (
			c            = context.Background()
			id           = int64(0)
			generateDate = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelMcnDataSummary(c, id, generateDate)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
