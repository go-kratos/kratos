package archive

import (
	"context"
	"fmt"
	xsql "go-common/library/database/sql"
	"reflect"
	"testing"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

func TestMusics(t *testing.T) {
	convey.Convey("Musics", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("uat db ok", func(ctx convey.C) {
			_, err := d.AllMusics(c)
			ctx.Convey("Then err should be nil.bizs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Convey("db error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Query", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (*xsql.Rows, error) {
				return nil, fmt.Errorf("db.Query error")
			})
			defer guard.Unpatch()
			_, err := d.AllMusics(c)
			ctx.Convey("Then err should be nil.bizs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
