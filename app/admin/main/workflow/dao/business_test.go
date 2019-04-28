package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoinit(t *testing.T) {
	convey.Convey("init", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestDaoBatchLastBusRecIDs(t *testing.T) {
	convey.Convey("BatchLastBusRecIDs", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			oids     = []int64{1}
			business = int8(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			bids, err := d.BatchLastBusRecIDs(c, oids, business)
			ctx.Convey("Then err should be nil.bids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(bids, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBusinessRecs(t *testing.T) {
	convey.Convey("BusinessRecs", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			bids = []int64{1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			bs, err := d.BusinessRecs(c, bids)
			ctx.Convey("Then err should be nil.bs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(bs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoLastBusRec(t *testing.T) {
	convey.Convey("LastBusRec", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			business = int8(1)
			oid      = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			bs, err := d.LastBusRec(c, business, oid)
			ctx.Convey("Then err should be nil.bs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(bs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBatchBusRecByCids(t *testing.T) {
	convey.Convey("BatchBusRecByCids", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			cids = []int64{1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			cidToBus, err := d.BatchBusRecByCids(c, cids)
			ctx.Convey("Then err should be nil.cidToBus should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(cidToBus, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBusObjectByGids(t *testing.T) {
	convey.Convey("BusObjectByGids", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			gids = []int64{1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			gidToBus, err := d.BusObjectByGids(c, gids)
			ctx.Convey("Then err should be nil.gidToBus should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(gidToBus, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAllMetas(t *testing.T) {
	convey.Convey("AllMetas", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.AllMetas(c)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoLoadRole(t *testing.T) {
	convey.Convey("LoadRole", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			role, err := d.LoadRole(c)
			ctx.Convey("Then err should be nil.role should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(role, convey.ShouldNotBeNil)
			})
		})
	})
}
