package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoMFilter(t *testing.T) {
	var (
		c    = context.Background()
		msgs = []string{"22娘", "unit test", "ut"}
	)
	convey.Convey("MFilter", t, func(ctx convey.C) {
		checked, err := d.MFilter(c, msgs)
		ctx.Convey("Then err should be nil.checked should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(checked, convey.ShouldNotBeEmpty)
		})
	})
}
