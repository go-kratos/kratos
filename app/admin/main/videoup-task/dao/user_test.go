package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetUsernameAndRole(t *testing.T) {
	var (
		c    = context.TODO()
		uids = []int64{1, 74, 241}
	)
	convey.Convey("GetUsernameAndRole", t, func(ctx convey.C) {
		list, err := d.GetUsernameAndRole(c, uids)
		ctx.Convey("Then err should be nil.list should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(list, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoGetUsernameAndDepartment(t *testing.T) {
	var (
		c    = context.TODO()
		uids = []int64{1, 74, 241}
	)
	convey.Convey("GetUsernameAndDepartment", t, func(ctx convey.C) {
		list, err := d.GetUsernameAndDepartment(c, uids)
		ctx.Convey("Then err should be nil.list should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(list, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoGetUsername(t *testing.T) {
	var (
		c    = context.TODO()
		uids = []int64{1, 74, 241}
	)
	convey.Convey("GetUsername", t, func(ctx convey.C) {
		list, err := d.GetUsername(c, uids)
		ctx.Convey("Then err should be nil.list should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(list, convey.ShouldNotBeNil)
		})
	})
}
