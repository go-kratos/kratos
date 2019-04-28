package dao

import (
	"context"
	"go-common/app/interface/main/dm2/model"
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

func TestDaoGetSubtitleIds(t *testing.T) {
	convey.Convey("GetSubtitleIds", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			tp  = int32(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			subtitlIds, err := testDao.GetSubtitleIds(c, oid, tp)
			ctx.Convey("Then err should be nil.subtitlIds should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(subtitlIds, convey.ShouldNotBeNil)
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

func TestDaoGetSubtitleDraft(t *testing.T) {
	convey.Convey("GetSubtitleDraft", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			tp  = int32(0)
			mid = int64(0)
			lan = uint8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			subtitle, err := testDao.GetSubtitleDraft(c, oid, tp, mid, lan)
			ctx.Convey("Then err should be nil.subtitle should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(subtitle, convey.ShouldNotBeNil)
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
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.GetSubtitle(c, oid, subtitleID)
		})
	})
}

func TestDaoAddSubtitle(t *testing.T) {
	convey.Convey("AddSubtitle", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			subtitle = &model.Subtitle{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.AddSubtitle(c, subtitle)
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

func TestDaoTxUpdateSubtitle(t *testing.T) {
	convey.Convey("TxUpdateSubtitle", t, func(ctx convey.C) {
		var (
			tx, _    = testDao.BeginBiliDMTrans(c)
			subtitle = &model.Subtitle{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.TxUpdateSubtitle(tx, subtitle)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTxAddSubtitlePub(t *testing.T) {
	convey.Convey("TxAddSubtitlePub", t, func(ctx convey.C) {
		var (
			c           = context.Background()
			tx, _       = testDao.BeginBiliDMTrans(c)
			subtitlePub = &model.SubtitlePub{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.TxAddSubtitlePub(tx, subtitlePub)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTxGetSubtitleOne(t *testing.T) {
	convey.Convey("TxGetSubtitleOne", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = testDao.BeginBiliDMTrans(c)
			oid   = int64(0)
			tp    = int32(0)
			lan   = uint8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.TxGetSubtitleOne(tx, oid, tp, lan)
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

func TestDaoSubtitleLanAdd(t *testing.T) {
	convey.Convey("SubtitleLanAdd", t, func(ctx convey.C) {
		var (
			c           = context.Background()
			subtitleLan = &model.SubtitleLan{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.SubtitleLanAdd(c, subtitleLan)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpsertWaveFrom(t *testing.T) {
	convey.Convey("UpsertWaveFrom", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			waveForm = &model.WaveForm{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.UpsertWaveFrom(c, waveForm)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoGetWaveForm(t *testing.T) {
	convey.Convey("GetWaveForm", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			tp  = int32(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			waveForm, err := testDao.GetWaveForm(c, oid, tp)
			ctx.Convey("Then err should be nil.waveForm should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(waveForm, convey.ShouldNotBeNil)
			})
		})
	})
}
