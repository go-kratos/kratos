package dao

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

	Convey("test AddOrderAdmin", t, func() {
		err := d.AddOrderAdmin(&oa)
		So(err, ShouldBeNil)
	})

	Convey("test QueryOrderAdmin", t, func() {
		var admin *model.OrderAdmin
		admin, _ = d.QueryOrderAdmin(oa.UserName)
		So(admin, ShouldNotBeNil)
	})
}
