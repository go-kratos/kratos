package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoTryLockReviewNotify(t *testing.T) {
	convey.Convey("TryLockReviewNotify", t, func() {
		now := time.Now()
		no := time.Date(2018, 1, 7, 20, now.Minute(), now.Second(), 651387237, time.UTC)
		p1, err := d.TryLockReviewNotify(context.Background(), no)
		convey.So(err, convey.ShouldBeNil)
		convey.So(p1, convey.ShouldNotBeNil)
	})
}

//func TestDaoreviewAuditNotifyLockKey(t *testing.T) {
//	convey.Convey("reviewAuditNotifyLockKey", t, func(ctx convey.C) {
//		var (
//			no = time.Now()
//		)
//		ctx.Convey("When everything goes positive", func(ctx convey.C) {
//			p1 := reviewAuditNotifyLockKey(no)
//			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
//				ctx.So(p1, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}
//
//func TestDaopingRedis(t *testing.T) {
//	convey.Convey("pingRedis", t, func(ctx convey.C) {
//		var (
//			c = context.Background()
//		)
//		ctx.Convey("When everything goes positive", func(ctx convey.C) {
//			err := d.pingRedis(c)
//			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
//				ctx.So(err, convey.ShouldBeNil)
//			})
//		})
//	})
//}
