package email

import (
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestEmailSendMail(t *testing.T) {
	convey.Convey("SendMail", t, func(ctx convey.C) {
		var (
			date    = time.Now()
			body    = ""
			subject = ""
			send    = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SendMail(date, body, subject, send)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestEmailSendMailAttach(t *testing.T) {
	convey.Convey("SendMailAttach", t, func(ctx convey.C) {
		var (
			filename = ""
			subject  = ""
			send     = []string{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SendMailAttach(filename, subject, send)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
