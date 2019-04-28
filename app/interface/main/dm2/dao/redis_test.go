package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyDM(t *testing.T) {
	convey.Convey("keyDM", t, func(ctx convey.C) {
		var (
			tp  = int32(0)
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			key := keyDM(tp, oid)
			ctx.Convey("Then key should not be nil.", func(ctx convey.C) {
				ctx.So(key, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyBroadcastLimit(t *testing.T) {
	convey.Convey("keyBroadcastLimit", t, func(ctx convey.C) {
		var (
			tp  = int32(0)
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			key := keyBroadcastLimit(tp, oid)
			ctx.Convey("Then key should not be nil.", func(ctx convey.C) {
				ctx.So(key, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDMCache(t *testing.T) {
	convey.Convey("DMCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tp  = int32(0)
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := testDao.DMCache(c, tp, oid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoExpireDMCache(t *testing.T) {
	convey.Convey("ExpireDMCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tp  = int32(0)
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ok, err := testDao.ExpireDMCache(c, tp, oid)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTrimDMCache(t *testing.T) {
	convey.Convey("TrimDMCache", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tp    = int32(0)
			oid   = int64(0)
			count = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.TrimDMCache(c, tp, oid, count)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoIncrPubCnt(t *testing.T) {
	convey.Convey("IncrPubCnt", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(0)
			color    = int64(0)
			mode     = int32(0)
			fontsize = int32(0)
			ip       = ""
			msg      = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.IncrPubCnt(c, mid, color, mode, fontsize, ip, msg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoPubCnt(t *testing.T) {
	convey.Convey("PubCnt", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(0)
			color    = int64(0)
			mode     = int32(0)
			fontsize = int32(0)
			ip       = ""
			msg      = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			count, err := testDao.PubCnt(c, mid, color, mode, fontsize, ip, msg)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoIncrCharPubCnt(t *testing.T) {
	convey.Convey("IncrCharPubCnt", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.IncrCharPubCnt(c, mid, oid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCharPubCnt(t *testing.T) {
	convey.Convey("CharPubCnt", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			count, err := testDao.CharPubCnt(c, mid, oid)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelCharPubCnt(t *testing.T) {
	convey.Convey("DelCharPubCnt", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.DelCharPubCnt(c, mid, oid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoBroadcastLimit(t *testing.T) {
	convey.Convey("BroadcastLimit", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			oid      = int64(0)
			tp       = int32(0)
			count    = int(0)
			interval = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.BroadcastLimit(c, oid, tp, count, interval)
		})
	})
}
