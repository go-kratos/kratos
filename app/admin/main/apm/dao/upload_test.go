package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoauthorize(t *testing.T) {
	convey.Convey("authorize", t, func(ctx convey.C) {
		var (
			key    = d.c.Bfs.Key
			secret = d.c.Bfs.Secret
			method = "PUT"
			bucket = d.c.Bfs.Bucket
			expire = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			authorization := authorize(key, secret, method, bucket, expire)
			ctx.Convey("Then authorization should not be nil.", func(ctx convey.C) {
				ctx.So(authorization, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUploadProxy(t *testing.T) {
	convey.Convey("UploadProxy", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			fileType = ""
			expire   = int64(1)
			body     = []byte("")
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			httpMock("PUT", d.c.Bfs.Addr+"/"+d.c.Bfs.Bucket).Reply(200).SetHeader("Code", "200")
			url, err := d.UploadProxy(c, fileType, expire, body)
			ctx.Convey("Then err should be nil.url should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(url, convey.ShouldNotBeNil)
			})
		})
	})
}
