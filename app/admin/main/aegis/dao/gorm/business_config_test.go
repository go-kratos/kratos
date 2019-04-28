package gorm

import (
	"context"
	"go-common/app/admin/main/aegis/model/business"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestGormGetConfigs(t *testing.T) {
	convey.Convey("GetConfigs", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			bizid = int64(1)
		)
		ctx.Convey("success", func(ctx convey.C) {
			cfgs, err := d.GetConfigs(c, bizid)
			ctx.Convey("Then err should be nil.cfgs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(cfgs, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("empty", func(ctx convey.C) {
			_, err := d.GetConfigs(c, -1)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestGormGetConfig(t *testing.T) {
	convey.Convey("GetConfig", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			bizid = int64(0)
			tp    = int8(-1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			config, err := d.GetConfig(c, bizid, tp)
			ctx.Convey("Then err should be nil.config should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(config, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestGormActiveConfigs(t *testing.T) {
	convey.Convey("ActiveConfigs", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			configs, err := d.ActiveConfigs(c)
			ctx.Convey("Then err should be nil.configs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(configs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestGormAddBizConfig(t *testing.T) {
	convey.Convey("AddBizConfig", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			cfg = &business.BizCFG{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			lastid, err := d.AddBizConfig(c, cfg)
			ctx.Convey("Then err should be nil.lastid should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(lastid, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestGormEditBizConfig(t *testing.T) {
	convey.Convey("EditBizConfig", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			cfg = &business.BizCFG{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.EditBizConfig(c, cfg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
