package dao

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDaonotifyKey(t *testing.T) {
	convey.Convey("notifyKey", t, func(convCtx convey.C) {
		var (
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := notifyKey(mid)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

//func TestDaoNotifyPurgeCache(t *testing.T) {
//	convey.Convey("NotifyPurgeCache", t, func(convCtx convey.C) {
//		var (
//			c      = context.Background()
//			mid    = int64(0)
//			action = ""
//		)
//		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
//			err := d.NotifyPurgeCache(c, mid, action)
//			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
//				convCtx.So(err, convey.ShouldBeNil)
//			})
//		})
//	})
//}
