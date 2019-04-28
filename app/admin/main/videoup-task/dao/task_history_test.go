package dao

import (
	"context"
	"go-common/app/admin/main/videoup-task/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddTaskHis(t *testing.T) {
	convey.Convey("AddTaskHis", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			pool   = int8(0)
			action = int8(0)
			taskID = int64(0)
			cid    = int64(0)
			uid    = int64(0)
			utime  = int64(0)
			result = int16(0)
			reason = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.AddTaskHis(c, pool, action, taskID, cid, uid, utime, result, reason)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoMulAddTaskHis(t *testing.T) {
	convey.Convey("MulAddTaskHis", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			tls    = []*model.TaskForLog{}
			action = int8(0)
			uid    = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.MulAddTaskHis(c, tls, action, uid)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
