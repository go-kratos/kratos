package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoDelArchiveSpy(t *testing.T) {
	convey.Convey("DelArchiveSpy", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			limit = int64(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.DelArchiveSpy(c, limit)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelUpSpy(t *testing.T) {
	convey.Convey("DelUpSpy", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			limit = int64(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.DelUpSpy(c, limit)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAvBreachRecord(t *testing.T) {
	convey.Convey("AvBreachRecord", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(0)
			limit = int64(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			last, ds, err := d.AvBreachRecord(c, id, limit)
			ctx.Convey("Then err should be nil.last,ds should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ds, convey.ShouldNotBeNil)
				ctx.So(last, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUps(t *testing.T) {
	convey.Convey("Ups", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{1, 2, 3, 4, 5, 6}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			cs, err := d.Ups(c, mids)
			ctx.Convey("Then err should be nil.cs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(cs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAvs(t *testing.T) {
	convey.Convey("Avs", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			date = time.Now()
			aids = []int64{1, 2, 3, 4, 5, 6}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			cs, err := d.Avs(c, date, aids)
			ctx.Convey("Then err should be nil.cs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(cs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoPlayCount(t *testing.T) {
	convey.Convey("PlayCount", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{1, 2, 3, 4, 5, 6}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			cs, err := d.PlayCount(c, mids)
			ctx.Convey("Then err should be nil.cs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(cs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertCheatUps(t *testing.T) {
	convey.Convey("InsertCheatUps", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			values = "(2, '2018-06-23', 'test', 100, 100, 100, 100, 3)"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.InsertCheatUps(c, values)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertCheatArchives(t *testing.T) {
	convey.Convey("InsertCheatArchives", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			values = "(1, 2, 'test', '2018-06-23', 100, 100, 100, 100, 100)"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.InsertCheatArchives(c, values)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
