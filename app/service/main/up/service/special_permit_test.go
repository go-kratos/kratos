package service

import (
	"testing"

	"go-common/library/net/http/blademaster"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceSpecialGroupPermit(t *testing.T) {
	convey.Convey("SpecialGroupPermit", t, func(ctx convey.C) {
		var (
			c       = &blademaster.Context{}
			groupID = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := s.SpecialGroupPermit(c, groupID)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
