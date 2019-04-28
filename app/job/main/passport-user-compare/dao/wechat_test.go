package dao

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSendWechat(t *testing.T) {
	convey.Convey("SendWechat", t, func(ctx convey.C) {
		var (
			param map[string]int64
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SendWechat(param)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaosign(t *testing.T) {
	convey.Convey("sign", t, func(ctx convey.C) {
		var (
			params map[string]string
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := d.sign(params)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
