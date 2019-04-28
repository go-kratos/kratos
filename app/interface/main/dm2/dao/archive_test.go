package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaotypesURI(t *testing.T) {
	convey.Convey("typesURI", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := testDao.typesURI()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTypeMapping(t *testing.T) {
	convey.Convey("TypeMapping", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rmap, err := testDao.TypeMapping(c)
			ctx.Convey("Then err should be nil.rmap should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rmap, convey.ShouldNotBeNil)
			})
		})
	})
}
