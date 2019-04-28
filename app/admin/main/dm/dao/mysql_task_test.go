package dao

import (
	"context"
	"testing"

	"go-common/app/admin/main/dm/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoTaskList(t *testing.T) {
	convey.Convey("TaskList", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			taskSQL = []string{}
			pn      = int64(1)
			ps      = int64(50)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tasks, total, err := testDao.TaskList(c, taskSQL, pn, ps)
			ctx.Convey("Then err should be nil.tasks,total should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
				t.Logf("====%d", total)
				ctx.So(tasks, convey.ShouldNotBeNil)
				// t.Logf("====%+v", tasks[0])
			})
		})
	})
}

func TestDaoAddTask(t *testing.T) {
	convey.Convey("AddTask", t, func(ctx convey.C) {
		var (
			tx, _ = testDao.BeginBiliDMTrans(context.Background())
			v     = &model.AddTaskArg{
				Title: "test",
				Start: "2016-10-30 16:12:21",
				End:   "2018-10-30 16:12:21",
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			taskID, err := testDao.AddTask(tx, v, 1)
			ctx.Convey("Then err should be nil.taskID should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(taskID, convey.ShouldNotBeNil)
				t.Logf("====%d", taskID)
				tx.Commit()
			})
		})
	})
}

func TestDaoAddSubTask(t *testing.T) {
	convey.Convey("AddSubTask", t, func(ctx convey.C) {
		var (
			tx, _     = testDao.BeginBiliDMTrans(context.Background())
			taskID    = int64(8)
			operation = int32(0)
			start     = "2018-10-30 16:12:21"
			rate      = int32(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			id, err := testDao.AddSubTask(tx, taskID, operation, start, rate)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
				t.Logf("====%d", id)
			})
			tx.Commit()
		})
	})
}

func TestDaoEditTaskState(t *testing.T) {
	convey.Convey("EditTaskState", t, func(ctx convey.C) {
		var (
			c = context.Background()
			v = &model.EditTasksStateArg{
				IDs:   "1,7",
				State: 1,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			affected, err := testDao.EditTaskState(c, v)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTaskView(t *testing.T) {
	convey.Convey("TaskView", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			task, err := testDao.TaskView(c, id)
			ctx.Convey("Then err should be nil.task should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(task, convey.ShouldNotBeNil)
				t.Logf("====%+v", task)
			})
		})
	})
}

func TestDaoSubTask(t *testing.T) {
	convey.Convey("SubTask", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(8)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			subTask, err := testDao.SubTask(c, id)
			ctx.Convey("Then err should be nil.subTask should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(subTask, convey.ShouldNotBeNil)
				t.Logf("====%+v", subTask)
			})
		})
	})
}

func TestDaoReviewTask(t *testing.T) {
	convey.Convey("ReviewTask", t, func(ctx convey.C) {
		var (
			c = context.Background()
			v = &model.ReviewTaskArg{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			affected, err := testDao.ReviewTask(c, v)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}
func TestDaoEditTaskPriority(t *testing.T) {
	convey.Convey("EditTaskPriority", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			ids      = "12"
			priority = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			affected, err := testDao.EditTaskPriority(c, ids, priority)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}
