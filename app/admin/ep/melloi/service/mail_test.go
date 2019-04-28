package service

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	receiver = "hujianping"
	subject  = "test of melloi"
	content  = "this is test content"
)

func Test_Mail(t *testing.T) {
	Convey("test QueryOrder", t, func() {
		err := s.SendMail(receiver, subject, content)
		So(err, ShouldBeNil)
	})
}
