package dao

import (
	"context"
	"testing"

	"go-common/app/admin/main/tag/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoRelationsByTid(t *testing.T) {
	var (
		c     = context.TODO()
		tid   = int64(0)
		start = int32(0)
		end   = int32(0)
	)
	convey.Convey("RelationsByTid", t, func(ctx convey.C) {
		res, oids, err := d.RelationsByTid(c, tid, start, end)
		ctx.Convey("Then err should be nil.res,oids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(oids, convey.ShouldHaveLength, 0)
			ctx.So(res, convey.ShouldHaveLength, 0)
		})
	})
}

func TestDaoRelationsByOid(t *testing.T) {
	var (
		c     = context.TODO()
		oid   = int64(0)
		tp    = int32(0)
		start = int32(0)
		end   = int32(0)
	)
	convey.Convey("RelationsByOid", t, func(ctx convey.C) {
		res, tids, err := d.RelationsByOid(c, oid, tp, start, end)
		ctx.Convey("Then err should be nil.res,tids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(tids, convey.ShouldHaveLength, 0)
			ctx.So(res, convey.ShouldHaveLength, 0)
		})
	})
}

func TestDaoTagResCount(t *testing.T) {
	var (
		c   = context.TODO()
		tid = int64(0)
	)
	convey.Convey("TagResCount", t, func(ctx convey.C) {
		count, err := d.TagResCount(c, tid)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}

func TestDaoResTagCount(t *testing.T) {
	var (
		c   = context.TODO()
		oid = int64(0)
		tp  = int32(0)
	)
	convey.Convey("ResTagCount", t, func(ctx convey.C) {
		count, err := d.ResTagCount(c, oid, tp)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}

func TestDaoRelation(t *testing.T) {
	var (
		c     = context.TODO()
		oid   = int64(0)
		tid   = int64(0)
		tp    = int32(0)
		state = int32(0)
	)
	convey.Convey("Relation", t, func(ctx convey.C) {
		res, err := d.Relation(c, oid, tid, tp, state)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestDaoTxInsertResTag(t *testing.T) {
	var relation = &model.Relation{
		Oid: 2,
		Tid: 1,
	}
	convey.Convey("TxInsertResTag", t, func(ctx convey.C) {
		tx, err := d.BeginTran(context.TODO())
		if err != nil {
			return
		}
		id, err := d.TxInsertResTag(tx, relation)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(id, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
		tx.Rollback()
	})
}

func TestDaoTxInsertTagRes(t *testing.T) {
	var relation = &model.Relation{
		Oid: 2,
		Tid: 1,
	}
	convey.Convey("TxInsertTagRes", t, func(ctx convey.C) {
		tx, err := d.BeginTran(context.TODO())
		if err != nil {
			return
		}
		id, err := d.TxInsertTagRes(tx, relation)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(id, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
		tx.Rollback()
	})
}

func TestDaoTxUpdateAttrResTag(t *testing.T) {
	var (
		tid  = int64(0)
		oid  = int64(0)
		tp   = int32(0)
		attr = int32(0)
	)
	convey.Convey("TxUpdateAttrResTag", t, func(ctx convey.C) {
		tx, err := d.BeginTran(context.TODO())
		if err != nil {
			return
		}
		affect, err := d.TxUpdateAttrResTag(tx, tid, oid, tp, attr)
		ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affect, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
		tx.Rollback()
	})
}

func TestDaoTxUpdateAttrTagRes(t *testing.T) {
	var (
		tid  = int64(0)
		oid  = int64(0)
		tp   = int32(0)
		attr = int32(0)
	)
	convey.Convey("TxUpdateAttrTagRes", t, func(ctx convey.C) {
		tx, err := d.BeginTran(context.TODO())
		if err != nil {
			return
		}
		rowsCount, err := d.TxUpdateAttrTagRes(tx, tid, oid, tp, attr)
		ctx.Convey("Then err should be nil.rowsCount should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(rowsCount, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
		tx.Rollback()
	})
}

func TestDaoTxUpdateStateResTag(t *testing.T) {
	var (
		tid   = int64(0)
		oid   = int64(0)
		tp    = int32(0)
		state = int32(0)
	)
	convey.Convey("TxUpdateStateResTag", t, func(ctx convey.C) {
		tx, err := d.BeginTran(context.TODO())
		if err != nil {
			return
		}
		affect, err := d.TxUpdateStateResTag(tx, tid, oid, tp, state)
		ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affect, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
		tx.Rollback()
	})
}

func TestDaoTxUpdateStateTagRes(t *testing.T) {
	var (
		tid   = int64(0)
		oid   = int64(0)
		tp    = int32(0)
		state = int32(0)
	)
	convey.Convey("TxUpdateStateTagRes", t, func(ctx convey.C) {
		tx, err := d.BeginTran(context.TODO())
		if err != nil {
			return
		}
		affect, err := d.TxUpdateStateTagRes(tx, tid, oid, tp, state)
		ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affect, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
		tx.Rollback()
	})
}
