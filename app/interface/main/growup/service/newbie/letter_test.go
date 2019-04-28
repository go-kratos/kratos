package newbie

import (
	"context"
	"go-common/app/interface/main/growup/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestNewbieLetter(t *testing.T) {
	convey.Convey("Letter", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			req = &model.NewbieLetterReq{Aid: 10110467, Mid: 27515398}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := s.Letter(c, req)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
