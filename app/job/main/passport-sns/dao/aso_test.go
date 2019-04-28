package dao

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"reflect"
	"testing"

	"go-common/app/job/main/passport-sns/model"
	xsql "go-common/library/database/sql"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

func TestDao_AddSnsLog(t *testing.T) {
	convey.Convey("AddSnsLog", t, func(ctx convey.C) {
		var (
			c = context.Background()
			a = &model.SnsLog{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			mock := monkey.PatchInstanceMethod(reflect.TypeOf(d.asoDB), "Exec", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (res sql.Result, err error) {
				return driver.RowsAffected(1), nil
			})
			defer mock.Unpatch()
			affected, err := d.AddSnsLog(c, a)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
			ctx.Convey("Then affected should not be nil.", func(ctx convey.C) {
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})

	})
}

func TestDao_AsoAccountSns(t *testing.T) {
	convey.Convey("AsoAccountSns", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			start = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.AsoAccountSns(c, start)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDao_AsoAccountSnsAll(t *testing.T) {
	convey.Convey("AsoAccountSnsAll", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			start = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.AsoAccountSnsAll(c, start)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
