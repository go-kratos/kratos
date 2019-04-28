package archive

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

func TestArchiveHistoryList(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(888952460)
		aid = int64(10110560)
		ip  = "127.0.0.1"
	)
	convey.Convey("HistoryList", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("GET", d.hList).Reply(200).JSON(`{"code":20001}`)
		historys, err := d.HistoryList(c, mid, aid, ip)
		ctx.Convey("Then err should be nil.historys should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(historys, convey.ShouldBeNil)
		})
	})
}

func TestArchiveHistoryView(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(888952460)
		hid = int64(0)
		ip  = "127.0.0.1"
	)
	convey.Convey("HistoryView", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("GET", d.hView).Reply(200).JSON(`{"code":20001}`)
		history, err := d.HistoryView(c, mid, hid, ip)
		ctx.Convey("Then err should be nil.history should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(history, convey.ShouldBeNil)
		})
	})
}
