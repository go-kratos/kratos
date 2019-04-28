package dao

import (
	"context"
	"go-common/app/admin/main/dm/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaohitSubtitle(t *testing.T) {
	convey.Convey("hitSubtitle", t, func(ctx convey.C) {
		var (
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := testDao.hitSubtitle(oid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetSubtitles(t *testing.T) {
	convey.Convey("GetSubtitles", t, func(ctx convey.C) {
		var (
			c           = context.Background()
			oid         = int64(0)
			subtitleIds = []int64{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.GetSubtitles(c, oid, subtitleIds)
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
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.GetSubtitle(c, oid, subtitleID)
		})
	})
}

func TestDaoUpdateSubtitle(t *testing.T) {
	convey.Convey("UpdateSubtitle", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			subtitle = &model.Subtitle{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.UpdateSubtitle(c, subtitle)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCountSubtitleDraft(t *testing.T) {
	convey.Convey("CountSubtitleDraft", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			mid = int64(0)
			lan = uint8(0)
			tp  = int32(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			count, err := testDao.CountSubtitleDraft(c, oid, mid, lan, tp)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxUpdateSubtitle(t *testing.T) {
	convey.Convey("TxUpdateSubtitle", t, func(ctx convey.C) {
		var (
			tx, _    = testDao.BeginBiliDMTrans(context.TODO())
			subtitle = &model.Subtitle{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.TxUpdateSubtitle(tx, subtitle)
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTxGetSubtitleID(t *testing.T) {
	convey.Convey("TxGetSubtitleID", t, func(ctx convey.C) {
		var (
			tx, _ = testDao.BeginBiliDMTrans(context.TODO())
			oid   = int64(0)
			tp    = int32(0)
			lan   = uint8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.TxGetSubtitleID(tx, oid, tp, lan)
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTxUpdateSubtitlePub(t *testing.T) {
	convey.Convey("TxUpdateSubtitlePub", t, func(ctx convey.C) {
		var (
			tx, _       = testDao.BeginBiliDMTrans(context.TODO())
			subtitlePub = &model.SubtitlePub{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.TxUpdateSubtitlePub(tx, subtitlePub)
		})
		ctx.Reset(func() {
			tx.Commit()
		})

	})
}

func TestDaoSubtitleLans(t *testing.T) {
	convey.Convey("SubtitleLans", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			subtitleLans, err := testDao.SubtitleLans(c)
			ctx.Convey("Then err should be nil.subtitleLans should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(subtitleLans, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddSubtitleSubject(t *testing.T) {
	convey.Convey("AddSubtitleSubject", t, func(ctx convey.C) {
		var (
			c               = context.Background()
			subtitleSubject = &model.SubtitleSubject{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.AddSubtitleSubject(c, subtitleSubject)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoGetSubtitleSubject(t *testing.T) {
	convey.Convey("GetSubtitleSubject", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			aid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			subtitleSubject, err := testDao.GetSubtitleSubject(c, aid)
			ctx.Convey("Then err should be nil.subtitleSubject should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(subtitleSubject, convey.ShouldNotBeNil)
			})
		})
	})
}
