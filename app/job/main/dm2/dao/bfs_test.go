package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoBfsData(t *testing.T) {
	convey.Convey("BfsData", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			bfsURL = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			testDao.BfsData(c, bfsURL)
		})
	})
}

func TestDaoBfsDmUpload(t *testing.T) {
	convey.Convey("BfsDmUpload", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			fileName = ""
			bs       = []byte("231231231231")
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			testDao.BfsDmUpload(c, fileName, bs)
		})
	})
}
