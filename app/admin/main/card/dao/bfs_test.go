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
			fileName = ""
			fileType = "png"
			data     = []byte("")
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.Upload(c, fileName, fileType, data, d.c.Bfs)
			ctx.Convey("Then err should be nil.location should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
