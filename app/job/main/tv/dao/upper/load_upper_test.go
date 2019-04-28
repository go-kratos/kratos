package upper

import (
	"context"
	"testing"

	ugcMdl "go-common/app/job/main/tv/model/ugc"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpperupperMetaKey(t *testing.T) {
	var MID = int64(0)
	convey.Convey("upperMetaKey", t, func(c convey.C) {
		p1 := upperMetaKey(MID)
		c.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestUpperLoadUpMeta(t *testing.T) {
	var (
		ctx = context.Background()
		mid = int64(0)
	)
	convey.Convey("LoadUpMeta", t, func(c convey.C) {
		upper, err := d.LoadUpMeta(ctx, mid)
		c.Convey("Then err should be nil.upper should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(upper, convey.ShouldNotBeNil)
		})
	})
}

func TestUpperupMetaCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
	)
	convey.Convey("upMetaCache", t, func(ctx convey.C) {
		upper, err := d.upMetaCache(c, mid)
		ctx.Convey("Then err should be nil.upper should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(upper, convey.ShouldNotBeNil)
		})
	})
}

func TestUpperupMetaDB(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
	)
	convey.Convey("upMetaDB", t, func(ctx convey.C) {
		upper, err := d.upMetaDB(c, mid)
		ctx.Convey("Then err should be nil.upper should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(upper, convey.ShouldNotBeNil)
		})
	})
}

func TestUpperaddUpMetaCache(t *testing.T) {
	var (
		upper = &ugcMdl.Upper{}
	)
	convey.Convey("addUpMetaCache", t, func(ctx convey.C) {
		err := d.addUpMetaCache(context.Background(), upper)
		ctx.Convey("Then err should be nil.upper should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
