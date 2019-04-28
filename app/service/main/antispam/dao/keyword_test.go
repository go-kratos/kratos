package dao

import (
	"context"
	"fmt"
	"go-common/app/service/main/antispam/util"
	"math/rand"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func testKeywordDaoImplGetRubbish(t *testing.T) {
	var (
		c    = context.TODO()
		cond = &Condition{Pagination: &util.Pagination{CurPage: 1, PerPage: 10}, Tags: []string{"reply"}, Area: "reply", Offset: "1", State: "0", HitCounts: "0", StartTime: "2018-8-1 16:36:48", EndTime: "2018-8-21 16:36:48"}
	)
	convey.Convey("GetRubbish", t, func(ctx convey.C) {
		_, err := kwi.GetRubbish(c, cond)
		ctx.Convey("Then err should be nil.keywords should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func testKeywordDaoImplGetByOffsetLimit(t *testing.T) {
	var (
		c    = context.TODO()
		cond = &Condition{Pagination: &util.Pagination{CurPage: 1, PerPage: 10}, Tags: []string{"reply"}, Area: "reply", Offset: "1", State: "0", HitCounts: "0", StartTime: "2018-8-1 16:36:48", EndTime: "2018-8-21 16:36:48"}
	)
	convey.Convey("GetByOffsetLimit", t, func(ctx convey.C) {
		keywords, err := kwi.GetByOffsetLimit(c, cond)
		ctx.Convey("Then err should be nil.keywords should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(keywords, convey.ShouldNotBeNil)
		})
	})
}

func testKeywordDaoImplGetByCond(t *testing.T) {
	var (
		c    = context.TODO()
		cond = &Condition{Pagination: &util.Pagination{CurPage: 1, PerPage: 10}, Tags: []string{"reply"}, Area: "reply", Offset: "1", State: "0", HitCounts: "0", StartTime: "2018-8-1 16:36:48", EndTime: "2018-8-21 16:36:48"}
	)
	convey.Convey("GetByCond", t, func(ctx convey.C) {
		keywords, totalCounts, err := kwi.GetByCond(c, cond)
		ctx.Convey("Then err should be nil.keywords,totalCounts should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(totalCounts, convey.ShouldNotBeNil)
			ctx.So(keywords, convey.ShouldNotBeNil)
		})
	})
}

func testKeywordDaoImplGetByAreaAndContents(t *testing.T) {
	var (
		c    = context.TODO()
		cond = &Condition{Pagination: &util.Pagination{CurPage: 1, PerPage: 10}, Tags: []string{"reply"}, Area: "reply", Offset: "1", State: "0", HitCounts: "0", StartTime: "2018-8-1 16:36:48", EndTime: "2018-8-21 16:36:48"}
	)
	convey.Convey("GetByAreaAndContents", t, func(ctx convey.C) {
		p1, err := kwi.GetByAreaAndContents(c, cond)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func testKeywordDaoImplGetByAreaAndContent(t *testing.T) {
	var (
		c    = context.TODO()
		cond = &Condition{Pagination: &util.Pagination{CurPage: 1, PerPage: 10}, Tags: []string{"reply"}, Area: "reply", Offset: "1", State: "0", HitCounts: "0", StartTime: "2018-8-1 16:36:48", EndTime: "2018-8-21 16:36:48"}
	)
	convey.Convey("GetByAreaAndContent", t, func(ctx convey.C) {
		p1, err := kwi.GetByAreaAndContent(c, cond)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestKeywordDaoImplUpdate(t *testing.T) {
	var (
		c = context.TODO()
		k = &Keyword{ID: 1, Content: fmt.Sprint(rand.Int63())}
	)
	convey.Convey("Update", t, func(ctx convey.C) {
		p1, err := kwi.Update(c, k)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestKeywordDaoImplInsert(t *testing.T) {
	var (
		c = context.TODO()
		k = &Keyword{Content: fmt.Sprint(rand.Int63())}
	)
	convey.Convey("Insert", t, func(ctx convey.C) {
		p1, err := kwi.Insert(c, k)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestKeywordDaoImplDeleteByIDs(t *testing.T) {
	var (
		c   = context.TODO()
		ids = []int64{1}
	)
	convey.Convey("DeleteByIDs", t, func(ctx convey.C) {
		p1, err := kwi.DeleteByIDs(c, ids)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestKeywordDaoImplGetByID(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(1)
	)
	convey.Convey("GetByID", t, func(ctx convey.C) {
		p1, err := kwi.GetByID(c, id)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestKeywordDaoImplGetByIDs(t *testing.T) {
	var (
		c   = context.TODO()
		ids = []int64{1, 2, 3}
	)
	convey.Convey("GetByIDs", t, func(ctx convey.C) {
		p1, err := kwi.GetByIDs(c, ids)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}
