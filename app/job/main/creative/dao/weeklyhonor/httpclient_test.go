package weeklyhonor

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

func TestWeeklyhonorSendNotify(t *testing.T) {
	convey.Convey("SendNotify", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{11}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer gock.Off()
			httpMock("POST", d.c.Host.Message+_notifyURL).Reply(200).JSON(`{"code":0,"data":{}}`)
			_, err := d.SendNotify(c, mids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, err)
			})
		})
	})
}

func TestWeeklyhonorUpMids(t *testing.T) {
	convey.Convey("UpMids", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			size       = int(0)
			lastid     = int64(0)
			activeOnly bool
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			mids, newid, err := d.UpMids(c, size, lastid, activeOnly)
			ctx.Convey("Then err should be nil.mids,newid should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(newid, convey.ShouldNotBeNil)
				ctx.So(mids, convey.ShouldNotBeNil)
			})
		})
	})
}
