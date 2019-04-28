package danmu

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

func TestDanmuGetAdvDmPurchases(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		ip  = "127.0.0.1"
	)
	convey.Convey("GetAdvDmPurchases", t, func(ctx convey.C) {
		danmus, err := d.GetAdvDmPurchases(c, mid, ip)
		ctx.Convey("Then err should be nil.danmus should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(danmus, convey.ShouldNotBeNil)
		})
	})
}

func TestDanmuPassAdvDmPurchase(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		id  = int64(1)
		ip  = "127.0.0.1"
	)
	convey.Convey("PassAdvDmPurchase", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("POST", d.advDmPurchasePassURL).Reply(200).JSON(`{"code":20043,"data":""}`)
		err := d.PassAdvDmPurchase(c, mid, id, ip)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestDanmuDenyAdvDmPurchase(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		id  = int64(1)
		ip  = "127.0.0.1"
	)
	convey.Convey("DenyAdvDmPurchase", t, func(ctx convey.C) {
		err := d.DenyAdvDmPurchase(c, mid, id, ip)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDanmuCancelAdvDmPurchase(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		id  = int64(1)
		ip  = "127.0.0.1"
	)
	convey.Convey("CancelAdvDmPurchase", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("POST", d.advDmPurchaseCancelURL).Reply(200).JSON(`{"code":0,"data":""}`)
		err := d.CancelAdvDmPurchase(c, mid, id, ip)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
