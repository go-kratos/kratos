package web

import (
	"context"
	"net/url"
	"testing"

	"go-common/app/interface/main/web-goblin/model/web"

	"github.com/smartystreets/goconvey/convey"
)

func TestWebRecruit(t *testing.T) {
	convey.Convey("Recruit", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			params = url.Values{}
			ru     = &web.Params{
				Route: "v1/jobs",
			}
		)
		params.Set("mode", "social")
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			httpMock("GET", "https://api.mokahr.com/v1/jobs/bilibili").Reply(200).JSON(`{jobs:[], "total": 245}`)
			res, err := d.Recruit(c, params, ru)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
