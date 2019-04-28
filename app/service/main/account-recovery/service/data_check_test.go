package service

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceCompareInfo(t *testing.T) {
	var (
		c   = context.Background()
		rid = int64(240)
	)
	convey.Convey("CompareInfo", t, func(ctx convey.C) {
		err := s.CompareInfo(c, rid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
