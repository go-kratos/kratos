package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAllTips(t *testing.T) {
	var (
		c   = context.TODO()
		now = int64(20)
	)
	convey.Convey("AllTips", t, func(ctx convey.C) {
		_, err := d.AllTips(c, now)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})

	})
}
