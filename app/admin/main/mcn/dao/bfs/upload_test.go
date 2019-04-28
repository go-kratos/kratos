package bfs

import (
	"context"
	"io"
	"testing"

	"github.com/smartystreets/goconvey/convey"

	"gopkg.in/h2non/gock.v1"
)

func TestBfsUpload(t *testing.T) {
	convey.Convey("Upload", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			fileName = ""
			fileType = ""
			expire   = int64(0)
			body     io.Reader
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("PUT", d.bfs).Reply(200).SetHeader("code", "200").SetHeader("Location", "baidu")
			location, err := d.Upload(c, fileName, fileType, expire, body)
			ctx.Convey("Then err should be nil.location should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(location, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestBfsauthorize(t *testing.T) {
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
