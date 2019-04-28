package dao

import (
	"context"
	"go-common/app/interface/main/credit/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaocaseInfoKey(t *testing.T) {
	convey.Convey("caseInfoKey", t, func(convCtx convey.C) {
		var (
			cid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := caseInfoKey(cid)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaovoteCaseInfoKey(t *testing.T) {
	convey.Convey("voteCaseInfoKey", t, func(convCtx convey.C) {
		var (
			mid = int64(0)
			cid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := voteCaseInfoKey(mid, cid)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaocaseVoteTopKey(t *testing.T) {
	convey.Convey("caseVoteTopKey", t, func(convCtx convey.C) {
		var (
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := caseVoteTopKey(mid)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSetCaseInfoCache(t *testing.T) {
	convey.Convey("SetCaseInfoCache", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			cid = int64(0)
			bc  = &model.BlockedCase{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.SetCaseInfoCache(c, cid, bc)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCaseInfoCache(t *testing.T) {
	convey.Convey("CaseInfoCache", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			cid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			bc, err := d.CaseInfoCache(c, cid)
			convCtx.Convey("Then err should be nil.bc should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(bc, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSetVoteInfoCache(t *testing.T) {
	convey.Convey("SetVoteInfoCache", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			cid = int64(0)
			vi  = &model.VoteInfo{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.SetVoteInfoCache(c, mid, cid, vi)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoVoteInfoCache(t *testing.T) {
	convey.Convey("VoteInfoCache", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			cid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			vi, err := d.VoteInfoCache(c, mid, cid)
			convCtx.Convey("Then err should be nil.vi should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(vi, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCaseVoteTopCache(t *testing.T) {
	convey.Convey("CaseVoteTopCache", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(-1)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			bs, err := d.CaseVoteTopCache(c, mid)
			convCtx.Convey("Then err should be nil.bs should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(len(bs), convey.ShouldBeGreaterThanOrEqualTo, 0)
			})
		})
	})
}

func TestDaoSetCaseVoteTopCache(t *testing.T) {
	convey.Convey("SetCaseVoteTopCache", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			bs  = []*model.BlockedCase{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.SetCaseVoteTopCache(c, mid, bs)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
