package service

import (
	"context"
	"go-common/app/admin/main/up/conf"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceUpload(t *testing.T) {
	convey.Convey("Upload", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			fileName = ""
			fileType = ""
			no       = time.Now()
			body     = []byte("")
			bfs      = &conf.Bfs{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			location, err := s.Upload(c, fileName, fileType, no, body, bfs)
			ctx.Convey("Then err should be nil.location should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(location, convey.ShouldNotBeNil)
			})
		})
	})
}
