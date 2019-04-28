package newbiedao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestNewbiedaoGetRelations(t *testing.T) {
	convey.Convey("GetRelations", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(27515398)
			fids = []int64{389088, 6810019, 4578433}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.GetRelations(c, mid, fids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
