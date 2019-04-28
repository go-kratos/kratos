package like

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

func TestLikeSubject(t *testing.T) {
	convey.Convey("Subject", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sid = int64(10193)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			n, err := d.Subject(c, sid)
			ctx.Convey("Then err should be nil.n should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(n, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeInOnlinelog(t *testing.T) {
	convey.Convey("InOnlinelog", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			sid   = int64(10193)
			aid   = int64(1)
			stage = int64(1)
			yes   = int64(1)
			no    = int64(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.InOnlinelog(c, sid, aid, stage, yes, no)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeSubjectList(t *testing.T) {
	convey.Convey("SubjectList", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			types = []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 14, 13, 15, 16, 17, 18, 19}
			ts    = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.SubjectList(c, types, ts)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeSubjectTotalStat(t *testing.T) {
	convey.Convey("SubjectTotalStat", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sid = int64(10338)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rs, err := d.SubjectTotalStat(c, sid)
			ctx.Convey("Then err should be nil.rs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeAddLotteryTimes(t *testing.T) {
	convey.Convey("AddLotteryTimes", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sid = int64(315)
			mid = int64(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("GET", d.addLotteryTimesURL).Reply(200).SetHeaders(map[string]string{
				"Code": "0",
			})
			err := d.AddLotteryTimes(c, sid, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
