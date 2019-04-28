package dao

import (
	"context"
	"testing"

	"go-common/app/interface/main/web/model"

	"github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

func TestDaoWxHot(t *testing.T) {
	convey.Convey("WxHot", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("GET", d.wxHotURL).Reply(200).JSON(`{"code":0,"list":[{"aid":111,"score":10},{"aid":2222,"score":20}]}`)
			aids, err := d.WxHot(c)
			ctx.Convey("Then err should be nil.aids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(aids, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSetWxHotCache(t *testing.T) {
	convey.Convey("SetWxHotCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			arcs = []*model.WxArchive{{Aid: 1111}, {Aid: 2222}}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetWxHotCache(c, arcs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoWxHotCache(t *testing.T) {
	convey.Convey("WxHotCache", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			arcs, err := d.WxHotCache(c)
			ctx.Convey("Then err should be nil.arcs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.Printf("%+v", arcs)
			})
		})
	})
}

func TestDaoWxHotBakCache(t *testing.T) {
	convey.Convey("WxHotBakCache", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			arcs, err := d.WxHotBakCache(c)
			ctx.Convey("Then err should be nil.arcs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.Printf("%+v", arcs)
			})
		})
	})
}
