package academy

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAcademyArchives(t *testing.T) {
	convey.Convey("Archives", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(10110127)
			bs    = int(1)
			limit = int(10)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.Archives(c, id, bs, limit)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, err)
			})
		})
	})
}

func TestAcademyUPHotByAIDs(t *testing.T) {
	convey.Convey("UPHotByAIDs", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			hots = map[int64]int64{
				10110127: 11,
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.UPHotByAIDs(c, hots)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, err)
			})
		})
	})
}
