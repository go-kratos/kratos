package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoLoadSecret(t *testing.T) {
	var c = context.TODO()
	convey.Convey("LoadSecret", t, func(ctx convey.C) {
		res, err := d.LoadSecret(c)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
