package reply

import (
	"context"
	model "go-common/app/job/main/reply/model/reply"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestReplyPubEvent(t *testing.T) {
	convey.Convey("PubEvent", t, func(ctx convey.C) {
		var (
			action = ""
			mid    = int64(0)
			sub    = &model.Subject{}
			rp     = &model.Reply{}
			report = &model.Report{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.PubEvent(context.Background(), action, mid, sub, rp, report)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
