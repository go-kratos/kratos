package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func testRuleDaoImplGetByCond(t *testing.T) {
	var (
		c    = context.TODO()
		cond = &Condition{}
	)
	convey.Convey("GetByCond", t, func(ctx convey.C) {
		rules, totalCounts, err := rdi.GetByCond(c, cond)
		ctx.Convey("Then err should be nil.rules,totalCounts should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(totalCounts, convey.ShouldNotBeNil)
			ctx.So(rules, convey.ShouldNotBeNil)
		})
	})
}

func testRuleDaoImplUpdate(t *testing.T) {
	var (
		c = context.TODO()
		r = &Rule{}
	)
	convey.Convey("Update", t, func(ctx convey.C) {
		p1, err := rdi.Update(c, r)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func testRuleDaoImplInsert(t *testing.T) {
	var (
		c = context.TODO()
		r = &Rule{}
	)
	convey.Convey("Insert", t, func(ctx convey.C) {
		p1, err := rdi.Insert(c, r)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func testRuleDaoImplGetByID(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(0)
	)
	convey.Convey("GetByID", t, func(ctx convey.C) {
		p1, err := rdi.GetByID(c, id)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func testRuleDaoImplDaoGetByIDs(t *testing.T) {
	var (
		c   = context.TODO()
		ids = []int64{}
	)
	convey.Convey("GetByIDs", t, func(ctx convey.C) {
		p1, err := rdi.GetByIDs(c, ids)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func testRuleDaoImplGetByAreaAndLimitType(t *testing.T) {
	var (
		c    = context.TODO()
		cond = &Condition{}
	)
	convey.Convey("GetByAreaAndLimitType", t, func(ctx convey.C) {
		p1, err := rdi.GetByAreaAndLimitType(c, cond)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func testRuleDaoImplGetByAreaAndTypeAndScope(t *testing.T) {
	var (
		c    = context.TODO()
		cond = &Condition{}
	)
	convey.Convey("GetByAreaAndTypeAndScope", t, func(ctx convey.C) {
		p1, err := rdi.GetByAreaAndTypeAndScope(c, cond)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func testRuleDaoImplGetByArea(t *testing.T) {
	var (
		c    = context.TODO()
		cond = &Condition{}
	)
	convey.Convey("GetByArea", t, func(ctx convey.C) {
		p1, err := rdi.GetByArea(c, cond)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}
