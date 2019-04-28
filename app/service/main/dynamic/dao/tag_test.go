package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoHot(t *testing.T) {
	convey.Convey("Hot", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		convey.Convey("When http request gets code == 0", func(ctx convey.C) {
			httpMock("GET", d.hotURI).Reply(200).JSON(`{"code":0,"message":"0","ttl":1,"data":{"121":[25638,830,147345]}}`)
			res, err := d.Hot(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(len(res), convey.ShouldBeGreaterThan, 0)
			})
		})
	})
}

func TestDaoRids(t *testing.T) {
	convey.Convey("Rids", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		convey.Convey("When http request gets code == 0", func(ctx convey.C) {
			httpMock("GET", d.pridURI).Reply(200).JSON(`{"code":0,"message":"0","ttl":1,"data":[1,2,3,5]}`)
			res, err := d.Rids(c)
			ctx.Convey("Then err should be nil.res should greater 0.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(len(res), convey.ShouldBeGreaterThan, 0)
			})
		})
	})
}
