package monitor

import (
	"context"
	"gopkg.in/h2non/gock.v1"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestMonitorArchiveAttr(t *testing.T) {
	convey.Convey("ArchiveAttr", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			aid = int64(1212)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			defer gock.OffAll()
			httpMock("GET", d.URLArcAddit).Reply(200).JSON(`{"code":0,"data":{}}`)
			_, err := d.ArchiveAttr(c, aid)
			convCtx.Convey("Then err should be nil.addit should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
