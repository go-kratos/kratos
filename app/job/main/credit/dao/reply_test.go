package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoRegReply(t *testing.T) {
	convey.Convey("RegReply", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(1)
			tid = int8(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.RegReply(c, id, tid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelReply(t *testing.T) {
	convey.Convey("DelReply", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			rpid = ""
			tp   = ""
			oid  = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelReply(c, rpid, tp, oid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpReplyState(t *testing.T) {
	convey.Convey("UpReplyState", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			oid   = int64(0)
			rpid  = int64(0)
			tp    = int8(0)
			state = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.UpReplyState(c, oid, rpid, tp, state)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
