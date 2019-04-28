package http

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestHttpauthorize(t *testing.T) {
	convey.Convey("authorize", t, func(ctx convey.C) {
		var (
			key    = ""
			secret = ""
			method = ""
			bucket = ""
			file   = ""
			expire = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			authorization := authorize(key, secret, method, bucket, file, expire)
			ctx.Convey("Then authorization should not be nil.", func(ctx convey.C) {
				ctx.So(authorization, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestHttpUpload(t *testing.T) {
	convey.Convey("Upload", t, func(ctx convey.C) {
		var (
			c = context.TODO()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.Upload(c, "", "", 0, []byte{})
			ctx.Convey("Then authorization should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
