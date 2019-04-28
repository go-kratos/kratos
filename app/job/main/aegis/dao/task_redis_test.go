package dao

import (
	"context"
	"fmt"
	"testing"

	"go-common/app/job/main/aegis/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSetTask(t *testing.T) {
	convey.Convey("SetTask", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err1 := d.SetTask(c, task1)
			err2 := d.SetTask(c, task2)
			err3 := d.SetTask(c, task3)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err1, convey.ShouldBeNil)
				ctx.So(err2, convey.ShouldBeNil)
				ctx.So(err3, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoGetTask(t *testing.T) {
	convey.Convey("GetTask", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			task, err := d.GetTask(c, id)
			ctx.Convey("Then err should be nil.task should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(task, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoPushPublicTask(t *testing.T) {
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

func TestDaoRemovePublicTask(t *testing.T) {
	convey.Convey("RemovePublicTask", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.RemovePublicTask(c, 1, 1, 1)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoPushPersonalTask(t *testing.T) {
	convey.Convey("PushPersonalTask", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			businessID = int64(1)
			flowID     = int64(1)
			uid        = int64(1)
			taskid     = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.PushPersonalTask(c, businessID, flowID, uid, taskid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoRemovePersonalTask(t *testing.T) {
	convey.Convey("RemovePersonalTask", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			businessID = int64(1)
			flowID     = int64(1)
			uid        = int64(1)
			taskid     = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.RemovePersonalTask(c, businessID, flowID, uid, taskid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoPushDelayTask(t *testing.T) {
	convey.Convey("PushDelayTask", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			businessID = int64(1)
			flowID     = int64(1)
			uid        = int64(1)
			taskid     = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.PushDelayTask(c, businessID, flowID, uid, taskid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoRemoveDelayTask(t *testing.T) {
	convey.Convey("RemoveDelayTask", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			businessID = int64(1)
			flowID     = int64(1)
			uid        = int64(1)
			taskid     = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.RemoveDelayTask(c, businessID, flowID, uid, taskid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetWeight(t *testing.T) {
	convey.Convey("SetWeight", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			businessID = int64(1)
			flowID     = int64(1)
			id         = int64(1)
			weight     = int64(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.SetWeight(c, businessID, flowID, id, weight)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			})
		})
	})
}

func TestDaoGetWeight(t *testing.T) {
	convey.Convey("GetWeight", t, func(ctx convey.C) {
	})
}

func TestDaoTopWeights(t *testing.T) {
	convey.Convey("TopWeights", t, func(ctx convey.C) {
	})
}

func TestDaoCreateUnionSet(t *testing.T) {
	convey.Convey("CreateUnionSet", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			businessID = int64(1)
			flowID     = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.CreateUnionSet(c, businessID, flowID)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoRangeUinonSet(t *testing.T) {
	convey.Convey("RangeUinonSet", t, func(ctx convey.C) {
	})
}

func TestDaoDeleteUinonSet(t *testing.T) {
	convey.Convey("DeleteUinonSet", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			businessID = int64(1)
			flowID     = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DeleteUinonSet(c, businessID, flowID)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoIncresByField(t *testing.T) {
	convey.Convey("IncresByField", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.IncresByField(c, 1, 1, 1, model.Dispatch, 1)

			err = d.IncresByField(c, 1, 1, 1, model.Release, 1)
			err = d.IncresByField(c, 1, 1, 1, model.Submit, 1)
			err = d.IncresByField(c, 1, 1, 1, model.Delay, 1)
			err = d.IncresByField(c, 1, 1, 1, fmt.Sprintf(model.RscState, 1), 1)
			err = d.IncresByField(c, 1, 1, 1, model.UseTime, 112)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoFlushReport(t *testing.T) {
	convey.Convey("FlushReport", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			data, err := d.FlushReport(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				for key, val := range data {
					fmt.Println("key:", key)
					fmt.Println("val:", string(val))
				}
			})
		})
	})
}
