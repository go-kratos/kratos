package dao

import (
	"context"
	"testing"

	"go-common/app/admin/main/tag/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoRankCount(t *testing.T) {
	var (
		c    = context.TODO()
		prid = int64(0)
		rid  = int64(0)
		tp   = int32(0)
	)
	convey.Convey("RankCount", t, func(ctx convey.C) {
		res, err := d.RankCount(c, prid, rid, tp)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestDaoRankTop(t *testing.T) {
	var (
		c    = context.TODO()
		prid = int64(0)
		rid  = int64(0)
		tp   = int32(0)
	)
	convey.Convey("RankTop", t, func(ctx convey.C) {
		top, topMap, tids, err := d.RankTop(c, prid, rid, tp)
		ctx.Convey("Then err should be nil.top,topMap,tids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(tids, convey.ShouldNotBeNil)
			ctx.So(topMap, convey.ShouldNotBeNil)
			ctx.So(top, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoRankFilter(t *testing.T) {
	var (
		c    = context.TODO()
		prid = int64(0)
		rid  = int64(0)
		tp   = int32(0)
	)
	convey.Convey("RankFilter", t, func(ctx convey.C) {
		filter, filterMap, tids, err := d.RankFilter(c, prid, rid, tp)
		ctx.Convey("Then err should be nil.filter,filterMap,tids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(tids, convey.ShouldHaveLength, 0)
			ctx.So(filterMap, convey.ShouldHaveLength, 0)
			ctx.So(filter, convey.ShouldHaveLength, 0)
		})
	})
}

func TestDaoTxUpdateRankCount(t *testing.T) {
	var (
		id          = int64(0)
		topCount    = int(0)
		filterCount = int(0)
		viewCount   = int(0)
	)
	convey.Convey("TxUpdateRankCount", t, func(ctx convey.C) {
		tx, err := d.BeginTran(context.TODO())
		if err != nil {
			return
		}
		affect, err := d.TxUpdateRankCount(tx, id, topCount, filterCount, viewCount)
		ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affect, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
		tx.Rollback()
	})
}

func TestDaoTxAddRankCount(t *testing.T) {
	var (
		prid        = int64(0)
		rid         = int64(0)
		tp          = int32(0)
		topCount    = int(0)
		filterCount = int(0)
		viewCount   = int(0)
	)
	convey.Convey("TxAddRankCount", t, func(ctx convey.C) {
		tx, err := d.BeginTran(context.TODO())
		if err != nil {
			return
		}
		id, err := d.TxAddRankCount(tx, prid, rid, tp, topCount, filterCount, viewCount)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(id, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
		tx.Rollback()
	})
}

func TestDaoTxRemoveRankTop(t *testing.T) {
	var (
		prid = int64(0)
		rid  = int64(0)
		tp   = int32(0)
	)
	convey.Convey("TxRemoveRankTop", t, func(ctx convey.C) {
		tx, err := d.BeginTran(context.TODO())
		if err != nil {
			return
		}
		affect, err := d.TxRemoveRankTop(tx, prid, rid, tp)
		ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affect, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
		tx.Rollback()
	})
}

func TestDaoTxRemoveRankFilter(t *testing.T) {
	var (
		prid = int64(0)
		rid  = int64(0)
		tp   = int32(0)
	)
	convey.Convey("TxRemoveRankFilter", t, func(ctx convey.C) {
		tx, err := d.BeginTran(context.TODO())
		if err != nil {
			return
		}
		affect, err := d.TxRemoveRankFilter(tx, prid, rid, tp)
		ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affect, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
		tx.Rollback()
	})
}

func TestDaoTxRemoveRankResult(t *testing.T) {
	var (
		prid = int64(0)
		rid  = int64(0)
		tp   = int32(0)
	)
	convey.Convey("TxRemoveRankResult", t, func(ctx convey.C) {
		tx, err := d.BeginTran(context.TODO())
		if err != nil {
			return
		}
		affect, err := d.TxRemoveRankResult(tx, prid, rid, tp)
		ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affect, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
		tx.Rollback()
	})
}

func TestDaoTxInsertRankTop(t *testing.T) {
	var (
		prid    = int64(0)
		rid     = int64(0)
		tp      = int32(0)
		rankTop = []*model.RankTop{
			{
				Tid:   1,
				TName: "ut",
			},
		}
	)
	convey.Convey("TxInsertRankTop", t, func(ctx convey.C) {
		tx, err := d.BeginTran(context.TODO())
		if err != nil {
			return
		}
		id, err := d.TxInsertRankTop(tx, prid, rid, tp, rankTop)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(id, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
		tx.Rollback()
	})
}

func TestDaoTxInsertRankFilter(t *testing.T) {
	var (
		prid       = int64(0)
		rid        = int64(0)
		tp         = int32(0)
		rankFilter = []*model.RankFilter{
			{
				Tid:   1,
				TName: "ut",
			},
		}
	)
	convey.Convey("TxInsertRankFilter", t, func(ctx convey.C) {
		tx, err := d.BeginTran(context.TODO())
		if err != nil {
			return
		}
		id, err := d.TxInsertRankFilter(tx, prid, rid, tp, rankFilter)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(id, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
		tx.Rollback()
	})
}

func TestDaoTxInsertRankResult(t *testing.T) {
	var (
		prid      = int64(0)
		rid       = int64(0)
		tp        = int32(0)
		rankLists = []*model.RankResult{
			{
				Tid:   1,
				TName: "ut",
			},
		}
	)
	convey.Convey("TxInsertRankResult", t, func(ctx convey.C) {
		tx, err := d.BeginTran(context.TODO())
		if err != nil {
			return
		}
		id, err := d.TxInsertRankResult(tx, prid, rid, tp, rankLists)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(id, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
		tx.Rollback()
	})
}
