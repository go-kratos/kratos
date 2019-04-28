package v2

import (
	"context"
	"go-common/app/infra/config/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestV2Hosts(t *testing.T) {
	var (
		c   = context.TODO()
		svr = "11133_fat1_sh001"
	)
	convey.Convey("Hosts", t, func(ctx convey.C) {
		hosts, err := d.Hosts(c, svr)
		ctx.Convey("Then err should be nil.hosts should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(hosts, convey.ShouldNotBeNil)
		})
	})
}

func TestV2SetHost(t *testing.T) {
	var (
		c    = context.TODO()
		host = &model.Host{}
		svr  = "11133_fat1_sh001"
	)
	convey.Convey("SetHost", t, func(ctx convey.C) {
		err := d.SetHost(c, host, svr)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestV2ClearHost(t *testing.T) {
	var (
		c   = context.TODO()
		svr = "11133_fat1_sh001"
	)
	convey.Convey("ClearHost", t, func(ctx convey.C) {
		err := d.ClearHost(c, svr)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestV2pingRedis(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("pingRedis", t, func(ctx convey.C) {
		err := d.pingRedis(c)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
