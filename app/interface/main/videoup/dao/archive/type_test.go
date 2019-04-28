package archive

import (
	"context"
	"testing"

	"gopkg.in/h2non/gock.v1"

	"github.com/smartystreets/goconvey/convey"
)

func TestArchiveTypeMapping(t *testing.T) {
	convey.Convey("TypeMapping", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)

		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("GET", d.typesURI).Reply(200).JSON(`{"code":20001}`)
			rmap, err := d.TypeMapping(c)
			ctx.Convey("Then err should be nil.rmap should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(rmap, convey.ShouldBeNil)
			})
		})
	})
}
