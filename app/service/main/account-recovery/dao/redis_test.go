package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSetLinkMailCount(t *testing.T) {
	var (
		c        = context.Background()
		linkMail = "2459593393@qq.com"
	)
	convey.Convey("SetLinkMailCount", t, func(ctx convey.C) {
		state, err := d.SetLinkMailCount(c, linkMail)
		ctx.Convey("Then err should be nil.state should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(state, convey.ShouldNotBeNil)
		})
	})
}

func TestDaogetSubtime(t *testing.T) {
	convey.Convey("getSubtime", t, func(ctx convey.C) {
		subtime := getSubtime()
		ctx.Convey("Then subtime should not be nil.", func(ctx convey.C) {
			ctx.So(subtime, convey.ShouldNotBeNil)
		})
	})
}

func TestDaokeyCaptcha(t *testing.T) {
	var (
		mid      = int64(0)
		linkMail = "2459593393@qq.com"
	)
	convey.Convey("keyCaptcha", t, func(ctx convey.C) {
		p1 := keyCaptcha(mid, linkMail)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSetCaptcha(t *testing.T) {
	var (
		c        = context.Background()
		code     = "1234"
		mid      = int64(1)
		linkMail = "2459593393@qq.com"
	)
	convey.Convey("SetCaptcha", t, func(ctx convey.C) {
		err := d.SetCaptcha(c, code, mid, linkMail)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoGetEMailCode(t *testing.T) {
	var (
		c        = context.Background()
		mid      = int64(1)
		linkMail = "2459593393@qq.com"
	)
	convey.Convey("GetEMailCode", t, func(ctx convey.C) {
		code, err := d.GetEMailCode(c, mid, linkMail)
		ctx.Convey("Then err should be nil.code should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(code, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelEMailCode(t *testing.T) {
	var (
		c        = context.Background()
		mid      = int64(1)
		linkMail = "2459593393@qq.com"
	)
	convey.Convey("GetEMailCode", t, func(ctx convey.C) {
		err := d.DelEMailCode(c, mid, linkMail)
		ctx.Convey("Then err should be nil.code should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoPingRedis(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("PingRedis", t, func(ctx convey.C) {
		err := d.PingRedis(c)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
