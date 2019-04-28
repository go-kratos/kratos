package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

func TestDaotagAidKey(t *testing.T) {
	convey.Convey("tagAidKey", t, func(ctx convey.C) {
		var (
			tid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := tagAidKey(tid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTagAids(t *testing.T) {
	convey.Convey("TagAids", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("GET", d.c.Host.Data+_tagFeedURL).Reply(200).JSON(`{"code":0,"data":[1111,2222],"total":2}`)
			res, err := d.TagAids(c, tid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSetTagAidsBakCache(t *testing.T) {
	convey.Convey("SetTagAidsBakCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			tid  = int64(2222)
			aids = []int64{111, 222}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetTagAidsBakCache(c, tid, aids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoTagAidsBakCache(t *testing.T) {
	convey.Convey("TagAidsBakCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tid = int64(2222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.TagAidsBakCache(c, tid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
