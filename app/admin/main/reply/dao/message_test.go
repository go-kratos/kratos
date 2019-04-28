package dao

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMessage(t *testing.T) {
	var (
		mid   = int64(1)
		title = "test title"
		msg   = "test msg"
		now   = time.Now()
		c     = context.Background()
	)
	Convey("test send a message", t, WithDao(func(d *Dao) {
		d.SendReplyDelMsg(c, mid, title, msg, now)
		So(msg, ShouldEqual, "test msg")
		d.SendReportAcceptMsg(c, mid, title, msg, now)
		So(title, ShouldEqual, "test title")
	}))
}
