package tag

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestTagInsertUpTagYear(t *testing.T) {
	convey.Convey("InsertUpTagYear", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			vals = "(111, 100)"
			col  = "tag1"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.InsertUpTagYear(c, vals, col)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
