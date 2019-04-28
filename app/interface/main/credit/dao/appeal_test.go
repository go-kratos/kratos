package dao

import (
	"context"
	"testing"

	gock "gopkg.in/h2non/gock.v1"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddAppeal(t *testing.T) {
	convey.Convey("AddAppeal", t, func(convCtx convey.C) {
		var (
			c        = context.Background()
			tid      = int64(0)
			btid     = int64(0)
			oid      = int64(0)
			mid      = int64(0)
			business = int64(0)
			content  = ""
			reason   = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			defer gock.OffAll()
			httpMock("POST", d.addAppealURL).Reply(200).JSON(`{"code": 0}`)
			err := d.AddAppeal(c, tid, btid, oid, mid, business, content, reason)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAppealList(t *testing.T) {
	convey.Convey("AppealList", t, func(convCtx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(0)
			business = int(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			defer gock.OffAll()
			httpMock("GET", d.appealListURL).Reply(200).JSON(`{"code": 0,"data":[]}`)
			as, err := d.AppealList(c, mid, business)
			convCtx.Convey("Then err should be nil.as should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(as, convey.ShouldNotBeNil)
			})
			httpMock("GET", d.appealListURL).Reply(200).JSON(`{"code": 0}`)
			as, err = d.AppealList(c, mid, business)
			convCtx.Convey("Then err should be nil.as should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(as, convey.ShouldBeNil)
			})
		})
	})
}
