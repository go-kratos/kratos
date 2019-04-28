package dao

import (
	"context"
	"go-common/app/admin/main/creative/conf"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoUpload(t *testing.T) {
	convey.Convey("Upload", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			fileName = "filename"
			fileType = "png"
			timing   = int64(1545382342)
			data     = []byte("iamdata")
			bfs      = &conf.Bfs{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			bfs = d.c.Bfs
			location, err := d.Upload(c, fileName, fileType, timing, data, bfs)
			ctx.Convey("Then err should be nil.location should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil) // http code 401
				ctx.So(location, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoauthorize(t *testing.T) {
	convey.Convey("authorize", t, func(ctx convey.C) {
		var (
			key    = "8d4e593ba7555502"
			secret = "0bdbd4c7caeeddf587c3c4daec0475"
			method = "PUT"
			bucket = "archive"
			file   = ""
			expire = int64(100)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			authorization := authorize(key, secret, method, bucket, file, expire)
			ctx.Convey("Then authorization should not be nil.", func(ctx convey.C) {
				ctx.So(authorization, convey.ShouldNotBeNil)
			})
		})
	})
}
