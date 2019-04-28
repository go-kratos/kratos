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

func TestDaoPosRecs(t *testing.T) {
	convey.Convey("PosRecs", t, func(ctx convey.C) {
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			res, err := d.PosRecs(context.Background(), time.Now())
			ctx.Convey("Error should be nil, res should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When db.Query gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Query", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (*xsql.Rows, error) {
				return nil, fmt.Errorf("db.Query error")
			})
			defer guard.Unpatch()
			_, err := d.PosRecs(context.Background(), time.Now())
			ctx.Convey("Error should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRecContents(t *testing.T) {
	convey.Convey("RecContents", t, func(ctx convey.C) {
		ctx.Convey("When everything is correct", func() {
			res, aids, err := d.RecContents(context.Background(), time.Now())
			ctx.Convey("Error should be nil, res, aids should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
				ctx.So(aids, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When db.Query gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Query", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (*xsql.Rows, error) {
				return nil, fmt.Errorf("db.Query error")
			})
			defer guard.Unpatch()
			_, _, err := d.RecContents(context.Background(), time.Now())
			ctx.Convey("Error should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
