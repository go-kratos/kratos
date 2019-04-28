package gorm

import (
	"context"
	"testing"

	"go-common/app/admin/main/aegis/model/business"

	"github.com/smartystreets/goconvey/convey"
)

func TestEnableBusiness(t *testing.T) {
	convey.Convey("EnableBusiness", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.EnableBusiness(c, 0)
			ctx.Convey("Then err should be nil.cfgs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDisableBusiness(t *testing.T) {
	convey.Convey("DisableBusiness", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DisableBusiness(c, 0)
			ctx.Convey("Then err should be nil.cfgs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestBusiness(t *testing.T) {
	convey.Convey("Business", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.Business(c, 0)
			ctx.Convey("Then err should be nil.cfgs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestBusinessList(t *testing.T) {
	convey.Convey("BusinessList", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.BusinessList(c, 1, []int64{0}, true)
			ctx.Convey("Then err should be nil.cfgs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestUpdateBusiness(t *testing.T) {
	convey.Convey("UpdateBusiness", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.UpdateBusiness(c, &business.Business{})
			ctx.Convey("Then err should be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
