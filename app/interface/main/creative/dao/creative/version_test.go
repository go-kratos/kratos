package creative

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestCreativeAllByTypes(t *testing.T) {
	var (
		c      = context.TODO()
		tyInts = []int{7, 8}
	)
	convey.Convey("AllByTypes", t, func(ctx convey.C) {
		vss, err := d.AllByTypes(c, tyInts)
		ctx.Convey("Then err should be nil.vss should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(vss, convey.ShouldNotBeNil)
		})
	})
}

func TestCreativeLatestByType(t *testing.T) {
	var (
		c   = context.TODO()
		tid = int(7)
	)
	convey.Convey("LatestByType", t, func(ctx convey.C) {
		vs, err := d.LatestByType(c, tid)
		ctx.Convey("Then err should be nil.vs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(vs, convey.ShouldNotBeNil)
		})
	})
}
