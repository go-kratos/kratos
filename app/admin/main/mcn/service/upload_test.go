package service

import (
	"context"
	"io"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceUpload(t *testing.T) {
	convey.Convey("Upload", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			fileName = ""
			fileType = ""
			expire   = int64(0)
			body     io.Reader
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			location, err := s.Upload(c, fileName, fileType, expire, body)
			ctx.Convey("Then err should be nil.location should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(location, convey.ShouldNotBeNil)
			})
		})
	})
}
