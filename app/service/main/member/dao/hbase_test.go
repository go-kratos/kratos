package dao

//import (
//	"context"
//	"testing"
//
//	"github.com/smartystreets/goconvey/convey"
//)
//
//func TestDaomoralRowKey(t *testing.T) {
//	var (
//		mid = int64(0)
//		ts  = int64(0)
//		tid uint64
//	)
//	convey.Convey("moralRowKey", t, func(ctx convey.C) {
//		p1 := moralRowKey(mid, ts, tid)
//		ctx.Convey("p1 should not be nil", func(ctx convey.C) {
//			ctx.So(p1, convey.ShouldNotBeNil)
//		})
//	})
//}
//
//func TestDaoAddMoralLog(t *testing.T) {
//	var (
//		c       = context.Background()
//		mid     = int64(0)
//		ts      = int64(0)
//		s       = []byte("abcd")
//		content = map[string][]byte{"ts": s}
//	)
//	convey.Convey("AddMoralLog", t, func(ctx convey.C) {
//		err := d.AddMoralLog(c, mid, ts, content)
//		ctx.Convey("Error should be nil", func(ctx convey.C) {
//			ctx.So(err, convey.ShouldBeNil)
//		})
//	})
//}
//
//func TestDaogenID(t *testing.T) {
//	convey.Convey("genID", t, func(ctx convey.C) {
//		p1 := genID()
//		ctx.Convey("p1 should not be nil", func(ctx convey.C) {
//			ctx.So(p1, convey.ShouldNotBeNil)
//		})
//	})
//}
