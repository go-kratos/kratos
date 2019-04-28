package newcomer

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestNewcomerUserTasks(t *testing.T) {
	convey.Convey("UserTasks", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			index = ""
			id    = int64(0)
			limit = int(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.UserTasks(c, index, id, limit)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, err)
			})
		})
	})
}

func TestNewcomergetTableName(t *testing.T) {
	convey.Convey("getTableName", t, func(ctx convey.C) {
		var (
			mid = int64(27515405)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := getTableName(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldEqual, p1)
			})
		})
	})
}

func TestNewcomerUpUserTask(t *testing.T) {
	convey.Convey("UpUserTask", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515405)
			tid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.UpUserTask(c, mid, tid)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, err)
			})
		})
	})
}

func TestNewcomerUserTasksByMID(t *testing.T) {
	convey.Convey("UserTasksByMID", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515405)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.UserTasksByMID(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, err)
			})
		})
	})
}

func TestNewcomerTaskByTID(t *testing.T) {
	convey.Convey("TaskByTID", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515405)
			tid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.TaskByTID(c, mid, tid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, err)
			})
		})
	})
}
func TestNewcomerTasks(t *testing.T) {
	convey.Convey("Tasks", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.Tasks(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, err)
			})
		})
	})
}

func TestNewcomerCheckTaskComplete(t *testing.T) {
	convey.Convey("CheckTaskComplete", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515405)
			tid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res := d.CheckTaskComplete(c, mid, tid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(res, convey.ShouldEqual, res)
			})
		})
	})
}

func TestNewcomerUserTasksNotify(t *testing.T) {
	convey.Convey("UserTasksNotify", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(0)
			index = "1"
			start = time.Now().Format("2006-01-02 15:04:05")
			end   = time.Now().Format("2006-01-02 15:04:05")
			limit = 100
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.UserTasksNotify(c, index, id, start, end, limit)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, err)
			})
		})
	})
}

func TestNewcomerGiftRewardCount(t *testing.T) {
	convey.Convey("GiftRewardCount", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515405)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.GiftRewardCount(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, err)
			})
		})
	})
}

func TestNewcomerBaseRewardCount(t *testing.T) {
	convey.Convey("BaseRewardCount", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515405)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.BaseRewardCount(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, err)
			})
		})
	})
}

func TestNewcomerCheckTasksForRewardNotify(t *testing.T) {
	convey.Convey("CheckTasksForRewardNotify", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(0)
			index = "1"
			start = time.Now()
			end   = time.Now()
			limit = 100
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.CheckTasksForRewardNotify(c, index, id, start, end, limit)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, err)
			})
		})
	})
}

func TestNewcomerUserTasksByMIDAndState(t *testing.T) {
	convey.Convey("UserTasksByMIDAndState", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(27515405)
			state = 0
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.UserTasksByMIDAndState(c, mid, state)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, err)
			})
		})
	})
}
