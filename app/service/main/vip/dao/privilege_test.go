package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

// go test  -test.v -test.run TestDaoPrivilegeList
func TestDaoPrivilegeList(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("PrivilegeList", t, func(ctx convey.C) {
		res, err := d.PrivilegeList(c)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("res should not be nil", func(ctx convey.C) {
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoPrivilegeResourcesList(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("PrivilegeResourcesList", t, func(ctx convey.C) {
		_, err := d.PrivilegeResourcesList(c)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
