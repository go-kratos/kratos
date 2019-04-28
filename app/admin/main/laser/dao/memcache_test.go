package dao

import (
	"context"
	"go-common/app/admin/main/laser/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyTaskInfo(t *testing.T) {
	convey.Convey("keyTaskInfo", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyTaskInfo(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTaskInfoCache(t *testing.T) {
	convey.Convey("TaskInfoCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.TaskInfoCache(c, mid)
			ctx.Convey("Then err should be nil.ti should not be nil.", func(ctx convey.C) {
			})
		})
	})
}

func TestDaoAddTaskInfoCache(t *testing.T) {
	convey.Convey("AddTaskInfoCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			ti  = &model.TaskInfo{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.AddTaskInfoCache(c, mid, ti)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {

			})
		})
	})
}

func TestDaoRemoveTaskInfoCache(t *testing.T) {
	convey.Convey("RemoveTaskInfoCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.RemoveTaskInfoCache(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {

			})
		})
	})
}
