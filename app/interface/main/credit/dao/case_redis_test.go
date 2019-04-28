package dao

import (
	"context"
	"go-common/app/interface/main/credit/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaovoteIndexKey(t *testing.T) {
	convey.Convey("voteIndexKey", t, func(convCtx convey.C) {
		var (
			cid   = int64(0)
			otype = int8(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := voteIndexKey(cid, otype)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaocaseIndexKey(t *testing.T) {
	convey.Convey("caseIndexKey", t, func(convCtx convey.C) {
		var (
			cid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := caseIndexKey(cid)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoVoteOpIdxCache(t *testing.T) {
	convey.Convey("VoteOpIdxCache", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			cid   = int64(0)
			start = int64(0)
			end   = int64(0)
			otype = int8(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			ids, err := d.VoteOpIdxCache(c, cid, start, end, otype)
			convCtx.Convey("Then err should be nil.ids should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(ids, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoExpireVoteIdx(t *testing.T) {
	convey.Convey("ExpireVoteIdx", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			cid   = int64(0)
			otype = int8(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			ok, err := d.ExpireVoteIdx(c, cid, otype)
			convCtx.Convey("Then err should be nil.ok should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoLenVoteIdx(t *testing.T) {
	convey.Convey("LenVoteIdx", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			cid   = int64(0)
			otype = int8(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			count, err := d.LenVoteIdx(c, cid, otype)
			convCtx.Convey("Then err should be nil.count should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCaseOpIdxCache(t *testing.T) {
	convey.Convey("CaseOpIdxCache", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			cid   = int64(0)
			start = int64(0)
			end   = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			ids, err := d.CaseOpIdxCache(c, cid, start, end)
			convCtx.Convey("Then err should be nil.ids should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(ids, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoLenCaseIdx(t *testing.T) {
	convey.Convey("LenCaseIdx", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			cid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			count, err := d.LenCaseIdx(c, cid)
			convCtx.Convey("Then err should be nil.count should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoExpireCaseIdx(t *testing.T) {
	convey.Convey("ExpireCaseIdx", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			cid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			ok, err := d.ExpireCaseIdx(c, cid)
			convCtx.Convey("Then err should be nil.ok should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoLoadVoteOpIdxs(t *testing.T) {
	convey.Convey("LoadVoteOpIdxs", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			cid   = int64(0)
			otype = int8(0)
			idx   = []int64{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.LoadVoteOpIdxs(c, cid, otype, idx)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoLoadCaseIdxs(t *testing.T) {
	convey.Convey("LoadCaseIdxs", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			cid = int64(0)
			ops = []*model.Opinion{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.LoadCaseIdxs(c, cid, ops)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelCaseIdx(t *testing.T) {
	convey.Convey("DelCaseIdx", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			cid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.DelCaseIdx(c, cid)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelVoteIdx(t *testing.T) {
	convey.Convey("DelVoteIdx", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			cid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.DelVoteIdx(c, cid)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

// TestVoteOpIdxCache .
func TestVoteOpIdxCache(t *testing.T) {
	var (
		idx         = []int64{1, 2, 3, 4}
		cid   int64 = 639
		c           = context.TODO()
		otype int8  = 1
	)
	convey.Convey("return someting", t, func(convCtx convey.C) {
		d.LoadVoteOpIdxs(c, cid, otype, idx)
		ids, err := d.VoteOpIdxCache(c, cid, 1, 3, otype)
		convey.So(err, convey.ShouldBeNil)
		convey.So(ids, convey.ShouldNotBeNil)
		convey.Convey("get vote count", func(convCtx convey.C) {
			count, err := d.LenVoteIdx(c, cid, otype)
			convey.So(err, convey.ShouldBeNil)
			convey.So(count, convey.ShouldEqual, 4)
		})
	})
}

// TestCaseOpIdxCache .
func TestCaseOpIdxCache(t *testing.T) {
	var (
		c         = context.TODO()
		cid int64 = 2
	)
	convey.Convey("return someting", t, func(ctx convey.C) {
		ops, err := d.Opinions(c, []int64{631, 633})
		convey.So(err, convey.ShouldBeNil)
		convey.So(ops, convey.ShouldNotBeNil)
		convey.Convey("gets case data from cache", func(ctx convey.C) {
			d.LoadCaseIdxs(c, cid, ops)
			ids, err := d.CaseOpIdxCache(c, cid, 0, 3)
			convey.So(err, convey.ShouldBeNil)
			convey.So(ids, convey.ShouldNotBeNil)
		})
		convey.Convey("count cid from case cache", func(ctx convey.C) {
			count, err := d.LenCaseIdx(c, cid)
			convey.So(err, convey.ShouldBeNil)
			convey.So(count, convey.ShouldEqual, len(ops))
		})
	})
}
