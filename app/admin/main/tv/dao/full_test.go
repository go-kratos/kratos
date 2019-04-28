package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoFullImport(t *testing.T) {
	var (
		c     = context.Background()
		build = int(0)
	)
	convey.Convey("FullImport", t, func(ctx convey.C) {
		ctx.Convey("Http code err", func(ctx convey.C) {
			httpMock("GET", d.fullURL).Reply(-400).JSON(``)
			_, err := d.FullImport(c, build)
			ctx.So(err, convey.ShouldNotBeNil)
		})
		ctx.Convey("Business code err", func(ctx convey.C) {
			httpMock("GET", d.fullURL).Reply(200).JSON(`{"code":-400}`)
			_, err := d.FullImport(c, build)
			ctx.So(err, convey.ShouldNotBeNil)
		})
		ctx.Convey("Everything goes well", func(ctx convey.C) {
			httpMock("GET", d.fullURL).Reply(200).JSON(`{"code":0,"data":[{"id":1}]}`)
			data, err := d.FullImport(c, build)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(data, convey.ShouldNotBeNil)
		})
	})
}
