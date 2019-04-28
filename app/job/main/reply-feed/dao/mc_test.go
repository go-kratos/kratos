package dao

import (
	"context"
	"testing"

	"go-common/app/job/main/reply-feed/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoPingMc(t *testing.T) {
	convey.Convey("PingMc", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.PingMc(context.Background())
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaokeyReplyStat(t *testing.T) {
	convey.Convey("keyReplyStat", t, func(ctx convey.C) {
		var (
			rpID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyReplyStat(rpID)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSetReplyStatMc(t *testing.T) {
	convey.Convey("SetReplyStatMc", t, func(ctx convey.C) {
		var (
			rs = &model.ReplyStat{RpID: 0}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetReplyStatMc(context.Background(), rs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoReplyStatsMc(t *testing.T) {
	convey.Convey("ReplyStatsMc", t, func(ctx convey.C) {
		var (
			rpIDs       = []int64{}
			hitRpIDs    = []int64{}
			missedRpIDs = []int64{}
			c           = context.Background()
		)
		for i := 0; i < 5000; i++ {
			d.RemReplyStatMc(c, int64(i))
		}
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			for i := 0; i < 5000; i++ {
				if i%2 == 0 {
					err := d.SetReplyStatMc(c, &model.ReplyStat{RpID: int64(i)})
					if err != nil {
						t.Fatal(err)
					}
					hitRpIDs = append(hitRpIDs, int64(i))
				} else {
					missedRpIDs = append(missedRpIDs, int64(i))
				}
				rpIDs = append(rpIDs, int64(i))
			}
			rsMap, missIDs, err := d.ReplyStatsMc(c, rpIDs)
			ctx.Convey("Then err should be nil.rsMap,missIDs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(len(rsMap), convey.ShouldEqual, len(hitRpIDs))
				ctx.So(len(missIDs), convey.ShouldEqual, len(missedRpIDs))
			})
			for _, rpID := range hitRpIDs {
				if err = d.RemReplyStatMc(c, rpID); err != nil {
					t.Fatal(err)
				}
			}
		})
	})
}
