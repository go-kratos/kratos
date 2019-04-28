package mysql

import (
	"context"
	"testing"
	"time"

	"go-common/app/admin/main/aegis/model/common"
	modtask "go-common/app/admin/main/aegis/model/task"

	"github.com/smartystreets/goconvey/convey"
)

func TestMysqlTaskFromDB(t *testing.T) {
	convey.Convey("TaskFromDB", t, func(ctx convey.C) {
		var (
			id = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			task, err := d.TaskFromDB(cntx, id)
			ctx.Convey("Then err should be nil.task should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(task, convey.ShouldBeNil)
			})
		})
	})
}

func TestMysqlDispatchByID(t *testing.T) {
	convey.Convey("DispatchByID", t, func(ctx convey.C) {
		var (
			mtasks map[int64]*modtask.Task
			ids    = []int64{0}
			args   = interface{}(int64(0))
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			missids, err := d.DispatchByID(cntx, mtasks, ids, args)
			ctx.Convey("Then err should be nil.missids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(missids, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMysqlDBDispatch(t *testing.T) {
	convey.Convey("DBDispatch", t, func(ctx convey.C) {
		var (
			opt = &modtask.NextOptions{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, count, err := d.DBDispatch(cntx, opt)
			ctx.Convey("Then err should be nil.tasks,count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMysqlRelease(t *testing.T) {
	convey.Convey("Release", t, func(ctx convey.C) {
		var (
			opt = &common.BaseOptions{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.Release(cntx, opt, true)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMysqlSeize(t *testing.T) {
	convey.Convey("Seize", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mapids = map[int64]int64{1: 1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.Seize(c, mapids)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMysqlDelay(t *testing.T) {
	convey.Convey("Delay", t, func(ctx convey.C) {
		var (
			opt = &modtask.DelayOptions{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.Delay(cntx, opt)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMysqlListCheckUnSeized(t *testing.T) {
	convey.Convey("ListCheckUnSeized", t, func(ctx convey.C) {
		var (
			mtasks = map[int64]*modtask.Task{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.ListCheckUnSeized(cntx, mtasks, []int64{})
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMysqlListCheckSeized(t *testing.T) {
	convey.Convey("ListCheckSeized", t, func(ctx convey.C) {
		var (
			mtasks = map[int64]*modtask.Task{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.ListCheckSeized(cntx, mtasks, []int64{}, int64(1))
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMysqlListCheckDelay(t *testing.T) {
	convey.Convey("ListCheckDelay", t, func(ctx convey.C) {
		var (
			mtasks = map[int64]*modtask.Task{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.ListCheckDelay(cntx, mtasks, []int64{}, int64(1))
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMysqlListTasks(t *testing.T) {
	convey.Convey("ListTasks", t, func(ctx convey.C) {
		opt := &modtask.ListOptions{
			State: 4,
		}
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, _, err := d.ListTasks(context.TODO(), opt)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMysqllistCheck(t *testing.T) {
	convey.Convey("listCheck", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.listCheck(context.TODO(), "state=1", []int64{1})
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMysqlConsumerOn(t *testing.T) {
	convey.Convey("ConsumerOn", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			opt = &common.BaseOptions{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.ConsumerOn(c, opt)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMysqlConsumerOff(t *testing.T) {
	convey.Convey("ConsumerOff", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			opt = &common.BaseOptions{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.ConsumerOff(c, opt)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMysqlIsConsumerOn(t *testing.T) {
	convey.Convey("IsConsumerOn", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			opt = &common.BaseOptions{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			on, err := d.IsConsumerOn(c, opt)
			ctx.Convey("Then err should be nil.on should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(on, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMysqlconsumer(t *testing.T) {
	convey.Convey("consumer", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			opt    = &common.BaseOptions{}
			action = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.consumer(c, opt, action)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMysqlConsumerStat(t *testing.T) {
	convey.Convey("ConsumerStat", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			bizid  = int64(0)
			flowid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			items, err := d.ConsumerStat(c, bizid, flowid)
			ctx.Convey("Then err should be nil.items should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(items, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMysqlOnlines(t *testing.T) {
	convey.Convey("Onlines", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			opt = &common.BaseOptions{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			uids, err := d.Onlines(c, opt)
			ctx.Convey("Then err should be nil.uids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(uids, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMysqlQueryTask(t *testing.T) {
	convey.Convey("QueryTask", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, _, err := d.QueryTask(context.TODO(), 0, time.Now(), 0, 0)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMysqlCountPersonal(t *testing.T) {
	convey.Convey("CountPersonal", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			opt = &common.BaseOptions{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.CountPersonal(c, opt)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMysqlQueryForSeize(t *testing.T) {
	convey.Convey("QueryForSeize", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.QueryForSeize(context.TODO(), 0, 0, 0, 0)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
