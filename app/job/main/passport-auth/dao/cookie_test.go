package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/job/main/passport-auth/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddCookie(t *testing.T) {
	convey.Convey("AddCookie", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			cookie  = &model.Cookie{}
			session = []byte("712b7a22,1535703191,c07e44d8")
			csrf    = []byte("0273f9216fa8d6d77c3dd5499a1d0d4a")
			ct      = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			affected, err := d.AddCookie(c, cookie, session, csrf, ct)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelCookie(t *testing.T) {
	convey.Convey("DelCookie", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			session = []byte("712b7a22,1535703191,c07e44d8")
			ct      = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			affected, err := d.DelCookie(c, session, ct)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddCookieDeleted(t *testing.T) {
	convey.Convey("AddCookieDeleted", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			cookie  = &model.Cookie{}
			session = []byte("712b7a22,1535703191,c07e44d8")
			csrf    = []byte("0273f9216fa8d6d77c3dd5499a1d0d4a")
			ct      = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			affected, err := d.AddCookieDeleted(c, cookie, session, csrf, ct)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}
