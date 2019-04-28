package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/job/main/dm2/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoTaskInfos(t *testing.T) {
	convey.Convey("TaskInfos", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			state = int32(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tasks, err := testDao.TaskInfos(c, state)
			ctx.Convey("Then err should be nil.tasks should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(tasks, convey.ShouldNotBeNil)
				for _, task := range tasks {
					t.Logf("%+v", task)
				}
			})
		})
	})
}

func TestDaoUpdateTask(t *testing.T) {
	convey.Convey("UpdateTask", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			task = &model.TaskInfo{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			affected, err := testDao.UpdateTask(c, task)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelDMs(t *testing.T) {
	convey.Convey("DelDMs", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			oid   = int64(221)
			dmids = []int64{719182141}
			state = int32(12)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			affected, err := testDao.DelDMs(c, oid, dmids, state)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUptSubTask(t *testing.T) {
	convey.Convey("UptSubTask", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			taskID   = int64(0)
			delCount = int64(0)
			end      = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			affected, err := testDao.UptSubTask(c, taskID, delCount, end)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTaskSearchRes(t *testing.T) {
	convey.Convey("TaskSearchRes", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			task = &model.TaskInfo{
				Topic: "http://berserker.bilibili.co/m-avenger/api/hive/status/query/148/672bc22888af701529e8b3052fd2c4a7/1541066053/1389520",
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, result, state, err := testDao.TaskSearchRes(c, task)
			ctx.Convey("Then err should be nil.result,state should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(state, convey.ShouldNotBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
				t.Log(result, state)
			})
		})
	})
}

func TestDaoUptSubjectCount(t *testing.T) {
	convey.Convey("UptSubjectCount", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tp    = int32(1)
			oid   = int64(0)
			count = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			affected, err := testDao.UptSubjectCount(c, tp, oid, count)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSendWechatWorkMsg(t *testing.T) {
	convey.Convey("SendWechatWorkMsg", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			content = "test"
			title   = "test"
			users   = []string{"fengduzhen"}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := testDao.SendWechatWorkMsg(c, content, title, users)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSubTask(t *testing.T) {
	convey.Convey("SubTask", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(32)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := testDao.SubTask(c, id)
			ctx.Convey("Then err should be nil.subTask should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
