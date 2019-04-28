package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoUPGroups(t *testing.T) {
	var (
		c    = context.TODO()
		mids = []int64{27515256, 27515304}
	)
	convey.Convey("UPGroups", t, func(ctx convey.C) {
		httpMock("GET", d.upGroupURL).Reply(200).JSON(`{"code":0,"data":{"items":[{"group_id":1,"group_name":"name","group_tag":"tag","font_color":"font","bg_color":"bg","note":"note","mid":11}]}}`)
		_, err := d.UPGroups(c, mids)
		ctx.Convey("Then err should be nil.groups should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
