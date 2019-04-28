package service

import (
	"context"
	"testing"

	"go-common/app/job/main/dm2/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestServicetaskResProc(t *testing.T) {
	convey.Convey("taskResProc", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			svr.taskResProc()
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestServicetaskDelProc(t *testing.T) {
	convey.Convey("taskDelProc", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			svr.taskDelProc()
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestServicetaskDelDM(t *testing.T) {
	convey.Convey("taskDelDM", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			task = &model.TaskInfo{
				Result: "http://berserker.bilibili.co/avenger/download/hdfs?path=/api/hive/query/148/672bc22888af701529e8b3052fd2c4a7/1543546463/1547966/result",
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			delCount, lastIndex, state, err := svr.taskDelDM(c, task)
			ctx.Convey("Then err should be nil.delCount,pause should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(state, convey.ShouldNotBeNil)
				ctx.So(delCount, convey.ShouldNotBeNil)
				t.Log(delCount, lastIndex, state)
			})
		})
	})
}
