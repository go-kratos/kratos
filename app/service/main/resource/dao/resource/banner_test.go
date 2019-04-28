package resource

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	xsql "go-common/library/database/sql"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

func TestResourceBanner(t *testing.T) {
	convey.Convey("Banner", t, func(ctx convey.C) {
		ctx.Convey("When Everything is correct", func(ctx convey.C) {
			_, err := d.Banner(context.Background())
			ctx.Convey("Error shouold be nil, res should not be empty(No Data)", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Convey("When db.Query gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Query", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (*xsql.Rows, error) {
				return nil, fmt.Errorf("db.Query error")
			})
			defer guard.Unpatch()
			_, err := d.Banner(context.Background())
			ctx.Convey("Error should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestResourceCategory(t *testing.T) {
	convey.Convey("Category", t, func(ctx convey.C) {
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			_, err := d.Category(context.Background())
			ctx.Convey("Error shouold be nil, res should not be empty(No Data)", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Convey("When db.Query gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Query", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (*xsql.Rows, error) {
				return nil, fmt.Errorf("db.Query error")
			})
			defer guard.Unpatch()
			_, err := d.Category(context.Background())
			ctx.Convey("Error should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestResourceLimit(t *testing.T) {
	convey.Convey("Limit", t, func(ctx convey.C) {
		res, err := d.Limit(context.Background())
		ctx.Convey("Error shouold be nil, res should not be empty", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeEmpty)
		})
	})
}
