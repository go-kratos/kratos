package share

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestShareredisKey(t *testing.T) {
	var (
		aid = int64(1)
	)
	convey.Convey("redisKey", t, func(ctx convey.C) {
		p1 := redisKey(aid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestShareredisFKey(t *testing.T) {
	var (
		mid  = int64(1)
		date = ""
	)
	convey.Convey("redisFKey", t, func(ctx convey.C) {
		p1 := redisFKey(mid, date)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestShareAddShare(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(1)
		aid = int64(1)
		ip  = ""
	)
	convey.Convey("AddShare", t, func(ctx convey.C) {
		ok, err := d.AddShare(c, mid, aid, ip)
		ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ok, convey.ShouldNotBeNil)
		})
	})
}

func TestShareHadFirstShare(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(1)
		aid = int64(1)
		ip  = ""
	)
	convey.Convey("HadFirstShare", t, func(ctx convey.C) {
		had, err := d.HadFirstShare(c, mid, aid, ip)
		ctx.Convey("Then err should be nil.had should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(had, convey.ShouldNotBeNil)
		})
	})
}
