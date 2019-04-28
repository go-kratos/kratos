package dao

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSendMail(t *testing.T) {
	var (
		body    = strconv.Itoa(rand.Intn(100))
		subject = "邮件测试hyy"
		send    = "2459593393@qq.com"
	)
	convey.Convey("SendMail", t, func(ctx convey.C) {
		err := d.SendMail(body, subject, send)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
