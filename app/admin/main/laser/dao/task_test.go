package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddTask(t *testing.T) {
	convey.Convey("AddTask", t, func(ctx convey.C) {
		var (
			c            = context.Background()
			mid          = int64(2233)
			username     = "yanjinbin"
			adminID      = int64(479)
			logDate      = "2018-11-28 20:19:09"
			contactEmail = "yanjinbin@qq.com"
			platform     = int(1)
			sourceType   = int(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.AddTask(c, mid, username, adminID, logDate, contactEmail, platform, sourceType)
			ctx.Convey("Then err should be nil.lastInsertID should not be nil.", func(ctx convey.C) {

			})
		})
	})
}

func TestDaoFindTask(t *testing.T) {
	convey.Convey("FindTask", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(0)
			state = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.FindTask(c, mid, state)
			ctx.Convey("Then err should be nil.no should not be nil.", func(ctx convey.C) {

			})
		})
	})
}

func TestDaoQueryTaskInfoByIDSQL(t *testing.T) {
	convey.Convey("QueryTaskInfoByIDSQL", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.QueryTaskInfoByIDSQL(c, id)
			ctx.Convey("Then err should be nil.no should not be nil.", func(ctx convey.C) {

			})
		})
	})
}

func TestDaoDeleteTask(t *testing.T) {
	convey.Convey("DeleteTask", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			taskID   = int64(0)
			username = ""
			adminID  = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DeleteTask(c, taskID, username, adminID)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpdateTask(t *testing.T) {
	convey.Convey("UpdateTask", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			taskID     = int64(0)
			state      = int(0)
			updateStmt = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.UpdateTask(c, taskID, state, updateStmt)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {

			})
		})
	})
}

func TestDaoQueryTask(t *testing.T) {
	convey.Convey("QueryTask", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			queryStmt = ""
			sort      = ""
			offset    = int(0)
			limit     = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.QueryTask(c, queryStmt, sort, offset, limit)
			ctx.Convey("Then err should be nil.tasks,count should not be nil.", func(ctx convey.C) {

			})
		})
	})
}
