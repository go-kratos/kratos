package dao

import (
	"context"
	"go-common/app/job/main/passport-auth/model"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddRefresh(t *testing.T) {
	convey.Convey("AddRefresh", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			no      = &model.Refresh{}
			refresh = []byte("9df38fe4b94a47baad001ad823b84110")
			token   = []byte("61c13e530b1418653e2fdc265b3f0fe6")
			ct      = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			affected, err := d.AddRefresh(c, no, refresh, token, ct)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelRefresh(t *testing.T) {
	convey.Convey("DelRefresh", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			refresh = []byte("9df38fe4b94a47baad001ad823b84110")
			ct      = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			affected, err := d.DelRefresh(c, refresh, ct)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoformatRefreshSuffix(t *testing.T) {
	convey.Convey("formatRefreshSuffix", t, func(ctx convey.C) {
		var (
			no = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := formatRefreshSuffix(no)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoformatByDate(t *testing.T) {
	convey.Convey("formatByDate", t, func(ctx convey.C) {
		var (
			year  = int(0)
			month = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := formatByDate(year, month)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
