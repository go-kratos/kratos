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

func TestUsers(t *testing.T) {
	convey.Convey("Users", t, func(ctx convey.C) {
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			res, err := d.Users(context.Background())
			ctx.Convey("Error should be nil, res should not be empty", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeEmpty)
			})
		})
		ctx.Convey("When db.Query gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Query", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (*xsql.Rows, error) {
				return nil, fmt.Errorf("db.Query error")
			})
			defer guard.Unpatch()
			_, err := d.Users(context.Background())
			ctx.Convey("Error should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

// Set d.close() to get reversal case
func TestUser(t *testing.T) {
	var id = int64(67)
	convey.Convey("User", t, func(ctx convey.C) {
		convey.Convey("When everything is correct,", func(ctx convey.C) {
			asgs, err := d.User(context.Background(), id)
			ctx.Convey("Error should be nil, asgs should not be nil(No Data)", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(asgs, convey.ShouldBeNil)
			})
		})
		convey.Convey("When set db closed", WithReopenDB(func(d *Dao) {
			d.Close()
			_, err := d.User(context.Background(), id)
			convey.Convey("Error should not be nil", func(ctx convey.C) {
				convey.So(err, convey.ShouldNotBeNil)
			})
		}))
	})
}
