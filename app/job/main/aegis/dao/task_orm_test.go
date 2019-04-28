package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/job/main/aegis/model"
	"go-common/library/sync/errgroup"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoTaskActiveConfigs(t *testing.T) {
	convey.Convey("TaskActiveConfigs", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			configs, err := d.TaskActiveConfigs(c)
			ctx.Convey("Then err should be nil.configs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(configs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTaskActiveConsumer(t *testing.T) {
	convey.Convey("TaskActiveConsumer", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			consumerCache, err := d.TaskActiveConsumer(c)
			ctx.Convey("Then err should be nil.consumerCache should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(consumerCache, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoResource(t *testing.T) {
	convey.Convey("Resource", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			rid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.Resource(c, rid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoTaskRelease(t *testing.T) {
	convey.Convey("TaskRelease", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mtime = time.Now().Add(-10 * time.Second)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.TaskRelease(c, mtime)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoReport(t *testing.T) {
	convey.Convey("Report", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rt := &model.Report{
				BusinessID: 1,
				FlowID:     1,
				UID:        1,
				Content:    []byte("sguyiuo"),
			}
			err := d.Report(c, rt)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoTaskClear(t *testing.T) {
	convey.Convey("TaskClear", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.TaskClear(c, time.Now().Add(-3*24*time.Hour), 1000)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCheckFlow(t *testing.T) {
	convey.Convey("CheckFlow", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.CheckFlow(c, 1, 1)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCreateTask(t *testing.T) {
	f := func() error {
		for i := 1; i < 100; i++ {
			task := &model.Task{
				BusinessID: int64(i),
				RID:        int64(i),
				FlowID:     int64(i),
			}
			if err := d.CreateTask(context.Background(), task); err != nil {
				return err
			}
		}
		return nil
	}
	wg := errgroup.Group{}

	wg.Go(f)
	wg.Go(f)
	wg.Go(f)

	if err := wg.Wait(); err != nil {
		t.Fail()
	}
}
