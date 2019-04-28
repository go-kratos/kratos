package dao

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	xsql "go-common/library/database/sql"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

func TestDaoTagBySsnID(t *testing.T) {
	convey.Convey("Given function TagBySsnID", t, func(ctx convey.C) {
		ctx.Convey("When everyting is correct", func(ctx convey.C) {
			tagMap, err := d.TagBySsnID(context.Background(), []int64{1, 2})
			ctx.Convey("Then err should be nil.tapMap should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(tagMap, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When d.dbMs.Query gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.dbMs), "Query", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (*xsql.Rows, error) {
				return nil, fmt.Errorf("d.dbMs.Query Error")
			})
			defer guard.Unpatch()
			_, err := d.TagBySsnID(context.Background(), []int64{1, 2})
			ctx.Convey("Then err should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
