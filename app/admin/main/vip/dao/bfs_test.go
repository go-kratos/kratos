package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoUpload(t *testing.T) {
	convey.Convey("Upload", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			fileName = "123"
			fileType = "image/png"
			data     = []byte("1231")
			bfs      = d.c.Bfs
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			location, err := d.Upload(c, fileName, fileType, data, bfs)
			ctx.Convey("Then err should be nil.location should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(location, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoauthorize(t *testing.T) {
	convey.Convey("authorize", t, func(ctx convey.C) {
		var (
			key    = ""
			secret = ""
			method = ""
			bucket = ""
			file   = ""
			expire = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			authorization := authorize(key, secret, method, bucket, file, expire)
			ctx.Convey("Then authorization should not be nil.", func(ctx convey.C) {
				ctx.So(authorization, convey.ShouldNotBeNil)
			})
		})
	})
}
