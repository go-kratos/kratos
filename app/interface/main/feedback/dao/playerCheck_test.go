package dao

import (
	"context"
	xsql "database/sql"
	"fmt"
	"reflect"
	"testing"

	"go-common/library/database/sql"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

func TestDaoInPlayCheck(t *testing.T) {
	convey.Convey("InPlayCheck", t, func(ctx convey.C) {
		var (
			c             = context.TODO()
			platform      = int(0)
			isp           = int(0)
			ipChangeTimes = int(0)
			mid           = int64(0)
			checkTime     = int64(0)
			aid           = int64(0)
			connectSpeed  = int64(0)
			ioSpeed       = int64(0)
			region        = ""
			school        = ""
			ip            = ""
			cdn           = ""
		)
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			rows, err := d.InPlayCheck(c, platform, isp, ipChangeTimes, mid, checkTime, aid, connectSpeed, ioSpeed, region, school, ip, cdn)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When d.dbMs.Query gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.dbMs), "Exec",
				func(_ *sql.DB, _ context.Context, _ string, _ ...interface{}) (xsql.Result, error) {
					return nil, fmt.Errorf("d.dbMs.Exec Error")
				})
			defer guard.Unpatch()
			_, err := d.InPlayCheck(c, platform, isp, ipChangeTimes, mid, checkTime, aid, connectSpeed, ioSpeed, region, school, ip, cdn)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
