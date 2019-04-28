package dao

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/gomail.v2"
)

var (
	m = gomail.NewMessage()
)

func Test_SendMail(t *testing.T) {
	Convey("test send mail", t, func() {
		m.SetHeader("To", "fengyifeng@bilibili.com")
		m.SetHeader("Subject", "unit test")
		d.SendMail(m)
	})
}
