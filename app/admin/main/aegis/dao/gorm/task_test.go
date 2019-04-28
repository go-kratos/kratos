package gorm

import (
	"context"
	"testing"

	"go-common/app/admin/main/aegis/model/common"
	taskmod "go-common/app/admin/main/aegis/model/task"

	"github.com/smartystreets/goconvey/convey"
)

func TestUndoStat(t *testing.T) {
	convey.Convey("UndoStat", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.UndoStat(c, 0, 0, 0)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestTaskStat(t *testing.T) {
	convey.Convey("TaskStat", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.TaskStat(c, 0, 0, 0)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestTaskMaxWeight(t *testing.T) {
	convey.Convey("MaxWeight", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.MaxWeight(c, 0, 0)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestTaskListSeized(t *testing.T) {
	convey.Convey("TaskListSeized", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, _, err := d.TaskListSeized(c, &taskmod.ListOptions{})
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestTxCloseTasks(t *testing.T) {
	convey.Convey("TxCloseTasks", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTx(c)
		)
		defer tx.Commit()
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.TxCloseTasks(tx, []int64{}, 0)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestTxSubmit(t *testing.T) {
	convey.Convey("TxSubmit", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTx(cntx)
			opt   = &taskmod.SubmitOptions{}
		)
		defer tx.Commit()
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.TxSubmit(tx, opt, 0)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestCloseTask(t *testing.T) {
	convey.Convey("CloseTask", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.CloseTask(cntx, 0)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestTaskByRID(t *testing.T) {
	convey.Convey("TaskByRID", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.TaskByRID(cntx, 0, 1)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestTaskListAssignd(t *testing.T) {
	convey.Convey("TaskListAssignd", t, func(ctx convey.C) {
		var (
			opt = &taskmod.ListOptions{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, _, err := d.TaskListAssignd(cntx, opt)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestTaskListDelayd(t *testing.T) {
	convey.Convey("TaskListDelayd", t, func(ctx convey.C) {
		var (
			opt = &taskmod.ListOptions{
				BaseOptions: common.BaseOptions{
					UID: 1,
				},
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, _, err := d.TaskListDelayd(cntx, opt)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
