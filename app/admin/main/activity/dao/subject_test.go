package dao

import (
	"context"
	"testing"

	"fmt"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddTags(t *testing.T) {
	convey.Convey("AddTags", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			tags = "keai"
			ip   = "10.256.36.68"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddTags(c, tags, ip)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				fmt.Printf("%+v", err)
			})
		})
	})
}
