package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAssignTask(t *testing.T) {
	convey.Convey("AssignTask", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			task = task1
		)
		task.AdminID = 1
		task.UID = 1
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.AssignTask(c, task)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoQueryTask(t *testing.T) {
	convey.Convey("QueryTask", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			state = int8(0)
			mtime = time.Now()
			id    = int64(0)
			limit = int64(1000)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tasks, lastid, err := d.QueryTask(c, state, mtime, id, limit)
			ctx.Convey("Then err should be nil.tasks,lastid should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(lastid, convey.ShouldNotBeNil)
				ctx.So(tasks, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSetWeightDB(t *testing.T) {
	convey.Convey("SetWeightDB", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			taskid = int64(1)
			weight = int64(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.SetWeightDB(c, taskid, weight)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
