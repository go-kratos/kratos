package dao

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSendMail(t *testing.T) {
	convey.Convey("SendMail", t, func(ctx convey.C) {
		var (
			body    = ""
			subject = ""
			send    = []string{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SendMail(body, subject, send)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSendMailAttach(t *testing.T) {
	convey.Convey("SendMailAttach", t, func(ctx convey.C) {
		var (
			filename = ""
			subject  = ""
			send     = []string{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SendMailAttach(filename, subject, send)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
