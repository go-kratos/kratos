package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoTagByNames(t *testing.T) {
	var (
		c     = context.TODO()
		names = []string{"22", "33"}
	)
	convey.Convey("TagByNames", t, func(ctx convey.C) {
		res, missed, err := d.TagByNames(c, names)
		ctx.Convey("Then err should be nil.res,missed should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(missed, convey.ShouldNotBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTags(t *testing.T) {
	var (
		c   = context.TODO()
		ids = []int64{1, 2, 3}
	)
	convey.Convey("Tags", t, func(ctx convey.C) {
		res, err := d.Tags(c, ids)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
