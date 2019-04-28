package newbiedao

import (
	"context"
	"go-common/app/interface/main/growup/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestNewbiedaoGetVideoUp(t *testing.T) {
	convey.Convey("GetVideoUp", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			req = &model.NewbieLetterReq{Aid: 10110467, Mid: 27515398}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			videoUpArchive, err := d.GetVideoUp(c, req.Aid)
			ctx.Convey("Then err should be nil.videoUpArchive should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(videoUpArchive, convey.ShouldNotBeNil)
			})
		})
	})
}
