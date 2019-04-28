package account

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAccountUploadImage(t *testing.T) {
	var (
		c        = context.Background()
		fileType = "jpg"
		bs       = []byte("dangerous")
	)
	convey.Convey("UploadImage", t, func(ctx convey.C) {
		location, err := d.UploadImage(c, fileType, bs, d.c.BFS)
		ctx.Convey("Then err should be nil.location should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(location, convey.ShouldNotBeNil)
		})
	})
}

func TestAccountauthorize(t *testing.T) {
	var (
		key    = ""
		secret = ""
		method = ""
		bucket = ""
		expire = int64(0)
	)
	convey.Convey("authorize", t, func(ctx convey.C) {
		authorization := authorize(key, secret, method, bucket, expire)
		ctx.Convey("Then authorization should not be nil.", func(ctx convey.C) {
			ctx.So(authorization, convey.ShouldNotBeNil)
		})
	})
}
