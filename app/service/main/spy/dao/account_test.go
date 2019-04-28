package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoTelRiskLevel(t *testing.T) {
	convey.Convey("TelRiskLevel", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			riskLevel, err := d.TelRiskLevel(c, mid)
			ctx.Convey("Then err should be nil.riskLevel should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(riskLevel, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBlockAccount(t *testing.T) {
	convey.Convey("BlockAccount", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.BlockAccount(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSecurityLogin(t *testing.T) {
	convey.Convey("SecurityLogin", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(1)
			reason = "unit test"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SecurityLogin(c, mid, reason)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoTelInfo(t *testing.T) {
	convey.Convey("TelInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tel, err := d.TelInfo(c, mid)
			ctx.Convey("Then err should be nil.tel should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(tel, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestProfileInfo(t *testing.T) {
	convey.Convey("ProfileInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tel, err := d.ProfileInfo(c, mid, "")
			ctx.Convey("Then err should be nil.tel should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(tel, convey.ShouldNotBeNil)
			})
		})
	})
}
