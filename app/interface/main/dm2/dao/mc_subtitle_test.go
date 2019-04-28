package dao

import (
	"context"
	"go-common/app/interface/main/dm2/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaosubtitleKey(t *testing.T) {
	convey.Convey("subtitleKey", t, func(ctx convey.C) {
		var (
			oid        = int64(0)
			subtitleID = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := testDao.subtitleKey(oid, subtitleID)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaosubtitleVideoKey(t *testing.T) {
	convey.Convey("subtitleVideoKey", t, func(ctx convey.C) {
		var (
			oid = int64(0)
			tp  = int32(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := testDao.subtitleVideoKey(oid, tp)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaosubtitleDraftKey(t *testing.T) {
	convey.Convey("subtitleDraftKey", t, func(ctx convey.C) {
		var (
			oid = int64(0)
			tp  = int32(0)
			mid = int64(0)
			lan = uint8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := testDao.subtitleDraftKey(oid, tp, mid, lan)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaosubtitleSubjectKey(t *testing.T) {
	convey.Convey("subtitleSubjectKey", t, func(ctx convey.C) {
		var (
			aid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := testDao.subtitleSubjectKey(aid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaosubtitleReportTagKey(t *testing.T) {
	convey.Convey("subtitleReportTagKey", t, func(ctx convey.C) {
		var (
			bid = int64(0)
			rid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := testDao.subtitleReportTagKey(bid, rid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSetVideoSubtitleCache(t *testing.T) {
	convey.Convey("SetVideoSubtitleCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			tp  = int32(0)
			res = &model.VideoSubtitleCache{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.SetVideoSubtitleCache(c, oid, tp, res)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoVideoSubtitleCache(t *testing.T) {
	convey.Convey("VideoSubtitleCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			tp  = int32(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := testDao.VideoSubtitleCache(c, oid, tp)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelVideoSubtitleCache(t *testing.T) {
	convey.Convey("DelVideoSubtitleCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			tp  = int32(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.DelVideoSubtitleCache(c, oid, tp)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSubtitleDraftCache(t *testing.T) {
	convey.Convey("SubtitleDraftCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			tp  = int32(0)
			mid = int64(0)
			lan = uint8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.SubtitleDraftCache(c, oid, tp, mid, lan)
		})
	})
}

func TestDaoSetSubtitleDraftCache(t *testing.T) {
	convey.Convey("SetSubtitleDraftCache", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			subtitle = &model.Subtitle{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.SetSubtitleDraftCache(c, subtitle)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelSubtitleDraftCache(t *testing.T) {
	convey.Convey("DelSubtitleDraftCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			tp  = int32(0)
			mid = int64(0)
			lan = uint8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.DelSubtitleDraftCache(c, oid, tp, mid, lan)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSubtitlesCache(t *testing.T) {
	convey.Convey("SubtitlesCache", t, func(ctx convey.C) {
		var (
			c           = context.Background()
			oid         = int64(0)
			subtitleIds = []int64{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, missed, err := testDao.SubtitlesCache(c, oid, subtitleIds)
			ctx.Convey("Then err should be nil.res,missed should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(missed, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSubtitleCache(t *testing.T) {
	convey.Convey("SubtitleCache", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			oid        = int64(0)
			subtitleID = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.SubtitleCache(c, oid, subtitleID)
		})
	})
}

func TestDaoSetSubtitleCache(t *testing.T) {
	convey.Convey("SetSubtitleCache", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			subtitle = &model.Subtitle{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.SetSubtitleCache(c, subtitle)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelSubtitleCache(t *testing.T) {
	convey.Convey("DelSubtitleCache", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			oid        = int64(0)
			subtitleID = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.DelSubtitleCache(c, oid, subtitleID)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetSubtitleSubjectCache(t *testing.T) {
	convey.Convey("SetSubtitleSubjectCache", t, func(ctx convey.C) {
		var (
			c               = context.Background()
			subtitleSubject = &model.SubtitleSubject{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.SetSubtitleSubjectCache(c, subtitleSubject)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSubtitleSubjectCache(t *testing.T) {
	convey.Convey("SubtitleSubjectCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			aid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			subtitleSubject, err := testDao.SubtitleSubjectCache(c, aid)
			ctx.Convey("Then err should be nil.subtitleSubject should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(subtitleSubject, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelSubtitleSubjectCache(t *testing.T) {
	convey.Convey("DelSubtitleSubjectCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			aid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.DelSubtitleSubjectCache(c, aid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSubtitleWorlFlowTagCache(t *testing.T) {
	convey.Convey("SubtitleWorlFlowTagCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			bid = int64(0)
			rid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.SubtitleWorlFlowTagCache(c, bid, rid)
		})
	})
}

func TestDaoSetSubtitleWorlFlowTagCache(t *testing.T) {
	convey.Convey("SetSubtitleWorlFlowTagCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			bid  = int64(0)
			rid  = int64(0)
			data = []*model.WorkFlowTag{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.SetSubtitleWorlFlowTagCache(c, bid, rid, data)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
