package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoUploadBfs(t *testing.T) {
	convey.Convey("UploadBfs", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			fileName = ""
			bs       = []byte("123")
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.UploadBfs(c, fileName, bs)
		})
	})
}
