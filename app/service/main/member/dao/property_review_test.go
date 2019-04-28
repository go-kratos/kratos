package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/member/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddUserMonitor(t *testing.T) {
	var (
		c        = context.TODO()
		mid      = int64(1)
		operator = "zhoujiahui"
		remark   = "test"
	)
	convey.Convey("AddUserMonitor", t, func(ctx convey.C) {
		p1 := d.AddUserMonitor(c, mid, operator, remark)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldBeNil)
		})
	})
}

func TestDaoIsInUserMonitor(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(1)
	)
	convey.Convey("IsInUserMonitor", t, func(ctx convey.C) {
		p1, p2 := d.IsInUserMonitor(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(p2, convey.ShouldBeNil)
		})
		ctx.Convey("p1 should not be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddPropertyReview(t *testing.T) {
	var (
		c = context.TODO()
		r = &model.UserPropertyReview{
			Mid:       2231365,
			Old:       "hahhah",
			New:       "dangerou",
			State:     model.ReviewStateWait,
			Property:  model.ReviewPropertySign,
			IsMonitor: true,
			Extra:     "{}",
		}
	)
	convey.Convey("AddPropertyReview", t, func(ctx convey.C) {
		p1 := d.AddPropertyReview(c, r)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldBeNil)
		})
	})
}

func TestDaoArchivePropertyReview(t *testing.T) {
	var (
		c        = context.TODO()
		mid      = int64(3)
		property = int8(1)
	)
	convey.Convey("ArchivePropertyReview", t, func(ctx convey.C) {
		p1 := d.ArchivePropertyReview(c, mid, property)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldBeNil)
		})
	})
}
