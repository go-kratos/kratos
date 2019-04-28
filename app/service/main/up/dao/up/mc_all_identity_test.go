package up

import (
	"context"
	"go-common/app/service/main/up/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpkeyIdentityAll(t *testing.T) {
	var (
		mid = int64(0)
	)
	convey.Convey("keyIdentityAll", t, func(ctx convey.C) {
		p1 := keyIdentityAll(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestUpIdentityAllCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(1)
		st  = &model.IdentifyAll{}
	)
	convey.Convey("IdentityAllCache", t, func(ctx convey.C) {
		err := d.AddIdentityAllCache(c, mid, st)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		st, err := d.IdentityAllCache(c, mid)
		ctx.Convey("Then err should be nil.st should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(st, convey.ShouldNotBeNil)
		})
	})
}

func TestUpAddIdentityAllCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		st  = &model.IdentifyAll{}
	)
	convey.Convey("AddIdentityAllCache", t, func(ctx convey.C) {
		err := d.AddIdentityAllCache(c, mid, st)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
