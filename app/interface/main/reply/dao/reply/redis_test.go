package reply

import (
	"context"
	"testing"
	"time"

	"go-common/app/interface/main/reply/conf"
	"go-common/app/interface/main/reply/model/reply"

	"github.com/smartystreets/goconvey/convey"
)

func TestReplyNewRedisDao(t *testing.T) {
	convey.Convey("NewRedisDao", t, func(ctx convey.C) {
		var (
			c = conf.Conf.Redis
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := NewRedisDao(c)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplykeyDialogIdx(t *testing.T) {
	convey.Convey("keyDialogIdx", t, func(ctx convey.C) {
		var (
			dialogID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
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
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyIdx(oid, tp, sort)
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
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyAuditIdx(oid, tp)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplykeyRtIdx(t *testing.T) {
	convey.Convey("keyRtIdx", t, func(ctx convey.C) {
		var (
			rpID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyRtIdx(rpID)
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
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
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
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyLike(rpID)
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
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
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
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
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
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
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
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyTopOid(tp)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyPing(t *testing.T) {
	convey.Convey("Ping", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Redis.Ping(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyAddFloorIndex(t *testing.T) {
	convey.Convey("AddFloorIndex", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			tp  = int8(0)
			rs  = &reply.Reply{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Redis.AddFloorIndex(c, oid, tp, rs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyAddCountIndex(t *testing.T) {
	convey.Convey("AddCountIndex", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			tp  = int8(0)
			rs  = &reply.Reply{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Redis.AddCountIndex(c, oid, tp, rs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyAddLikeIndex(t *testing.T) {
	convey.Convey("AddLikeIndex", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			tp  = int8(0)
			r   = &reply.Reply{}
			rpt = &reply.Report{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Redis.AddLikeIndex(c, oid, tp, r, rpt)
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
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Redis.DelIndex(c, rp)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyAddFloorIndexByRoot(t *testing.T) {
	convey.Convey("AddFloorIndexByRoot", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			root = int64(0)
			rs   = &reply.Reply{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Redis.AddFloorIndexByRoot(c, root, rs)
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
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
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
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
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
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
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
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rpIds, isEnd, err := d.Redis.Range(c, oid, tp, sort, start, end)
			ctx.Convey("Then err should be nil.rpIds,isEnd should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(isEnd, convey.ShouldNotBeNil)
				if len(rpIds) <= 0 {
					ctx.So(rpIds, convey.ShouldBeEmpty)
				}
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
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.Redis.CountReplies(c, oid, tp, sort)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyUserAuditReplies(t *testing.T) {
	convey.Convey("UserAuditReplies", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
			oid = int64(1)
			tp  = int8(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rpIds, err := d.Redis.UserAuditReplies(c, mid, oid, tp)
			ctx.Convey("Then err should be nil.rpIds should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rpIds, convey.ShouldBeEmpty)
			})
		})
	})
}

func TestReplyRangeByRoot(t *testing.T) {
	convey.Convey("RangeByRoot", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			root  = int64(0)
			start = int(0)
			end   = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rpIds, err := d.Redis.RangeByRoot(c, root, start, end)
			ctx.Convey("Then err should be nil.rpIds should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rpIds, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyRangeByRoots(t *testing.T) {
	convey.Convey("RangeByRoots", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			roots = []int64{}
			start = int(0)
			end   = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			mrpids, idx, miss, err := d.Redis.RangeByRoots(c, roots, start, end)
			ctx.Convey("Then err should be nil.mrpids,idx,miss should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(miss, convey.ShouldBeNil)
				ctx.So(idx, convey.ShouldBeNil)
				ctx.So(mrpids, convey.ShouldNotBeNil)
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
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ok, err := d.Redis.ExpireIndex(c, oid, tp, sort)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyExpireIndexByRoot(t *testing.T) {
	convey.Convey("ExpireIndexByRoot", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			root = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ok, err := d.Redis.ExpireIndexByRoot(c, root)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
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
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
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
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
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
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
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
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rank, err := d.Redis.RankIndex(c, oid, tp, rpID, sort)
			ctx.Convey("Then err should be nil.rank should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rank, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyRankIndexByRoot(t *testing.T) {
	convey.Convey("RankIndexByRoot", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			root = int64(0)
			rpID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rank, err := d.Redis.RankIndexByRoot(c, root, rpID)
			ctx.Convey("Then err should be nil.rank should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rank, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyOidHaveTop(t *testing.T) {
	convey.Convey("OidHaveTop", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			tp  = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ok, err := d.Redis.OidHaveTop(c, oid, tp)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
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
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			recent, daily, err := d.Redis.SpamReply(c, mid)
			ctx.Convey("Then err should be nil.recent,daily should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(daily, convey.ShouldNotBeNil)
				ctx.So(recent, convey.ShouldNotBeNil)
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
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			code, err := d.Redis.SpamAction(c, mid)
			ctx.Convey("Then err should be nil.code should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(code, convey.ShouldNotBeNil)
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
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ok, err := d.Redis.ExpireDialogIndex(c, dialogID)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyRangeRpsByDialog(t *testing.T) {
	convey.Convey("RangeRpsByDialog", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			dialog = int64(0)
			start  = int(0)
			end    = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rpIDs, err := d.Redis.RangeRpsByDialog(c, dialog, start, end)
			ctx.Convey("Then err should be nil.rpIDs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rpIDs, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyDialogDesc(t *testing.T) {
	convey.Convey("DialogDesc", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			dialog = int64(0)
			floor  = int(0)
			size   = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rpIDs, err := d.Redis.DialogDesc(c, dialog, floor, size)
			ctx.Convey("Then err should be nil.rpIDs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rpIDs, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyDialogMinMaxFloor(t *testing.T) {
	convey.Convey("DialogMinMaxFloor", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			dialog = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			minFloor, maxFloor, err := d.Redis.DialogMinMaxFloor(c, dialog)
			ctx.Convey("Then err should be nil.minFloor,maxFloor should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(maxFloor, convey.ShouldBeZeroValue)
				ctx.So(minFloor, convey.ShouldBeZeroValue)
			})
		})
	})
}

func TestReplyDialogByCursor(t *testing.T) {
	convey.Convey("DialogByCursor", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			dialog = int64(0)
			cursor = &reply.Cursor{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rpIDs, err := d.Redis.DialogByCursor(c, dialog, cursor)
			ctx.Convey("Then err should be nil.rpIDs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rpIDs, convey.ShouldBeNil)
			})
		})
	})
}

func TestDelReplyIncr(t *testing.T) {
	convey.Convey("test DelReplyIncr", t, func(ctx convey.C) {
		var (
			mid = int64(2233)
			c   = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Redis.DelReplyIncr(c, mid, true)
			ctx.Convey("d.Redis.DelReplyIncr up error should  be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
			err = d.Redis.DelReplyIncr(c, mid, false)
			ctx.Convey("d.Redis.DelReplyIncr error should  be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
			err = d.Redis.DelReplySpam(c, mid)
			ctx.Convey("d.Redis.DelReplySpam error should  be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
