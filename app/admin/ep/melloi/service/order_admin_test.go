package service

import (
	"testing"

	"go-common/app/admin/ep/melloi/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	oa = model.OrderAdmin{
		UserName: "hujianping",
	}
)

func Test_OrderAdmin(t *testing.T) {

	Convey("test AddOrderAdmin 60003", t, func() {
		err := s.AddOrderAdmin(&oa)
		So(err, ShouldNotBeNil)
	})

	Convey("test QueryOrderAdmin 60003", t, func() {
		var admin *model.OrderAdmin
		admin, _ = s.QueryOrderAdmin(oa.UserName)
		So(admin, ShouldNotBeNil)
	})

	Convey("test AddOrderAdmin", t, func() {
		var oak = model.OrderAdmin{
			UserName: "hukai",
		}
		err := s.AddOrderAdmin(&oak)
		So(err, ShouldBeNil)
	})
}
