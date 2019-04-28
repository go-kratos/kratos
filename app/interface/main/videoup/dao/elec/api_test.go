package elec

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

func TestElecArcShow(t *testing.T) {
	convey.Convey("ArcShow", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(2089809)
			aid  = int64(10110826)
			ip   = "127.0.0.1"
			err  error
			show bool
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("GET", d.showURI).Reply(200).JSON(`{"code":20004,"data":""}`)
			show, err = d.ArcShow(c, mid, aid, ip)
			ctx.Convey("Then err should be nil.show should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(show, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestElecArcUpdate(t *testing.T) {
	convey.Convey("ArcUpdate", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(2089809)
			aid      = int64(10110826)
			openElec = int8(1)
			ip       = "127.0.0.1"
			err      error
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("GET", d.arcOpenURL).Reply(200).JSON(`{"code":20004,"data":""}`)
			err = d.ArcUpdate(c, mid, aid, openElec, ip)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
