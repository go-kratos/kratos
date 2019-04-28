package dao

import (
	"context"
	"testing"

	"go-common/app/job/main/reply-feed/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestGenSQL(t *testing.T) {
	convey.Convey("genSQL", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := genSQL()
			t.Log(p1)
		})
	})
}

func TestDaoreportHit(t *testing.T) {
	convey.Convey("reportHit", t, func(ctx convey.C) {
		var (
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := reportHit(oid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoreplyHit(t *testing.T) {
	convey.Convey("replyHit", t, func(ctx convey.C) {
		var (
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := replyHit(oid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaosubjectHit(t *testing.T) {
	convey.Convey("subjectHit", t, func(ctx convey.C) {
		var (
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := subjectHit(oid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaosplit(t *testing.T) {
	convey.Convey("split", t, func(ctx convey.C) {
		var (
			s = []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := split(s, 5)
			p2 := split(s, 10)
			p3 := split(s, 20)
			p4 := split(s, 3)
			p5 := split(s, 1)
			p6 := split(s, 2)
			p7 := split(s, 4)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(len(p1), convey.ShouldEqual, 2)
				ctx.So(len(p2), convey.ShouldEqual, 1)
				ctx.So(len(p3), convey.ShouldEqual, 1)
				ctx.So(len(p4), convey.ShouldEqual, 4)
				ctx.So(len(p5), convey.ShouldEqual, 10)
				ctx.So(len(p6), convey.ShouldEqual, 5)
				ctx.So(len(p7), convey.ShouldEqual, 3)
			})
		})
	})
}

func TestDaoSlotStats(t *testing.T) {
	convey.Convey("SlotStats", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ss, err := d.SlotStats(context.Background())
			ctx.Convey("Then err should be nil.ss should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ss, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRpIDs(t *testing.T) {
	convey.Convey("RpIDs", t, func(ctx convey.C) {
		var (
			oid = int64(0)
			tp  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rpIDs, err := d.RpIDs(context.Background(), oid, tp)
			ctx.Convey("Then err should be nil.rpIDs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rpIDs, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoReportStatsByID(t *testing.T) {
	convey.Convey("ReportStatsByID", t, func(ctx convey.C) {
		var (
			oid   = int64(0)
			rpIDs = []int64{-1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			reportMap, err := d.ReportStatsByID(context.Background(), oid, rpIDs)
			ctx.Convey("Then err should be nil.reportMap should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(reportMap, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoReplyLHRCStatsByID(t *testing.T) {
	convey.Convey("ReplyLHRCStatsByID", t, func(ctx convey.C) {
		var (
			oid   = int64(0)
			rpIDs = []int64{-1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			replyMap, err := d.ReplyLHRCStatsByID(context.Background(), oid, rpIDs)
			ctx.Convey("Then err should be nil.replyMap should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(replyMap, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSubjectStats(t *testing.T) {
	convey.Convey("SubjectStats", t, func(ctx convey.C) {
		var (
			oid = int64(0)
			tp  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ctime, err := d.SubjectStats(context.Background(), oid, tp)
			ctx.Convey("Then err should be nil.ctime should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ctime, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoReplyLHRCStats(t *testing.T) {
	convey.Convey("ReplyLHRCStats", t, func(ctx convey.C) {
		var (
			oid = int64(0)
			tp  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			replyMap, err := d.ReplyLHRCStats(context.Background(), oid, tp)
			ctx.Convey("Then err should be nil.replyMap should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(replyMap, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSlotsMapping(t *testing.T) {
	convey.Convey("SlotsMapping", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			slotsMap, err := d.SlotsMapping(context.Background())
			ctx.Convey("Then err should be nil.slotsMap should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(slotsMap, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpsertStatistics(t *testing.T) {
	convey.Convey("UpsertStatistics", t, func(ctx convey.C) {
		var (
			name = ""
			date = int(0)
			hour = int(0)
			s    = &model.StatisticsStat{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.UpsertStatistics(context.Background(), name, date, hour, s)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
