package dao

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"go-common/library/queue/databus"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

func TestDaoDatabusAddLog(t *testing.T) {
	convey.Convey("DatabusAddLog", t, func(convCtx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(0)
			exp    = int64(0)
			toExp  = int64(0)
			ts     = int64(0)
			oper   = ""
			reason = ""
			ip     = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.plogDatabus), "Send", func(_ *databus.Databus, _ context.Context, _ string, _ interface{}) error {
				return nil
			})
			defer monkey.UnpatchAll()
			err := d.DatabusAddLog(c, mid, exp, toExp, ts, oper, reason, ip)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})

		convCtx.Convey("When everything goes negative", func(convCtx convey.C) {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.plogDatabus), "Send", func(_ *databus.Databus, _ context.Context, _ string, _ interface{}) error {
				return fmt.Errorf("Failed send data err")
			})
			defer monkey.UnpatchAll()
			err := d.DatabusAddLog(c, mid, exp, toExp, ts, oper, reason, ip)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
