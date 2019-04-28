package dao

import (
	"go-common/app/admin/ep/merlin/model"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	username = "fengyifenguitest@bilibili.com"
)

func Test_Mail_Log(t *testing.T) {
	Convey("test add mail log", t, func() {
		ml := &model.MailLog{
			ReceiverName: username,
			MailType:     1,
			SendContext:  "test add mail log",
		}
		err := d.InsertMailLog(ml)
		So(err, ShouldBeNil)
	})

	Convey("test find mail log", t, func() {
		mailLogs, err := d.FindMailLog(username)
		So(len(mailLogs), ShouldBeGreaterThan, 0)
		So(err, ShouldBeNil)
	})

	Convey("test delete mail log", t, func() {
		err := d.DelMailLog(username)
		So(err, ShouldBeNil)
	})

	Convey("test find mail log", t, func() {
		mailLogs, err := d.FindMailLog(username)
		So(len(mailLogs), ShouldEqual, 0)
		So(err, ShouldBeNil)
	})
}
