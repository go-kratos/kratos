package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoRankElecUpdateOrder(t *testing.T) {
	convey.Convey("RankElecUpdateOrder", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			avID   = int64(0)
			upMID  = int64(0)
			payMID = int64(0)
			ver    = int64(0)
			fee    = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.RankElecUpdateOrder(c, avID, upMID, payMID, ver, fee)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoRankElecUpdateMessage(t *testing.T) {
	convey.Convey("RankElecUpdateMessage", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			avID    = int64(0)
			upMID   = int64(0)
			payMID  = int64(0)
			ver     = int64(0)
			message = ""
			hidden  bool
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.RankElecUpdateMessage(c, avID, upMID, payMID, ver, message, hidden)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
