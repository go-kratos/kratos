package videoshot

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestVideoshothashKey(t *testing.T) {
	var (
		cid = int64(1)
	)
	convey.Convey("hashKey", t, func(ctx convey.C) {
		p1 := hashKey(cid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestVideoshotcache(t *testing.T) {
	var (
		c   = context.TODO()
		cid = int64(1)
	)
	convey.Convey("cache", t, func(ctx convey.C) {
		count, ver, err := d.cache(c, cid)
		ctx.Convey("Then err should be nil.count,ver should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ver, convey.ShouldNotBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestVideoshotaddCache(t *testing.T) {
	var (
		c     = context.TODO()
		cid   = int64(1)
		ver   = int(0)
		count = int(0)
	)
	convey.Convey("addCache", t, func(ctx convey.C) {
		err := d.addCache(c, cid, ver, count)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
