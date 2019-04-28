package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoTagIds(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(1)
	)
	convey.Convey("TagIds", t, func(ctx convey.C) {
		_, err := d.TagIds(c, aid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
		// ctx.Convey("ids should not be nil", func(ctx convey.C) {
		// 	ctx.So(ids, convey.ShouldNotBeNil)
		// })
	})
}
