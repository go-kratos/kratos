package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoCookie(t *testing.T) {
	var (
		c          = context.TODO()
		sd         = []byte("5ebd8ebb,1530838806,2c8d0678")
		sdNotExist = []byte("9f1c9145,1536117849,c9fb62a2")
		ct, _      = time.Parse("01/02/2006", "07/27/2018")
	)
	convey.Convey("Cookie", t, func(ctx convey.C) {
		res, session, err := d.Cookie(c, sd, ct)
		ctx.Convey("Then err should be nil.res,session should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(session, convey.ShouldNotBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
		res2, session2, err2 := d.Cookie(c, sdNotExist, ct)
		ctx.Convey("Then err should be nil.res,session should be nil.", func(ctx convey.C) {
			ctx.So(err2, convey.ShouldBeNil)
			ctx.So(session2, convey.ShouldBeNil)
			ctx.So(res2, convey.ShouldBeNil)
		})
	})
}
