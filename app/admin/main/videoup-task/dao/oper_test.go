package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddVideoOper(t *testing.T) {
	var (
		c         = context.TODO()
		aid       = int64(1)
		adminID   = int64(421)
		vid       = int64(1)
		attribute = int32(0)
		status    = int16(0)
		lastID    = int64(0)
		content   = "测试"
		remark    = "测试"
	)
	convey.Convey("AddVideoOper", t, func(ctx convey.C) {
		id, err := d.AddVideoOper(c, aid, adminID, vid, attribute, status, lastID, content, remark)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(id, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoVideoOpers(t *testing.T) {
	var (
		c   = context.TODO()
		vid = int64(1)
	)
	convey.Convey("VideoOpers", t, func(ctx convey.C) {
		op, uids, err := d.VideoOpers(c, vid)
		ctx.Convey("Then err should be nil.op,uids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(uids, convey.ShouldNotBeNil)
			ctx.So(op, convey.ShouldNotBeNil)
		})
	})
}
