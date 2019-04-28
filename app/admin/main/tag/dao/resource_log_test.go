package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoresourceLogSQL(t *testing.T) {
	var (
		role   = int32(0)
		action = int32(0)
	)
	convey.Convey("resourceLogSQL", t, func(ctx convey.C) {
		p1 := resourceLogSQL(role, action)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoResTagLogCount(t *testing.T) {
	var (
		c      = context.TODO()
		oid    = int64(0)
		tp     = int32(0)
		role   = int32(0)
		action = int32(0)
	)
	convey.Convey("ResTagLogCount", t, func(ctx convey.C) {
		count, err := d.ResTagLogCount(c, oid, tp, role, action)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoResourceLogs(t *testing.T) {
	var (
		c      = context.TODO()
		oid    = int64(0)
		tp     = int32(0)
		role   = int32(0)
		action = int32(0)
		start  = int32(0)
		end    = int32(0)
	)
	convey.Convey("ResourceLogs", t, func(ctx convey.C) {
		res, err := d.ResourceLogs(c, oid, tp, role, action, start, end)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldHaveLength, 0)
		})
	})
}

func TestDaoUpdateResLogState(t *testing.T) {
	var (
		c     = context.TODO()
		id    = int64(0)
		oid   = int64(0)
		tp    = int32(0)
		state = int32(0)
	)
	convey.Convey("UpdateResLogState", t, func(ctx convey.C) {
		affect, err := d.UpdateResLogState(c, id, oid, tp, state)
		ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affect, convey.ShouldNotBeNil)
		})
	})
}
