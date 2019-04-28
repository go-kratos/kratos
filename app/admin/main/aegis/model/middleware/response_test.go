package middleware

import (
	"testing"

	"fmt"
	"github.com/smartystreets/goconvey/convey"
	"go-common/app/admin/main/aegis/model"
	"go-common/app/admin/main/aegis/model/resource"
)

func TestMiddlewareRequest(t *testing.T) {
	convey.Convey("Request", t, func(ctx convey.C) {
		opt := new(model.SearchParams)
		opt2 := &model.SearchParams{
			Extra1: "2",
		}
		ds.Encode = false
		Request(opt, &ds)
		Request(opt2, &ds)

		ctx.Convey("No return values", func(ctx convey.C) {
			ctx.So(opt.Extra1, convey.ShouldEqual, "")
			ctx.So(opt2.Extra1, convey.ShouldEqual, ds.Cfg[0].Hitv)
		})
	})
}

func TestMiddlewareResponse(t *testing.T) {
	convey.Convey("Response", t, func(ctx convey.C) {
		opt := new(model.AuditInfo)
		opt2 := &model.AuditInfo{
			Resource: &resource.Res{Extra1: 4},
		}
		ds.Encode = true
		Response(opt, nil, nil, &ds)
		Response(opt2, nil, nil, &ds)

		ctx.Convey("No return values", func(ctx convey.C) {
			ctx.So(opt.Resource, convey.ShouldBeNil)
			ctx.So(fmt.Sprintf("%v", opt2.Resource.Extra1), convey.ShouldEqual, ds.Cfg[0].Mapv)
		})
	})
}
