package show

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	xsql "go-common/library/database/sql"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

// Test_SideBar test dao side bar
func TestDaoSideBar(t *testing.T) {
	convey.Convey("SidebBar", t, func(ctx convey.C) {
		ctx.Convey("When everyting is correct", func(ctx convey.C) {
			_, _, err := d.SideBar(context.Background(), time.Now())
			ctx.Convey("Error should be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Convey("When db.Query gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Query", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (*xsql.Rows, error) {
				return nil, fmt.Errorf("db.Query error")
			})
			defer guard.Unpatch()
			_, _, err := d.SideBar(context.Background(), time.Now())
			ctx.Convey("Error should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
