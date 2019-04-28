package dao

import (
	"context"
	"go-common/app/interface/main/dm/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddProtectApply(t *testing.T) {
	convey.Convey("AddProtectApply", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			pas = []*model.Pa{
				{
					ID:  123,
					CID: 1,
					UID: 1,
				},
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			affect, err := testDao.AddProtectApply(c, pas)
			convCtx.Convey("Then err should be nil.affect should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(affect, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoProtectApplyTime(t *testing.T) {
	convey.Convey("ProtectApplyTime", t, func(convCtx convey.C) {
		var (
			c    = context.Background()
			dmid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			no, err := testDao.ProtectApplyTime(c, dmid)
			convCtx.Convey("Then err should be nil.no should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(no, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoProtectApplies(t *testing.T) {
	convey.Convey("ProtectApplies", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			uid   = int64(123)
			aid   = int64(111)
			order = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			_, err := testDao.ProtectApplies(c, uid, aid, order)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoProtectAids(t *testing.T) {
	convey.Convey("ProtectAids", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			uid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := testDao.ProtectAids(c, uid)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUptPaStatus(t *testing.T) {
	convey.Convey("UptPaStatus", t, func(convCtx convey.C) {
		var (
			c      = context.Background()
			uid    = int64(123)
			ids    = "123"
			status = int(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			affect, err := testDao.UptPaStatus(c, uid, ids, status)
			convCtx.Convey("Then err should be nil.affect should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(affect, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoProtectApplyByIDs(t *testing.T) {
	convey.Convey("ProtectApplyByIDs", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			uid = int64(123)
			ids = "123"
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := testDao.ProtectApplyByIDs(c, uid, ids)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUptPaNoticeSwitch(t *testing.T) {
	convey.Convey("UptPaNoticeSwitch", t, func(convCtx convey.C) {
		var (
			c      = context.Background()
			uid    = int64(0)
			status = int(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			affect, err := testDao.UptPaNoticeSwitch(c, uid, status)
			convCtx.Convey("Then err should be nil.affect should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(affect, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoPaNoticeClose(t *testing.T) {
	convey.Convey("PaNoticeClose", t, func(convCtx convey.C) {
		var (
			c    = context.Background()
			uids = []int64{123}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := testDao.PaNoticeClose(c, uids)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoProtectApplyStatistics(t *testing.T) {
	convey.Convey("ProtectApplyStatistics", t, func(convCtx convey.C) {
		var (
			c = context.Background()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := testDao.ProtectApplyStatistics(c)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaotwoDayAgo22(t *testing.T) {
	convey.Convey("twoDayAgo22", t, func(convCtx convey.C) {
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := twoDayAgo22()
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoPaUsrStat(t *testing.T) {
	convey.Convey("PaUsrStat", t, func(convCtx convey.C) {
		var (
			c = context.Background()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := testDao.PaUsrStat(c)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
