package dao

import (
	"context"
	"go-common/app/interface/main/credit/model"
	"math/rand"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddBlockedCases(t *testing.T) {
	convey.Convey("AddBlockedCases", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			bc = []*model.ArgJudgeCase{}
			b  = &model.ArgJudgeCase{MID: 1}
		)
		bc = append(bc, b)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.AddBlockedCases(c, bc)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoInsVote(t *testing.T) {
	convey.Convey("InsVote", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = rand.Int63n(99999999)
			cid = rand.Int63n(99999999)
			no  = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.InsVote(c, mid, cid, no)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetvote(t *testing.T) {
	convey.Convey("Setvote", t, func(convCtx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(0)
			cid  = int64(0)
			vote = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.Setvote(c, mid, cid, vote)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetVoteTx(t *testing.T) {
	convey.Convey("SetVoteTx", t, func(convCtx convey.C) {
		var (
			tx, _ = d.BeginTran(context.Background())
			mid   = rand.Int63n(99999999)
			cid   = rand.Int63n(99999999)
			vote  = int8(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			affect, err := d.SetVoteTx(tx, mid, cid, vote)
			if err == nil {
				if err = tx.Commit(); err != nil {
					tx.Rollback()
				}
			} else {
				tx.Rollback()
			}
			convCtx.Convey("Then err should be nil.affect should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(affect, convey.ShouldBeGreaterThanOrEqualTo, 0)
			})
		})
	})
}

func TestDaoAddCaseReasonApply(t *testing.T) {
	convey.Convey("AddCaseReasonApply", t, func(convCtx convey.C) {
		var (
			c            = context.Background()
			mid          = int64(0)
			cid          = int64(0)
			applyType    = int8(0)
			originReason = int8(0)
			applyReason  = int8(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.AddCaseReasonApply(c, mid, cid, applyType, originReason, applyReason)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddCaseVoteTotal(t *testing.T) {
	convey.Convey("AddCaseVoteTotal", t, func(convCtx convey.C) {
		var (
			c       = context.Background()
			field   = "vote_rule"
			cid     = int64(3)
			voteNum = int8(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.AddCaseVoteTotal(c, field, cid, voteNum)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCaseInfo(t *testing.T) {
	convey.Convey("CaseInfo", t, func(convCtx convey.C) {
		var (
			c    = context.Background()
			cid  = int64(3)
			cid1 = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			r, err := d.CaseInfo(c, cid)
			convCtx.Convey("Then err should be nil.r should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(r, convey.ShouldNotBeNil)
			})
			r1, err := d.CaseInfo(c, cid1)
			convCtx.Convey("Then err should be nil.r should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(r1, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCountCaseVote(t *testing.T) {
	convey.Convey("CountCaseVote", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			r, err := d.CountCaseVote(c, mid)
			convCtx.Convey("Then err should be nil.r should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(r, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoIsVote(t *testing.T) {
	convey.Convey("IsVote", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			cid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			r, err := d.IsVote(c, mid, cid)
			convCtx.Convey("Then err should be nil.r should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(r, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoVoteInfo(t *testing.T) {
	convey.Convey("VoteInfo", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			cid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			r, err := d.VoteInfo(c, mid, cid)
			convCtx.Convey("Then err should be nil.r should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(r, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoLoadVoteIDsMid(t *testing.T) {
	convey.Convey("LoadVoteIDsMid", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			day = int(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			cases, err := d.LoadVoteIDsMid(c, mid, day)
			convCtx.Convey("Then err should be nil.cases should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(cases, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCaseVoteIDs(t *testing.T) {
	convey.Convey("CaseVoteIDs", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{3, 6}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			mbc, err := d.CaseVoteIDs(c, ids)
			convCtx.Convey("Then err should be nil.mbc should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(mbc, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCaseRelationIDCount(t *testing.T) {
	convey.Convey("CaseRelationIDCount", t, func(convCtx convey.C) {
		var (
			c          = context.Background()
			tp         = int8(0)
			relationID = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			count, err := d.CaseRelationIDCount(c, tp, relationID)
			convCtx.Convey("Then err should be nil.count should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCaseInfoIDs(t *testing.T) {
	convey.Convey("CaseInfoIDs", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{3}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			cases, err := d.CaseInfoIDs(c, ids)
			convCtx.Convey("Then err should be nil.cases should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(cases, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCaseVotesMID(t *testing.T) {
	convey.Convey("CaseVotesMID", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{1, 44}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			mvo, err := d.CaseVotesMID(c, ids)
			convCtx.Convey("Then err should be nil.mvo should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(mvo, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCaseVoteIDMID(t *testing.T) {
	convey.Convey("CaseVoteIDMID", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(111001692)
			pn  = int64(1)
			ps  = int64(1)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			vids, cids, err := d.CaseVoteIDMID(c, mid, pn, ps)
			convCtx.Convey("Then err should be nil.vids,cids should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(cids, convey.ShouldNotBeNil)
				convCtx.So(vids, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCaseVoteIDTop(t *testing.T) {
	convey.Convey("CaseVoteIDTop", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			vids, cids, err := d.CaseVoteIDTop(c, mid)
			convCtx.Convey("Then err should be nil.vids,cids should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(cids, convey.ShouldNotBeNil)
				convCtx.So(vids, convey.ShouldNotBeNil)
			})
		})
	})
}
