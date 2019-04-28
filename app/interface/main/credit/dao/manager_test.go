package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoManagers(t *testing.T) {
	convey.Convey("Managers", t, func(convCtx convey.C) {
		var (
			c = context.Background()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			manMap, err := d.Managers(c)
			convCtx.Convey("Then err should be nil.manMap should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(manMap, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoManager(t *testing.T) {
	convey.Convey("Manager", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			pn = int64(0)
			ps = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			mi, err := d.Manager(c, pn, ps)
			convCtx.Convey("Then err should be nil.mi should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(mi, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoManagerTotal(t *testing.T) {
	convey.Convey("ManagerTotal", t, func(convCtx convey.C) {
		var (
			c = context.Background()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			count, err := d.ManagerTotal(c)
			convCtx.Convey("Then err should be nil.count should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}
