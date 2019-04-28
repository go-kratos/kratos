package dao

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSendWechat(t *testing.T) {
	msg := "push strategy test wechat message"
	convey.Convey("SendWechat", t, func(ctx convey.C) {
		err := d.SendWechat(msg)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaosign(t *testing.T) {
	params := map[string]string{"a": "b"}
	convey.Convey("sign", t, func(ctx convey.C) {
		p1 := d.sign(params)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeEmpty)
		})
	})
}
