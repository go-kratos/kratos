package redis

import (
	"context"
	"testing"

	"go-common/app/admin/main/aegis/model/common"
	modtask "go-common/app/admin/main/aegis/model/task"

	"github.com/smartystreets/goconvey/convey"
)

var (
	task1 = &modtask.Task{
		ID:         1,
		BusinessID: 1,
		FlowID:     1,
		RID:        1,
		Gtime:      0,
		Weight:     10,
	}
	task2 = &modtask.Task{
		ID:         2,
		BusinessID: 1,
		FlowID:     1,
		RID:        2,
		Gtime:      0,
		Weight:     8,
	}
	task3 = &modtask.Task{
		ID:         3,
		BusinessID: 1,
		FlowID:     1,
		RID:        3,
		Gtime:      0,
		Weight:     8,
	}
)

func TestRedisSetTask(t *testing.T) {
	convey.Convey("SetTask", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetTask(c, task1, task2, task3)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestRedisPushPublicTask(t *testing.T) {
	convey.Convey("PushPublicTask", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.PushPublicTask(c, task1, task2, task3)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestRedisRemovePublicTask(t *testing.T) {
	convey.Convey("RemovePublicTask", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			opt = &common.BaseOptions{
				BusinessID: 1,
				FlowID:     1,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.RemovePublicTask(c, opt, 1, 2, 3)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestRedisPushPersonalTask(t *testing.T) {
	convey.Convey("PushPersonalTask", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			opt = &common.BaseOptions{
				BusinessID: 1,
				FlowID:     1,
				UID:        1,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.PushPersonalTask(c, opt, 1, 2, 3)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestRedisCountPersonalTask(t *testing.T) {
	convey.Convey("CountPersonalTask", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			opt = &common.BaseOptions{
				BusinessID: 1,
				FlowID:     1,
				UID:        1,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.CountPersonalTask(c, opt)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestRedisRangePersonalTask(t *testing.T) {
	convey.Convey("RangePersonalTask", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			opt = &modtask.ListOptions{
				BaseOptions: common.BaseOptions{
					BusinessID: 1,
					FlowID:     1,
					UID:        1},
				Pager: common.Pager{
					Pn: 1,
					Ps: 20,
				},
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, _, _, _, err := d.RangePersonalTask(c, opt)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestRedisRemovePersonalTask(t *testing.T) {
	convey.Convey("RemovePersonalTask", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			opt = &common.BaseOptions{
				BusinessID: 1,
				FlowID:     1,
				UID:        1,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.RemovePersonalTask(c, opt, 1, 2, 3)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestRedisPushDelayTask(t *testing.T) {
	convey.Convey("PushDelayTask", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			opt = &common.BaseOptions{
				BusinessID: 1,
				FlowID:     1,
				UID:        1,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.PushDelayTask(c, opt, 1, 2, 3)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestRedisRangeDealyTask(t *testing.T) {
	convey.Convey("RangeDealyTask", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			opt = &modtask.ListOptions{
				BaseOptions: common.BaseOptions{
					BusinessID: 1,
					FlowID:     1,
					UID:        1},
				Pager: common.Pager{
					Pn: 1,
					Ps: 20,
				},
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tasks, count, hitids, _, err := d.RangeDealyTask(c, opt)
			ctx.Convey("Then err should be nil.tasks,count,hitids,missids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(hitids, convey.ShouldNotBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
				ctx.So(tasks, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestRedisRelease(t *testing.T) {
	convey.Convey("Release", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			opt = &common.BaseOptions{
				BusinessID: 1,
				FlowID:     1,
				UID:        1,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Release(c, opt, true)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestRedisrangefunc(t *testing.T) {
	convey.Convey("rangefunc", t, func(ctx convey.C) {

		var (
			c   = context.Background()
			opt = &modtask.ListOptions{
				BaseOptions: common.BaseOptions{
					BusinessID: 1,
					FlowID:     1,
					UID:        1},
				Pager: common.Pager{
					Pn: 1,
					Ps: 20,
				},
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, _, _, _, err := d.rangefuncCluster(c, "public", opt)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestRedisSeize(t *testing.T) {
	convey.Convey("SeizeTask", t, func(ctx convey.C) {

		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, _, _, err := d.SeizeTask(c, 1, 1, 1, 10)
			ctx.Convey("Then err should be nil.tasks,count,hitids,missids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestRedisRemoveDelayTask(t *testing.T) {
	convey.Convey("RemoveDelayTask", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			opt = &common.BaseOptions{}
			ids = interface{}(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.RemoveDelayTask(c, opt, ids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestRedispushList(t *testing.T) {
	convey.Convey("pushList", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			key    = ""
			values = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.pushList(c, key, values)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestRedisremoveList(t *testing.T) {
	convey.Convey("removeList", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			key = ""
			ids = interface{}(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.removeList(c, key, ids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestRedisGetTask(t *testing.T) {
	convey.Convey("GetTask", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.GetTask(cntx, []int64{1, 2, 3})
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
