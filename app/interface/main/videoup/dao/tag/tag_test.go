package tag

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

func TestTagUpBind(t *testing.T) {
	convey.Convey("UpBind", t, func(ctx convey.C) {
		var (
			err        error
			c          = context.Background()
			mid        = int64(2089809)
			aid        = int64(10110826)
			tags       = "LOL"
			regionName = "游戏"
			ip         = "127.0.0.1"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("POST", d.upBindURL).Reply(200).JSON(`{"code":20042,"data":""}`)
			err = d.UpBind(c, mid, aid, tags, regionName, ip)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestTagCheck(t *testing.T) {
	convey.Convey("TagCheck", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			mid     = int64(2089809)
			tagName = "LOL"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("Get", d.TagCheckURL).Reply(200).JSON(`{"code":20042,"data":""}`)
			no, err := d.TagCheck(c, mid, tagName)
			ctx.Convey("Then err should be nil.no should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(no, convey.ShouldBeNil)
			})
		})
	})
}
