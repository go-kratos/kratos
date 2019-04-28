package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func testRegexpDaoImplGetByCond(t *testing.T) {
	var (
		c    = context.TODO()
		cond = &Condition{State: "0"}
	)
	convey.Convey("GetByCond", t, func(ctx convey.C) {
		regexps, totalCounts, err := regdi.GetByCond(c, cond)
		ctx.Convey("Then err should be nil.regexps,totalCounts should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(totalCounts, convey.ShouldNotBeNil)
			ctx.So(regexps, convey.ShouldNotBeNil)
		})
	})
}

func testRegexpDaoImplDaoUpdate(t *testing.T) {
	var (
		c = context.TODO()
		r = &Regexp{ID: 1, Name: "name", Area: 1, Content: "test"}
	)
	convey.Convey("Update", t, func(ctx convey.C) {
		p1, err := regdi.Update(c, r)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func testRegexpDaoImplInsert(t *testing.T) {
	var (
		c = context.TODO()
		r = &Regexp{ID: 1, Name: "name", Area: 1, Content: "test"}
	)
	convey.Convey("Insert", t, func(ctx convey.C) {
		p1, err := regdi.Insert(c, r)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func testRegexpDaoImplGetByID(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(1)
	)
	convey.Convey("GetByID", t, func(ctx convey.C) {
		p1, err := regdi.GetByID(c, id)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func testRegexpDaoImplGetByIDs(t *testing.T) {
	var (
		c   = context.TODO()
		ids = []int64{1}
	)
	convey.Convey("RegexpDaoImplGetByIDs", t, func(ctx convey.C) {
		p1, err := regdi.GetByIDs(c, ids)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func testRegexpDaoImplGetByContents(t *testing.T) {
	var (
		c        = context.TODO()
		contents = []string{"test"}
	)
	convey.Convey("GetByContents", t, func(ctx convey.C) {
		p1, err := regdi.GetByContents(c, contents)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func testRegexpDaoImplGetByAreaAndContent(t *testing.T) {
	var (
		c    = context.TODO()
		cond = &Condition{State: "0"}
	)
	convey.Convey("GetByAreaAndContent", t, func(ctx convey.C) {
		p1, err := regdi.GetByAreaAndContent(c, cond)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}
