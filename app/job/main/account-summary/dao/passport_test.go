package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoPassportProfile(t *testing.T) {
	var (
		// ctx = context.Background()
		mid = int64(0)
	)
	convey.Convey("PassportProfile", t, func(ctx convey.C) {
		p1, err := d.PassportProfile(context.TODO(), mid)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAsoAccountRegOrigin(t *testing.T) {
	var (
		// ctx = context.Background()
		mid = int64(0)
	)
	convey.Convey("AsoAccountRegOrigin", t, func(ctx convey.C) {
		p1, err := d.AsoAccountRegOrigin(context.TODO(), mid)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}
