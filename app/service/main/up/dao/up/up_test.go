package up

import (
	"context"
	"testing"

	"go-common/app/service/main/up/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpAddUp(t *testing.T) {
	var (
		c = context.Background()
		u = &model.Up{
			MID:       int64(2089809),
			Attribute: 0,
		}
	)
	convey.Convey("AddUp", t, func(ctx convey.C) {
		id, err := d.AddUp(c, u)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(id, convey.ShouldNotBeNil)
		})
	})
}

func TestUpRawUp(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
	)
	convey.Convey("RawUp", t, func(ctx convey.C) {
		u, err := d.RawUp(c, mid)
		ctx.Convey("Then err should be nil.u should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(u, convey.ShouldNotBeNil)
		})
	})
}

func TestUpInfoActivitys(t *testing.T) {
	var (
		c      = context.Background()
		lastID = int64(0)
		ps     = 100
	)
	convey.Convey("UpInfoActivitys", t, func(ctx convey.C) {
		u, err := d.UpInfoActivitys(c, lastID, ps)
		ctx.Convey("Then err should be nil.u should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(u, convey.ShouldNotBeNil)
		})
	})
}
