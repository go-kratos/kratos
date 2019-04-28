package mc

import (
	"context"
	"testing"

	"go-common/app/admin/main/aegis/model/common"

	"github.com/smartystreets/goconvey/convey"
)

func TestMcConsumerOn(t *testing.T) {
	convey.Convey("ConsumerOn", t, func(ctx convey.C) {
		var (
			opt = &common.BaseOptions{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.ConsumerOn(context.TODO(), opt)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMcConsumerOff(t *testing.T) {
	convey.Convey("ConsumerOff", t, func(ctx convey.C) {
		var (
			opt = &common.BaseOptions{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.ConsumerOff(context.TODO(), opt)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMcIsConsumerOn(t *testing.T) {
	convey.Convey("IsConsumerOn", t, func(ctx convey.C) {
		var (
			opt = &common.BaseOptions{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.IsConsumerOn(context.TODO(), opt)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMcmcKey(t *testing.T) {
	convey.Convey("mcKey", t, func(ctx convey.C) {
		var (
			opt = &common.BaseOptions{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := mcKey(opt)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMcroleKey(t *testing.T) {
	convey.Convey("roleKey", t, func(ctx convey.C) {
		var (
			opt = &common.BaseOptions{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := roleKey(opt)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMcSetRole(t *testing.T) {
	convey.Convey("SetRole", t, func(ctx convey.C) {
		var (
			c   = context.TODO()
			opt = &common.BaseOptions{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SetRole(c, opt, 0)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMcGetRole(t *testing.T) {
	convey.Convey("GetRole", t, func(ctx convey.C) {
		var (
			c   = context.TODO()
			opt = &common.BaseOptions{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, _, err := d.GetRole(c, opt)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
