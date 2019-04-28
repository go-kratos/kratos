package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/dm/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceTaskList(t *testing.T) {
	convey.Convey("TaskList", t, func(ctx convey.C) {
		var (
			c = context.Background()
			v = &model.TaskListArg{
				Pn: 1,
				Ps: 10,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := svr.TaskList(c, v)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
				for _, v := range res.Result {
					t.Logf("=====%+v", v)
				}

			})
		})
	})
}

func TestServiceAddTask(t *testing.T) {
	convey.Convey("AddTask", t, func(ctx convey.C) {
		var (
			c = context.Background()
			v = &model.AddTaskArg{
				IPs:    "172.1.1.1",
				Start:  "2016-10-30 16:12:21",
				End:    "2018-10-30 16:12:21",
				OpTime: "2018-10-30 16:12:21",
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := svr.AddTask(c, v)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestServiceEditTaskState(t *testing.T) {
	convey.Convey("EditTaskState", t, func(ctx convey.C) {
		var (
			c = context.Background()
			v = &model.EditTasksStateArg{
				IDs:   "7",
				State: 3,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := svr.EditTaskState(c, v)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestServiceTaskView(t *testing.T) {
	convey.Convey("TaskView", t, func(ctx convey.C) {
		var (
			c = context.Background()
			v = &model.TaskViewArg{
				ID: int64(7),
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			task, err := svr.TaskView(c, v)
			ctx.Convey("Then err should be nil.task should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(task, convey.ShouldNotBeNil)
				t.Logf("====%+v\n", task)
				t.Logf("====%+v", task.SubTask)
			})
		})
	})
}

func TestServiceReviewTask(t *testing.T) {
	convey.Convey("ReviewTask", t, func(ctx convey.C) {
		var (
			c = context.Background()
			v = &model.ReviewTaskArg{
				ID: 21,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := svr.ReviewTask(c, v)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
