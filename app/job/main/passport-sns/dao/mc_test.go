package dao

import (
	"context"
	"testing"

	"go-common/app/job/main/passport-sns/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDao_snsKey(t *testing.T) {
	convey.Convey("snsKey", t, func(ctx convey.C) {
		var (
			platform = model.PlatformQQStr
			mid      = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res := snsKey(platform, mid)
			ctx.Convey("Then res should not be nil.", func(ctx convey.C) {
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDao_SetSnsCache(t *testing.T) {
	convey.Convey("SetSnsCache", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(0)
			platform = model.PlatformQQStr
			qq       = &model.SnsProto{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SetSnsCache(c, mid, platform, qq)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDao_DelSnsCache(t *testing.T) {
	convey.Convey("DelSnsCache", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(0)
			platform = model.PlatformQQStr
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelSnsCache(c, mid, platform)
			ctx.Convey("Then err should be nil.info should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDao_GetUnionIDCache(t *testing.T) {
	convey.Convey("GetUnionIDCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			key = "test"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.GetUnionIDCache(c, key)
			ctx.Convey("Then err should be nil.res should be nil.", func(ctx convey.C) {
				ctx.So(res, convey.ShouldBeEmpty)
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
