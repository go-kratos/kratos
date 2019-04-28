package dao

import (
	"context"
	"testing"

	pushmdl "go-common/app/service/main/push/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoBusinesses(t *testing.T) {
	convey.Convey("Businesses", t, func(ctx convey.C) {
		res, err := d.Businesses(context.Background())
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func addTask() (id int64, err error) {
	t := &pushmdl.Task{APPID: 1}
	return d.AddTask(context.Background(), t)
}

func TestDaoAddTask(t *testing.T) {
	convey.Convey("AddTask", t, func(ctx convey.C) {
		id, err := addTask()
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(id, convey.ShouldBeGreaterThan, 0)
		})
	})
}

func TestDaoTask(t *testing.T) {
	id, _ := addTask()
	convey.Convey("Task", t, func(ctx convey.C) {
		no, err := d.Task(context.Background(), id)
		ctx.Convey("Then err should be nil.no should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(no, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoMaxSettingID(t *testing.T) {
	convey.Convey("MaxSettingID", t, func(ctx convey.C) {
		_, err := d.MaxSettingID(context.Background())
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSettingsByRange(t *testing.T) {
	var (
		start = int64(0)
		end   = int64(1000)
	)
	convey.Convey("SettingsByRange", t, func(ctx convey.C) {
		res, err := d.SettingsByRange(context.Background(), start, end)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
