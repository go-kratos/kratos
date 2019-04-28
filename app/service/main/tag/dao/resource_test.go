package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/tag/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoTxAddResource(t *testing.T) {
	convey.Convey("TxAddResource", t, func(ctx convey.C) {
		r := &model.Resource{
			Oid:  28843596,
			Tid:  1833,
			Type: 3,
		}
		tx, err := d.BeginTran(context.Background())
		if err != nil {
			return
		}
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.TxAddResource(tx, r)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldBeGreaterThanOrEqualTo, 0)
			})
		})
		tx.Commit()
	})
}

func TestDaoTxDelResource(t *testing.T) {
	convey.Convey("TxDelResource", t, func(ctx convey.C) {
		r := &model.Resource{
			Oid:  28843596,
			Tid:  1833,
			Type: 3,
		}
		tx, err := d.BeginTran(context.Background())
		if err != nil {
			return
		}
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.TxDelResource(tx, r)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldBeGreaterThanOrEqualTo, 0)
			})
		})
		tx.Commit()
	})
}

func TestDaoTxAddDefaultResource(t *testing.T) {
	convey.Convey("TxAddDefaultResource", t, func(ctx convey.C) {
		r := &model.Resource{
			Oid:  28843596,
			Tid:  1835,
			Type: 3,
		}
		tx, err := d.BeginTran(context.Background())
		if err != nil {
			return
		}
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.TxAddDefaultResource(tx, r)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldBeGreaterThanOrEqualTo, 0)
			})
		})
		tx.Commit()
	})
}

func TestDaoResourceDefault(t *testing.T) {
	convey.Convey("ResourceDefault", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(28843596)
			typ = int32(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ResourceDefault(c, oid, typ)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoResource(t *testing.T) {
	convey.Convey("Resource", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(28843596)
			tid = int64(1833)
			typ = int32(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			r, err := d.Resource(c, oid, tid, typ)
			ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(r, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoResources(t *testing.T) {
	convey.Convey("Resources", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(28843596)
			typ = int32(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Resources(c, oid, typ)
		})
	})
}

func TestDaoResourceMap(t *testing.T) {
	convey.Convey("ResourceMap", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(28843596)
			typ = int32(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ResourceMap(c, oid, typ)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoResOidsByTid(t *testing.T) {
	convey.Convey("ResOidsByTid", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tid   = int64(0)
			limit = int64(0)
			typ   = int32(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.ResOidsByTid(c, tid, limit, typ)
		})
	})
}

func TestDaoResourcesByTid(t *testing.T) {
	convey.Convey("ResourcesByTid", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tid   = int64(1833)
			limit = int64(10)
			typ   = int32(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ResourcesByTid(c, tid, limit, typ)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoResAllOidByTid(t *testing.T) {
	convey.Convey("ResAllOidByTid", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			tid    = int64(1833)
			typeID = int32(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			oids, err := d.ResAllOidByTid(c, tid, typeID)
			ctx.Convey("Then err should be nil.oids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(oids, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxUpResLike(t *testing.T) {
	convey.Convey("TxUpResLike", t, func(ctx convey.C) {
		var (
			oid  = int64(28843596)
			tid  = int64(1833)
			typ  = int32(3)
			incr = int32(0)
		)
		tx, err := d.BeginTran(context.Background())
		if err != nil {
			return
		}
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.TxUpResLike(tx, oid, tid, typ, incr)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldBeGreaterThanOrEqualTo, 0)
			})
		})
		tx.Commit()
	})
}

func TestDaoTxUpResHate(t *testing.T) {
	convey.Convey("TxUpResHate", t, func(ctx convey.C) {
		var (
			oid  = int64(28843596)
			tid  = int64(1833)
			typ  = int32(3)
			incr = int32(0)
		)
		tx, err := d.BeginTran(context.Background())
		if err != nil {
			return
		}
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.TxUpResHate(tx, oid, tid, typ, incr)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldBeGreaterThanOrEqualTo, 0)
			})
		})
		tx.Commit()
	})
}

func TestDaoTxUpResAttr(t *testing.T) {
	convey.Convey("TxUpResAttr", t, func(ctx convey.C) {
		var (
			oid = int64(28843596)
			tid = int64(1833)
			bit = uint(1)
			val = int32(1)
			typ = int32(3)
		)
		tx, err := d.BeginTran(context.Background())
		if err != nil {
			return
		}
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.TxUpResAttr(tx, oid, tid, bit, val, typ)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldBeGreaterThanOrEqualTo, 0)
			})
		})
		tx.Commit()
	})
}
