package http

import (
	"context"
	"reflect"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestHttpFilterMulti(t *testing.T) {
	convey.Convey("FilterMulti", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			area       = ""
			msg        = ""
			successRes = `{"code":0,"data":[{"level":10,"msg":"女装大佬"}]}`
			failRes    = `{"code":200,"data":[{"level":10,"msg":"女装大佬"}]}`
		)
		ctx.Convey("success", func(ctx convey.C) {
			httpMock("POST", d.c.Host.API+_filterURI).Reply(200).JSON(successRes)
			hits, err := d.FilterMulti(c, area, msg)
			ctx.Convey("Then err should be nil.hits should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(reflect.DeepEqual(hits, []string{"女装大佬"}), convey.ShouldEqual, true)
			})
		})
		ctx.Convey("request fail", func(ctx convey.C) {
			httpMock("POST", d.c.Host.API+_filterURI).Reply(504)
			_, err := d.FilterMulti(c, area, msg)
			ctx.Convey("Then err should be nil.hits should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("business fail", func(ctx convey.C) {
			httpMock("POST", d.c.Host.API+_filterURI).Reply(200).JSON(failRes)
			_, err := d.FilterMulti(c, area, msg)
			ctx.Convey("Then err should be nil.hits should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
