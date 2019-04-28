package service

import (
	"context"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestServiceAddTask(t *testing.T) {
	convey.Convey("AddTask", t, func(ctx convey.C) {
		var (
			c            = context.Background()
			mid          = int64(0)
			username     = ""
			adminID      = int64(0)
			logDate      = int64(0)
			contactEmail = ""
			platform     = int(0)
			sourceType   = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			s.AddTask(c, mid, username, adminID, logDate, contactEmail, platform, sourceType)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {

			})
		})
	})
}

func TestServiceDeleteTask(t *testing.T) {
	convey.Convey("DeleteTask", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			taskID   = int64(0)
			username = ""
			adminID  = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			s.DeleteTask(c, taskID, username, adminID)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {

			})
		})
	})
}

func TestServiceQueryTask(t *testing.T) {
	convey.Convey("QueryTask", t, func(ctx convey.C) {
		var (
			c            = context.Background()
			mid          = int64(0)
			logDateStart = int64(0)
			logDateEnd   = int64(0)
			sourceType   = int(0)
			platform     = int(0)
			state        = int(0)
			sortBy       = ""
			pageNo       = int(0)
			pageSize     = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			s.QueryTask(c, mid, logDateStart, logDateEnd, sourceType, platform, state, sortBy, pageNo, pageSize)
			ctx.Convey("Then err should be nil.tasks,count should not be nil.", func(ctx convey.C) {
			})
		})
	})
}

func TestServicebuildSortStmt(t *testing.T) {
	convey.Convey("buildSortStmt", t, func(ctx convey.C) {
		var (
			sortBy = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			sort := buildSortStmt(sortBy)
			ctx.Convey("Then sort should not be nil.", func(ctx convey.C) {
				ctx.Convey(sort, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceUpdateTask(t *testing.T) {
	convey.Convey("UpdateTask", t, func(ctx convey.C) {
		var (
			c            = context.Background()
			username     = ""
			adminID      = int64(0)
			taskID       = int64(0)
			mid          = int64(0)
			logDate      = int64(0)
			contactEmail = ""
			sourceType   = int(0)
			platform     = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			s.UpdateTask(c, username, adminID, taskID, mid, logDate, contactEmail, sourceType, platform)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {

			})
		})
	})
}
