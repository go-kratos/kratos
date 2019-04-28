package pgc

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetLabelCache(t *testing.T) {
	convey.Convey("GetLabelCache", t, func(cx convey.C) {
		cx.Convey("When everything goes positive", func(cx convey.C) {
			res, err := d.GetLabelCache(ctx)
			cx.Convey("Then err should be nil.res should not be nil.", func(cx convey.C) {
				cx.So(err, convey.ShouldBeNil)
				cx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
