package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/member/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoMoralLog(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(4780461)
		// ip  = ""
	)
	convey.Convey("MoralLog", t, func(ctx convey.C) {
		p1, p2 := d.MoralLog(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(p2, convey.ShouldBeNil)
		})
		ctx.Convey("p1 should not be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoMoralLogByID(t *testing.T) {
	var (
		c     = context.Background()
		logID = "test"
		// ip    = ""
	)
	convey.Convey("MoralLogByID", t, func(ctx convey.C) {
		p1, p2 := d.MoralLogByID(c, logID)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(p2, convey.ShouldEqual, -404)
		})
		ctx.Convey("p1 should not be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldBeNil)
		})
	})
}

func TestDaoDeleteMoralLog(t *testing.T) {
	var (
		c     = context.Background()
		logID = "test"
		// ip    = ""
	)
	convey.Convey("DeleteMoralLog", t, func(ctx convey.C) {
		p1 := d.DeleteMoralLog(c, logID)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldBeNil)
		})
	})
}

func TestDaodeleteLogReport(t *testing.T) {
	var (
		c        = context.Background()
		business = int(0)
		logID    = "test"
		ip       = ""
	)
	convey.Convey("deleteLogReport", t, func(ctx convey.C) {
		p1 := d.deleteLogReport(c, business, logID, ip)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldBeNil)
		})
	})
}

func TestDaoasMoralLog(t *testing.T) {
	var (
		res = &model.SearchResult{}
	)
	convey.Convey("asMoralLog", t, func(ctx convey.C) {
		p1 := asMoralLog(res)
		ctx.Convey("p1 should not be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}
func TestDaoMoral(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(4780461)
	)
	convey.Convey("Moral", t, func(ctx convey.C) {
		moral, err := d.Moral(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("moral should not be nil", func(ctx convey.C) {
			ctx.So(moral, convey.ShouldNotBeNil)
		})
	})
}
