package dao

import (
	"context"
	"testing"

	"go-common/app/admin/main/reply/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEvent(t *testing.T) {
	var (
		mid    = int64(1)
		sub    = &model.Subject{}
		rp     = &model.Reply{Content: &model.ReplyContent{}}
		report = &model.Report{}
		c      = context.Background()
	)
	Convey("pub a event", t, WithDao(func(d *Dao) {
		err := d.PubEvent(c, model.EventReportAdd, mid, sub, rp, report)
		So(err, ShouldBeNil)
	}))
}
