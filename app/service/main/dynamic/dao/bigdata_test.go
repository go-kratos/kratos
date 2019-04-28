package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoRegionArcs(t *testing.T) {
	convey.Convey("RegionArcs", t, func() {
		var (
			c        = context.Background()
			rid      = int32(1)
			remoteIP = "0.0.0.0"
		)
		convey.Convey("When http request gets 404 error", func(ctx convey.C) {
			httpMock("GET", d.regionURI).MatchParam("rid", "0").Reply(502).JSON(`{"code":-500}`)
			_, _, err := d.RegionArcs(c, rid, remoteIP)
			ctx.Convey("Then error should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})

		convey.Convey("When http request gets code == 0", func(ctx convey.C) {
			httpMock("GET", d.regionURI).MatchParam("rid", "1").Reply(200).JSON(`{"code":0,"message":"ok","data":{"avid":[28868622,29365735,29621166],"count":4602378}}`)
			aids, total, err := d.RegionArcs(c, rid, remoteIP)
			ctx.Convey("Then error should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldBeGreaterThan, 0)
				ctx.So(len(aids), convey.ShouldBeGreaterThan, 0)
			})
		})
	})
}

func TestDaoRegionTagArcs(t *testing.T) {
	convey.Convey("RegionTagArcs", t, func() {
		var (
			c        = context.Background()
			rid      = int32(136)
			tagID    = int64(4577)
			remoteIP = "0.0.0.0"
			params   = map[string]string{"rid": "136", "tag_id": "4577"}
		)
		convey.Convey("When http request gets code == 0", func(ctx convey.C) {
			httpMock("GET", d.regionTagURI).MatchParams(params).Reply(200).JSON(`{"code":0,"message":"ok","data":{"avid":[28868622,29365735,29621166],"count":4602378}}`)
			aids, err := d.RegionTagArcs(c, rid, tagID, remoteIP)
			ctx.Convey("Then error should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(len(aids), convey.ShouldBeGreaterThan, 0)
			})
		})
	})

}
