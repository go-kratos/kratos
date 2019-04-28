package dao

import (
	"context"
	"go-common/app/interface/main/dm2/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyMsgPubLock(t *testing.T) {
	convey.Convey("keyMsgPubLock", t, func(ctx convey.C) {
		var (
			mid      = int64(0)
			color    = int64(0)
			rnd      = int64(0)
			mode     = int32(0)
			fontsize = int32(0)
			ip       = ""
			msg      = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyMsgPubLock(mid, color, rnd, mode, fontsize, ip, msg)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyOidPubLock(t *testing.T) {
	convey.Convey("keyOidPubLock", t, func(ctx convey.C) {
		var (
			mid = int64(0)
			oid = int64(0)
			ip  = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyOidPubLock(mid, oid, ip)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyPubCntLock(t *testing.T) {
	convey.Convey("keyPubCntLock", t, func(ctx convey.C) {
		var (
			mid      = int64(0)
			color    = int64(0)
			mode     = int32(0)
			fontsize = int32(0)
			ip       = ""
			msg      = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyPubCntLock(mid, color, mode, fontsize, ip, msg)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyCharPubLock(t *testing.T) {
	convey.Convey("keyCharPubLock", t, func(ctx convey.C) {
		var (
			mid = int64(0)
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyCharPubLock(mid, oid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyXML(t *testing.T) {
	convey.Convey("keyXML", t, func(ctx convey.C) {
		var (
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyXML(oid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeySubject(t *testing.T) {
	convey.Convey("keySubject", t, func(ctx convey.C) {
		var (
			tp  = int32(0)
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keySubject(tp, oid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyAjax(t *testing.T) {
	convey.Convey("keyAjax", t, func(ctx convey.C) {
		var (
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyAjax(oid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyJudge(t *testing.T) {
	convey.Convey("keyJudge", t, func(ctx convey.C) {
		var (
			tp   = int8(0)
			oid  = int64(0)
			dmid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyJudge(tp, oid, dmid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyDMLimitMid(t *testing.T) {
	convey.Convey("keyDMLimitMid", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyDMLimitMid(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyAdvanceCmt(t *testing.T) {
	convey.Convey("keyAdvanceCmt", t, func(ctx convey.C) {
		var (
			mid  = int64(0)
			oid  = int64(0)
			mode = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyAdvanceCmt(mid, oid, mode)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyAdvLock(t *testing.T) {
	convey.Convey("keyAdvLock", t, func(ctx convey.C) {
		var (
			mid = int64(0)
			cid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyAdvLock(mid, cid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyHistory(t *testing.T) {
	convey.Convey("keyHistory", t, func(ctx convey.C) {
		var (
			tp        = int32(0)
			oid       = int64(0)
			timestamp = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyHistory(tp, oid, timestamp)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyHistoryIdx(t *testing.T) {
	convey.Convey("keyHistoryIdx", t, func(ctx convey.C) {
		var (
			tp    = int32(0)
			oid   = int64(0)
			month = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyHistoryIdx(tp, oid, month)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyDMMask(t *testing.T) {
	convey.Convey("keyDMMask", t, func(ctx convey.C) {
		var (
			tp   = int32(0)
			oid  = int64(0)
			plat = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyDMMask(tp, oid, plat)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSubjectCache(t *testing.T) {
	convey.Convey("SubjectCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tp  = int32(0)
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.SubjectCache(c, tp, oid)
		})
	})
}

func TestDaoSubjectsCache(t *testing.T) {
	convey.Convey("SubjectsCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			tp   = int32(0)
			oids = []int64{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			cached, missed, err := testDao.SubjectsCache(c, tp, oids)
			ctx.Convey("Then err should be nil.cached,missed should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(missed, convey.ShouldNotBeNil)
				ctx.So(cached, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddSubjectCache(t *testing.T) {
	convey.Convey("AddSubjectCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sub = &model.Subject{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.AddSubjectCache(c, sub)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelXMLCache(t *testing.T) {
	convey.Convey("DelXMLCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.DelXMLCache(c, oid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddXMLCache(t *testing.T) {
	convey.Convey("AddXMLCache", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			oid   = int64(0)
			value = []byte("")
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.AddXMLCache(c, oid, value)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoXMLCache(t *testing.T) {
	convey.Convey("XMLCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			data, err := testDao.XMLCache(c, oid)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAjaxDMCache(t *testing.T) {
	convey.Convey("AjaxDMCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			msgs, err := testDao.AjaxDMCache(c, oid)
			ctx.Convey("Then err should be nil.msgs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(msgs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddAjaxDMCache(t *testing.T) {
	convey.Convey("AddAjaxDMCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			msgs = []string{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.AddAjaxDMCache(c, oid, msgs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetDMJudgeCache(t *testing.T) {
	convey.Convey("SetDMJudgeCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			tp   = int8(0)
			oid  = int64(0)
			dmid = int64(0)
			l    = &model.JudgeDMList{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.SetDMJudgeCache(c, tp, oid, dmid, l)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDMJudgeCache(t *testing.T) {
	convey.Convey("DMJudgeCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			tp   = int8(0)
			oid  = int64(0)
			dmid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			l, err := testDao.DMJudgeCache(c, tp, oid, dmid)
			ctx.Convey("Then err should be nil.l should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(l, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddMsgPubLock(t *testing.T) {
	convey.Convey("AddMsgPubLock", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(0)
			color    = int64(0)
			rnd      = int64(0)
			mode     = int32(0)
			fontsize = int32(0)
			ip       = ""
			msg      = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.AddMsgPubLock(c, mid, color, rnd, mode, fontsize, ip, msg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoMsgPublock(t *testing.T) {
	convey.Convey("MsgPublock", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(0)
			color    = int64(0)
			rnd      = int64(0)
			mode     = int32(0)
			fontsize = int32(0)
			ip       = ""
			msg      = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			cached, err := testDao.MsgPublock(c, mid, color, rnd, mode, fontsize, ip, msg)
			ctx.Convey("Then err should be nil.cached should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(cached, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddOidPubLock(t *testing.T) {
	convey.Convey("AddOidPubLock", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			oid = int64(0)
			ip  = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.AddOidPubLock(c, mid, oid, ip)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoOidPubLock(t *testing.T) {
	convey.Convey("OidPubLock", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			oid = int64(0)
			ip  = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			cached, err := testDao.OidPubLock(c, mid, oid, ip)
			ctx.Convey("Then err should be nil.cached should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(cached, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddDMLimitCache(t *testing.T) {
	convey.Convey("AddDMLimitCache", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			mid     = int64(0)
			limiter = &model.Limiter{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.AddDMLimitCache(c, mid, limiter)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDMLimitCache(t *testing.T) {
	convey.Convey("DMLimitCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			limiter, err := testDao.DMLimitCache(c, mid)
			ctx.Convey("Then err should be nil.limiter should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(limiter, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddAdvanceCmtCache(t *testing.T) {
	convey.Convey("AddAdvanceCmtCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			mid  = int64(0)
			mode = ""
			adv  = &model.AdvanceCmt{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.AddAdvanceCmtCache(c, oid, mid, mode, adv)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAdvanceCmtCache(t *testing.T) {
	convey.Convey("AdvanceCmtCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			mid  = int64(0)
			mode = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			adv, err := testDao.AdvanceCmtCache(c, oid, mid, mode)
			ctx.Convey("Then err should be nil.adv should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(adv, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddAdvanceLock(t *testing.T) {
	convey.Convey("AddAdvanceLock", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			cid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			succeed := testDao.AddAdvanceLock(c, mid, cid)
			ctx.Convey("Then succeed should not be nil.", func(ctx convey.C) {
				ctx.So(succeed, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelAdvanceLock(t *testing.T) {
	convey.Convey("DelAdvanceLock", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			cid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.DelAdvanceLock(c, mid, cid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelAdvCache(t *testing.T) {
	convey.Convey("DelAdvCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(0)
			cid  = int64(0)
			mode = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.DelAdvCache(c, mid, cid, mode)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddHistoryCache(t *testing.T) {
	convey.Convey("AddHistoryCache", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			tp        = int32(0)
			oid       = int64(0)
			timestamp = int64(0)
			value     = []byte("")
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.AddHistoryCache(c, tp, oid, timestamp, value)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoHistoryCache(t *testing.T) {
	convey.Convey("HistoryCache", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			tp        = int32(0)
			oid       = int64(0)
			timestamp = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			data, err := testDao.HistoryCache(c, tp, oid, timestamp)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddHisIdxCache(t *testing.T) {
	convey.Convey("AddHisIdxCache", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tp    = int32(0)
			oid   = int64(0)
			month = ""
			dates = []string{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.AddHisIdxCache(c, tp, oid, month, dates)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoHistoryIdxCache(t *testing.T) {
	convey.Convey("HistoryIdxCache", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tp    = int32(0)
			oid   = int64(0)
			month = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			dates, err := testDao.HistoryIdxCache(c, tp, oid, month)
			ctx.Convey("Then err should be nil.dates should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(dates, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDMMaskCache(t *testing.T) {
	convey.Convey("DMMaskCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			tp   = int32(0)
			oid  = int64(0)
			plat = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.DMMaskCache(c, tp, oid, plat)
		})
	})
}

func TestDaoAddMaskCache(t *testing.T) {
	convey.Convey("AddMaskCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			tp   = int32(0)
			mask = &model.Mask{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.AddMaskCache(c, tp, mask)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
