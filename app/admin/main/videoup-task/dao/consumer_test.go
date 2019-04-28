package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoTaskUserCheckIn(t *testing.T) {
	convey.Convey("TaskUserCheckIn", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			uid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.TaskUserCheckIn(c, uid)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTaskUserCheckOff(t *testing.T) {
	convey.Convey("TaskUserCheckOff", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			uid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.TaskUserCheckOff(c, uid)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoConsumers(t *testing.T) {
	convey.Convey("Consumers", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			cms, err := d.Consumers(c)
			ctx.Convey("Then err should be nil.cms should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(cms, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoIsConsumerOn(t *testing.T) {
	convey.Convey("IsConsumerOn", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			uid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			state := d.IsConsumerOn(c, uid)
			ctx.Convey("Then state should not be nil.", func(ctx convey.C) {
				ctx.So(state, convey.ShouldNotBeNil)
			})
		})
	})
}
