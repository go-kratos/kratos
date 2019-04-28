package subtitle

import (
	"context"
	"go-common/app/interface/main/dm2/model"
	"go-common/app/interface/main/dm2/rpc/client"
	"go-common/library/ecode"
	"reflect"
	"testing"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

func TestSubtitleView(t *testing.T) {
	convey.Convey("View", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			aid = int64(24416433)
		)
		ctx.Convey("testSubtitleView", func(ctx convey.C) {
			mock := monkey.PatchInstanceMethod(reflect.TypeOf(d.sub), "SubtitleSubjectSubmitGet",
				func(_ *client.Service, _ context.Context, _ *model.ArgArchiveID) (res *model.SubtitleSubjectReply, err error) {
					return nil, ecode.CreativeSubtitleAPIErr
				})
			defer mock.Unpatch()
			res, err := d.View(c, aid)
			ctx.Convey("testSubtitleView", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldBeNil)
			})
		})
	})
}
