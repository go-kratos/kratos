package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAdvanceType(t *testing.T) {
	convey.Convey("AdvanceType", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			cid  = int64(0)
			mid  = int64(0)
			mode = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			typ, err := testDao.AdvanceType(c, cid, mid, mode)
			ctx.Convey("Then err should be nil.typ should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(typ, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAdvance(t *testing.T) {
	convey.Convey("Advance", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			id  = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.Advance(c, mid, id)
		})
	})
}

func TestDaoAdvances(t *testing.T) {
	convey.Convey("Advances", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			owner = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := testDao.Advances(c, owner)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBuyAdvance(t *testing.T) {
	convey.Convey("BuyAdvance", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(0)
			cid    = int64(0)
			owner  = int64(0)
			refund = int64(0)
			typ    = ""
			mode   = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			id, err := testDao.BuyAdvance(c, mid, cid, owner, refund, typ, mode)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateAdvType(t *testing.T) {
	convey.Convey("UpdateAdvType", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(0)
			typ = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			affect, err := testDao.UpdateAdvType(c, id, typ)
			ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affect, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelAdvance(t *testing.T) {
	convey.Convey("DelAdvance", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			affect, err := testDao.DelAdvance(c, id)
			ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affect, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAdvanceCmt(t *testing.T) {
	convey.Convey("AdvanceCmt", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			mid  = int64(0)
			mode = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			adv, err := testDao.AdvanceCmt(c, oid, mid, mode)
			ctx.Convey("Then err should be nil.adv should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(adv, convey.ShouldNotBeNil)
			})
		})
	})
}
