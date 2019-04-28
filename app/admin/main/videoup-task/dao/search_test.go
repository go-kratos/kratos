package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoOutTime(t *testing.T) {
	convey.Convey("OutTime", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			httpMock("GET", d.c.Host.Search+_searchURL).Reply(200).JSON(`{"code":0,"data":{"result":[{"uid":0}]}}`)
			mcases, err := d.OutTime(c, ids)
			ctx.Convey("Then err should be nil.mcases should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(mcases, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInQuitList(t *testing.T) {
	convey.Convey("InQuitList", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			uids = []int64{}
			bt   = ""
			et   = ""
		)

		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			httpMock("GET", d.c.Host.Search+_searchURL).Reply(200).JSON(`{"code":0,"data":{"result":[{"uid":0,"action":"0"}]}}`)
			l, err := d.InQuitList(c, uids, bt, et)
			ctx.Convey("Then err should be nil.l should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(l, convey.ShouldNotBeNil)
			})
		})
	})
}
