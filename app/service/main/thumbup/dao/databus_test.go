package dao

import (
	"context"
	"go-common/app/service/main/thumbup/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoPubStatDatabus(t *testing.T) {
	var (
		c        = context.TODO()
		business = "archive"
		mid      = int64(1)
		s        = &model.Stats{}
	)
	convey.Convey("PubStatDatabus", t, func(ctx convey.C) {
		err := d.PubStatDatabus(c, business, mid, s, 1)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
