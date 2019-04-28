package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSetnxTaskJob(t *testing.T) {
	convey.Convey("SetnxTaskJob", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			value = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ok, err := testDao.SetnxTaskJob(c, value)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetTaskJob(t *testing.T) {
	convey.Convey("GetTaskJob", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			testDao.GetTaskJob(c)
		})
	})
}

func TestDaoGetSetTaskJob(t *testing.T) {
	convey.Convey("GetSetTaskJob", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			value = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			old, err := testDao.GetSetTaskJob(c, value)
			ctx.Convey("Then err should be nil.old should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(old, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelTaskJob(t *testing.T) {
	convey.Convey("DelTaskJob", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := testDao.DelTaskJob(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
