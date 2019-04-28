package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddOpinionTx(t *testing.T) {
	convey.Convey("AddOpinionTx", t, func(convCtx convey.C) {
		var (
			tx, _   = d.BeginTran(context.Background())
			cid     = int64(0)
			opid    = int64(0)
			mid     = int64(0)
			content = ""
			attr    = int8(0)
			vote    = int8(0)
			state   = int8(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			affect, err := d.AddOpinionTx(tx, cid, opid, mid, content, attr, vote, state)
			if err == nil {
				if err = tx.Commit(); err != nil {
					tx.Rollback()
				}
			} else {
				tx.Rollback()
			}
			convCtx.Convey("Then err should be nil.affect should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(affect, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddLikes(t *testing.T) {
	convey.Convey("AddLikes", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{53, 55}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			affect, err := d.AddLikes(c, ids)
			convCtx.Convey("Then err should be nil.affect should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(affect, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddHates(t *testing.T) {
	convey.Convey("AddHates", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{53, 55}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			affect, err := d.AddHates(c, ids)
			convCtx.Convey("Then err should be nil.affect should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(affect, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelOpinion(t *testing.T) {
	convey.Convey("DelOpinion", t, func(convCtx convey.C) {
		var (
			c    = context.Background()
			opid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.DelOpinion(c, opid)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoOpinionIdx(t *testing.T) {
	convey.Convey("OpinionIdx", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			cid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			ops, err := d.OpinionIdx(c, cid)
			convCtx.Convey("Then err should be nil.ops should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(ops, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoOpinionCaseIdx(t *testing.T) {
	convey.Convey("OpinionCaseIdx", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			cid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			ops, err := d.OpinionCaseIdx(c, cid)
			convCtx.Convey("Then err should be nil.ops should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(ops, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoOpinions(t *testing.T) {
	convey.Convey("Opinions", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			opids = []int64{55, 207}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			ops, err := d.Opinions(c, opids)
			convCtx.Convey("Then err should be nil.ops should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(ops, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoOpContentMid(t *testing.T) {
	convey.Convey("OpContentMid", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			content, err := d.OpContentMid(c, mid)
			convCtx.Convey("Then err should be nil.content should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(content, convey.ShouldNotBeNil)
			})
		})
	})
}
