package dao

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoauthorize(t *testing.T) {
	convey.Convey("authorize", t, func(ctx convey.C) {
		var (
			key    = ""
			secret = ""
			method = ""
			bucket = ""
			expire = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			authorization := authorize(key, secret, method, bucket, expire)
			ctx.Convey("Then authorization should not be nil.", func(ctx convey.C) {
				ctx.So(authorization, convey.ShouldNotBeNil)
			})
		})
	})
}
