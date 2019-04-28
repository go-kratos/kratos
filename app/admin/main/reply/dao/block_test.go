package dao

import (
	"context"
	"go-common/app/admin/main/reply/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoBlockAccount(t *testing.T) {
	convey.Convey("BlockAccount", t, func(ctx convey.C) {
		var (
			c             = context.Background()
			mid           = int64(0)
			ftime         = int64(0)
			notify        bool
			freason       = int32(0)
			originTitle   = ""
			originContent = ""
			redirectURL   = ""
			adname        = ""
			remark        = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := _d.BlockAccount(c, mid, ftime, notify, freason, originTitle, originContent, redirectURL, adname, remark)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTransferArbitration(t *testing.T) {
	convey.Convey("TransferArbitration", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			rps    map[int64]*model.Reply
			rpts   map[int64]*model.Report
			adid   = int64(0)
			adname = ""
			titles map[int64]string
			links  map[int64]string
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := _d.TransferArbitration(c, rps, rpts, adid, adname, titles, links)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
