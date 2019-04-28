package reply

import (
	"context"
	"testing"
	"time"

	"go-common/app/job/main/reply/model/reply"

	"github.com/smartystreets/goconvey/convey"
)

func TestReplykeyDialogIdx(t *testing.T) {
	convey.Convey("keyDialogIdx", t, func(ctx convey.C) {
		var (
			dialogID = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyDialogIdx(dialogID)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplykeyIdx(t *testing.T) {
	convey.Convey("keyIdx", t, func(ctx convey.C) {
		var (
			oid  = int64(0)
			tp   = int8(0)
			sort = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyIdx(oid, tp, sort)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplykeyNewRtIdx(t *testing.T) {
	convey.Convey("keyNewRtIdx", t, func(ctx convey.C) {
		var (
			rpID = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyNewRtIdx(rpID)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplykeyAuditIdx(t *testing.T) {
	convey.Convey("keyAuditIdx", t, func(ctx convey.C) {
		var (
			oid = int64(0)
			tp  = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyAuditIdx(oid, tp)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplykeyRpt(t *testing.T) {
	convey.Convey("keyRpt", t, func(ctx convey.C) {
		var (
			mid = int64(0)
			now = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyRpt(mid, now)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplykeyLike(t *testing.T) {
	convey.Convey("keyLike", t, func(ctx convey.C) {
		var (
			rpID = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyLike(rpID)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplykeyUAct(t *testing.T) {
	convey.Convey("keyUAct", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyUAct(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplykeySpamRpRec(t *testing.T) {
	convey.Convey("keySpamRpRec", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keySpamRpRec(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplykeySpamRpDaily(t *testing.T) {
	convey.Convey("keySpamRpDaily", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keySpamRpDaily(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplykeySpamActRec(t *testing.T) {
	convey.Convey("keySpamActRec", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keySpamActRec(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplykeyTopOid(t *testing.T) {
	convey.Convey("keyTopOid", t, func(ctx convey.C) {
		var (
			tp = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyTopOid(tp)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplykeyNotifyCnt(t *testing.T) {
	convey.Convey("keyNotifyCnt", t, func(ctx convey.C) {
		var (
			oid = int64(0)
			typ = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyNotifyCnt(oid, typ)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplykeyMaxLikeCnt(t *testing.T) {
	convey.Convey("keyMaxLikeCnt", t, func(ctx convey.C) {
		var (
			rpid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyMaxLikeCnt(rpid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyDelAuditIndexs(t *testing.T) {
	convey.Convey("DelAuditIndexs", t, func(ctx convey.C) {
		var (
			rs = &reply.Reply{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Redis.DelAuditIndexs(context.Background(), rs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyAddAuditIndex(t *testing.T) {
	convey.Convey("AddAuditIndex", t, func(ctx convey.C) {
		var (
			rp = &reply.Reply{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Redis.AddAuditIndex(context.Background(), rp)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyAddFloorIndexEnd(t *testing.T) {
	convey.Convey("AddFloorIndexEnd", t, func(ctx convey.C) {
		var (
			oid = int64(0)
			tp  = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Redis.AddFloorIndexEnd(context.Background(), oid, tp)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyAddFloorIndex(t *testing.T) {
	convey.Convey("AddFloorIndex", t, func(ctx convey.C) {
		var (
			oid = int64(0)
			tp  = int8(0)
			rs  = &reply.Reply{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Redis.AddFloorIndex(context.Background(), oid, tp, rs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyAddCountIndexBatch(t *testing.T) {
	convey.Convey("AddCountIndexBatch", t, func(ctx convey.C) {
		var (
			oid = int64(0)
			tp  = int8(0)
			rs  = &reply.Reply{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Redis.AddCountIndexBatch(context.Background(), oid, tp, rs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyAddCountIndex(t *testing.T) {
	convey.Convey("AddCountIndex", t, func(ctx convey.C) {
		var (
			oid = int64(0)
			tp  = int8(0)
			rp  = &reply.Reply{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Redis.AddCountIndex(context.Background(), oid, tp, rp)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyAddLikeIndexBatch(t *testing.T) {
	convey.Convey("AddLikeIndexBatch", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			tp   = int8(0)
			rpts map[int64]*reply.Report
			rs   = &reply.Reply{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Redis.AddLikeIndexBatch(c, oid, tp, rpts, rs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyAddLikeIndex(t *testing.T) {
	convey.Convey("AddLikeIndex", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			tp   = int8(0)
			rpts map[int64]*reply.Report
			r    = &reply.Reply{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Redis.AddLikeIndex(c, oid, tp, rpts, r)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyAddIndex(t *testing.T) {
	convey.Convey("AddIndex", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			oid       = int64(0)
			tp        = int8(0)
			rpt       = &reply.Report{}
			rp        = &reply.Reply{}
			isRecover bool
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Redis.AddIndex(c, oid, tp, rpt, rp, isRecover)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyDelIndexBySortType(t *testing.T) {
	convey.Convey("DelIndexBySortType", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			rp       = &reply.Reply{}
			sortType = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Redis.DelIndexBySortType(c, rp, sortType)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyDelIndex(t *testing.T) {
	convey.Convey("DelIndex", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			rp = &reply.Reply{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Redis.DelIndex(c, rp)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyAddNewChildIndex(t *testing.T) {
	convey.Convey("AddNewChildIndex", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			root = int64(0)
			rs   = &reply.Reply{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Redis.AddNewChildIndex(c, root, rs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyAddTopOid(t *testing.T) {
	convey.Convey("AddTopOid", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			tp  = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Redis.AddTopOid(c, oid, tp)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyDelTopOid(t *testing.T) {
	convey.Convey("DelTopOid", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			tp  = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Redis.DelTopOid(c, oid, tp)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyAddLike(t *testing.T) {
	convey.Convey("AddLike", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			rpID = int64(0)
			ras  = &reply.Action{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Redis.AddLike(c, rpID, ras)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyDelLike(t *testing.T) {
	convey.Convey("DelLike", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			rpID = int64(0)
			ra   = &reply.Action{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Redis.DelLike(c, rpID, ra)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyExpireLike(t *testing.T) {
	convey.Convey("ExpireLike", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			rpID = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ok, err := d.Redis.ExpireLike(c, rpID)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyRange(t *testing.T) {
	convey.Convey("Range", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			oid   = int64(0)
			tp    = int8(0)
			sort  = int8(0)
			start = int(0)
			end   = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rpIds, err := d.Redis.Range(c, oid, tp, sort, start, end)
			ctx.Convey("Then err should be nil.rpIds should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rpIds, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyFloorEnd(t *testing.T) {
	convey.Convey("FloorEnd", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			tp  = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			score, found, err := d.Redis.FloorEnd(c, oid, tp)
			ctx.Convey("Then err should be nil.score,found should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(found, convey.ShouldNotBeNil)
				ctx.So(score, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyMinScore(t *testing.T) {
	convey.Convey("MinScore", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			tp   = int8(0)
			sort = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			score, err := d.Redis.MinScore(c, oid, tp, sort)
			ctx.Convey("Then err should be nil.score should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(score, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyCountReplies(t *testing.T) {
	convey.Convey("CountReplies", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			tp   = int8(0)
			sort = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			count, err := d.Redis.CountReplies(c, oid, tp, sort)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyExpireDialogIndex(t *testing.T) {
	convey.Convey("ExpireDialogIndex", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			dialogID = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ok, err := d.Redis.ExpireDialogIndex(c, dialogID)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyExpireIndex(t *testing.T) {
	convey.Convey("ExpireIndex", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			tp   = int8(0)
			sort = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ok, err := d.Redis.ExpireIndex(c, oid, tp, sort)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyExpireNewChildIndex(t *testing.T) {
	convey.Convey("ExpireNewChildIndex", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			root = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ok, err := d.Redis.ExpireNewChildIndex(c, root)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyAddDialogIndex(t *testing.T) {
	convey.Convey("AddDialogIndex", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			dialogID = int64(0)
			rps      = []*reply.Reply{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Redis.AddDialogIndex(c, dialogID, rps)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplySetUserReportCnt(t *testing.T) {
	convey.Convey("SetUserReportCnt", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(0)
			count = int(0)
			now   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Redis.SetUserReportCnt(c, mid, count, now)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyGetUserReportCnt(t *testing.T) {
	convey.Convey("GetUserReportCnt", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			now = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			count, err := d.Redis.GetUserReportCnt(c, mid, now)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyGetUserReportTTL(t *testing.T) {
	convey.Convey("GetUserReportTTL", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			now = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ttl, err := d.Redis.GetUserReportTTL(c, mid, now)
			ctx.Convey("Then err should be nil.ttl should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ttl, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyRankIndex(t *testing.T) {
	convey.Convey("RankIndex", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			tp   = int8(0)
			rpID = int64(0)
			sort = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rank, err := d.Redis.RankIndex(c, oid, tp, rpID, sort)
			ctx.Convey("Then err should be nil.rank should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rank, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyExpireUserAct(t *testing.T) {
	convey.Convey("ExpireUserAct", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ok, err := d.Redis.ExpireUserAct(c, mid)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyAddUserActs(t *testing.T) {
	convey.Convey("AddUserActs", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			mid     = int64(0)
			actions map[int64]int8
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Redis.AddUserActs(c, mid, actions)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyDelUserAct(t *testing.T) {
	convey.Convey("DelUserAct", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(0)
			rpID = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Redis.DelUserAct(c, mid, rpID)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyUserAct(t *testing.T) {
	convey.Convey("UserAct", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(0)
			rpID = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			act, err := d.Redis.UserAct(c, mid, rpID)
			ctx.Convey("Then err should be nil.act should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(act, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyUserActs(t *testing.T) {
	convey.Convey("UserActs", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(0)
			rpids = []int64{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			acts, err := d.Redis.UserActs(c, mid, rpids)
			ctx.Convey("Then err should be nil.acts should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(acts, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplySpamReply(t *testing.T) {
	convey.Convey("SpamReply", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rec, daily, err := d.Redis.SpamReply(c, mid)
			ctx.Convey("Then err should be nil.rec,daily should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(daily, convey.ShouldNotBeNil)
				ctx.So(rec, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplySpamAction(t *testing.T) {
	convey.Convey("SpamAction", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			code, err := d.Redis.SpamAction(c, mid)
			ctx.Convey("Then err should be nil.code should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(code, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyNotifyCnt(t *testing.T) {
	convey.Convey("NotifyCnt", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			typ = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			cnt, err := d.Redis.NotifyCnt(c, oid, typ)
			ctx.Convey("Then err should be nil.cnt should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(cnt, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplySetNotifyCnt(t *testing.T) {
	convey.Convey("SetNotifyCnt", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			typ = int8(0)
			cnt = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Redis.SetNotifyCnt(c, oid, typ, cnt)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyMaxLikeCnt(t *testing.T) {
	convey.Convey("MaxLikeCnt", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			rpid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			cnt, err := d.Redis.MaxLikeCnt(c, rpid)
			ctx.Convey("Then err should be nil.cnt should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(cnt, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplySetMaxLikeCnt(t *testing.T) {
	convey.Convey("SetMaxLikeCnt", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			rpid = int64(0)
			cnt  = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Redis.SetMaxLikeCnt(c, rpid, cnt)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
