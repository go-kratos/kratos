package archive

import (
	"testing"

	arcwar "go-common/app/service/main/archive/api"

	"github.com/smartystreets/goconvey/convey"
)

func TestArchivevalidView(t *testing.T) {
	var (
		vp = &arcwar.ViewReply{}
	)
	convey.Convey("validView", t, func(ctx convey.C) {
		valid := validView(vp, true)
		ctx.Convey("Then valid should not be nil.", func(ctx convey.C) {
			ctx.So(valid, convey.ShouldNotBeNil)
		})
	})
}
