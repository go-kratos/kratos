package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoTaskStatus(t *testing.T) {
	convey.Convey("TaskStatus", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			date = "2018-12-01"
			typ  = int(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			status, err := d.TaskStatus(c, date, typ)
			ctx.Convey("Then err should be nil.status should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(status, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateTaskStatus(t *testing.T) {
	convey.Convey("UpdateTaskStatus", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			date   = "2018-12-01"
			typ    = int(1)
			status = int(2)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.UpdateTaskStatus(c, date, typ, status)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertTaskStatus(t *testing.T) {
	convey.Convey("InsertTaskStatus", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			typ     = int(1)
			status  = int(2)
			date    = "2018-12-01"
			message = "ok"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.InsertTaskStatus(c, typ, status, date, message)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
