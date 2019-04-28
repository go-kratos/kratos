package dao

import (
	"context"
	"go-common/app/job/main/dm2/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaohitSubject(t *testing.T) {
	convey.Convey("hitSubject", t, func(ctx convey.C) {
		var (
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := testDao.hitSubject(oid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaohitIndex(t *testing.T) {
	convey.Convey("hitIndex", t, func(ctx convey.C) {
		var (
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := testDao.hitIndex(oid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaohitContent(t *testing.T) {
	convey.Convey("hitContent", t, func(ctx convey.C) {
		var (
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := testDao.hitContent(oid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaohitSubtile(t *testing.T) {
	convey.Convey("hitSubtile", t, func(ctx convey.C) {
		var (
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := testDao.hitSubtile(oid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddSubject(t *testing.T) {
	convey.Convey("AddSubject", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			tp       = int32(0)
			oid      = int64(0)
			pid      = int64(0)
			mid      = int64(0)
			maxlimit = int64(0)
			attr     = int32(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			testDao.AddSubject(c, tp, oid, pid, mid, maxlimit, attr)
		})
	})
}

func TestDaoUpdateSubAttr(t *testing.T) {
	convey.Convey("UpdateSubAttr", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			tp   = int32(0)
			oid  = int64(0)
			attr = int32(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			affect, err := testDao.UpdateSubAttr(c, tp, oid, attr)
			ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affect, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateSubMid(t *testing.T) {
	convey.Convey("UpdateSubMid", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tp  = int32(0)
			oid = int64(0)
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			affect, err := testDao.UpdateSubMid(c, tp, oid, mid)
			ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affect, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSubject(t *testing.T) {
	convey.Convey("Subject", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tp  = int32(0)
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			s, err := testDao.Subject(c, tp, oid)
			ctx.Convey("Then err should be nil.s should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(s, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateChildpool(t *testing.T) {
	convey.Convey("UpdateChildpool", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			tp        = int32(0)
			oid       = int64(0)
			childpool = int32(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			affect, err := testDao.UpdateChildpool(c, tp, oid, childpool)
			ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affect, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxIncrSubjectCount(t *testing.T) {
	convey.Convey("TxIncrSubjectCount", t, func(ctx convey.C) {
		var (
			tx, _     = testDao.BeginTran(c)
			tp        = int32(0)
			oid       = int64(0)
			acount    = int64(0)
			count     = int64(0)
			childpool = int32(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			affect, err := testDao.TxIncrSubjectCount(tx, tp, oid, acount, count, childpool)
			ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affect, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTxAddIndex(t *testing.T) {
	convey.Convey("TxAddIndex", t, func(ctx convey.C) {
		var (
			tx, _ = testDao.BeginTran(c)
			m     = &model.DM{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			id, err := testDao.TxAddIndex(tx, m)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoIndexs(t *testing.T) {
	convey.Convey("Indexs", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tp  = int32(0)
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			testDao.Indexs(c, tp, oid)
		})
	})
}

func TestDaoIndexsSeg(t *testing.T) {
	convey.Convey("IndexsSeg", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tp    = int32(0)
			oid   = int64(0)
			ps    = int64(0)
			pe    = int64(0)
			limit = int64(0)
			pool  = int32(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			testDao.IndexsSeg(c, tp, oid, ps, pe, limit, pool)
		})
	})
}

func TestDaoIndexsSegID(t *testing.T) {
	convey.Convey("IndexsSegID", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tp    = int32(0)
			oid   = int64(0)
			ps    = int64(0)
			pe    = int64(0)
			limit = int64(0)
			pool  = int32(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			testDao.IndexsSegID(c, tp, oid, ps, pe, limit, pool)
		})
	})
}

func TestDaoIndexsID(t *testing.T) {
	convey.Convey("IndexsID", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			tp   = int32(0)
			oid  = int64(0)
			pool = int32(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			dmids, err := testDao.IndexsID(c, tp, oid, pool)
			ctx.Convey("Then err should be nil.dmids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(dmids, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoIndexsByid(t *testing.T) {
	convey.Convey("IndexsByid", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tp    = int32(0)
			oid   = int64(0)
			dmids = []int64{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			testDao.IndexsByid(c, tp, oid, dmids)
		})
	})
}

func TestDaoIndexsByPool(t *testing.T) {
	convey.Convey("IndexsByPool", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			tp   = int32(0)
			oid  = int64(0)
			pool = int32(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			dms, dmids, err := testDao.IndexsByPool(c, tp, oid, pool)
			ctx.Convey("Then err should be nil.dms,dmids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(dmids, convey.ShouldNotBeNil)
				ctx.So(dms, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxAddContent(t *testing.T) {
	convey.Convey("TxAddContent", t, func(ctx convey.C) {
		var (
			tx, _ = testDao.BeginTran(c)
			oid   = int64(0)
			m     = &model.Content{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			id, err := testDao.TxAddContent(tx, oid, m)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTxAddContentSpecial(t *testing.T) {
	convey.Convey("TxAddContentSpecial", t, func(ctx convey.C) {
		var (
			tx, _ = testDao.BeginTran(c)
			m     = &model.ContentSpecial{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			id, err := testDao.TxAddContentSpecial(tx, m)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoContent(t *testing.T) {
	convey.Convey("Content", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			dmid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ct, err := testDao.Content(c, oid, dmid)
			ctx.Convey("Then err should be nil.ct should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ct, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoContents(t *testing.T) {
	convey.Convey("Contents", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			oid   = int64(0)
			dmids = []int64{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ctsMap, err := testDao.Contents(c, oid, dmids)
			ctx.Convey("Then err should be nil.ctsMap should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ctsMap, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoContentsSpecial(t *testing.T) {
	convey.Convey("ContentsSpecial", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			dmids = []int64{123}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			testDao.ContentsSpecial(c, dmids)
		})
	})
}

func TestDaoContentSpecial(t *testing.T) {
	convey.Convey("ContentSpecial", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			dmid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			contentSpe, err := testDao.ContentSpecial(c, dmid)
			ctx.Convey("Then err should be nil.contentSpe should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(contentSpe, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelDMHideState(t *testing.T) {
	convey.Convey("DelDMHideState", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			tp   = int32(0)
			oid  = int64(0)
			dmid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			affect, err := testDao.DelDMHideState(c, tp, oid, dmid)
			ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affect, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxIncrSubMCount(t *testing.T) {
	convey.Convey("TxIncrSubMCount", t, func(ctx convey.C) {
		var (
			tx, _ = testDao.BeginTran(c)
			tp    = int32(0)
			oid   = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			affect, err := testDao.TxIncrSubMCount(tx, tp, oid)
			ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affect, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoUpdateSubtitle(t *testing.T) {
	convey.Convey("UpdateSubtitle", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			subtitle = &model.Subtitle{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := testDao.UpdateSubtitle(c, subtitle)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoGetSubtitles(t *testing.T) {
	convey.Convey("GetSubtitles", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tp  = int32(0)
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			subtitles, err := testDao.GetSubtitles(c, tp, oid)
			ctx.Convey("Then err should be nil.subtitles should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(subtitles, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetSubtitle(t *testing.T) {
	convey.Convey("GetSubtitle", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			oid        = int64(0)
			subtitleID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			testDao.GetSubtitle(c, oid, subtitleID)
		})
	})
}

func TestDaoTxUpdateSubtitle(t *testing.T) {
	convey.Convey("TxUpdateSubtitle", t, func(ctx convey.C) {
		var (
			tx, _    = testDao.BeginTran(c)
			subtitle = &model.Subtitle{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			testDao.TxUpdateSubtitle(tx, subtitle)
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTxAddSubtitlePub(t *testing.T) {
	convey.Convey("TxAddSubtitlePub", t, func(ctx convey.C) {
		var (
			tx, _       = testDao.BeginTran(c)
			subtitlePub = &model.SubtitlePub{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			testDao.TxAddSubtitlePub(tx, subtitlePub)
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoMaskMids(t *testing.T) {
	convey.Convey("MaskMids", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			mids, err := testDao.MaskMids(c)
			ctx.Convey("Then err should be nil.mids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(mids, convey.ShouldNotBeNil)
			})
		})
	})
}
