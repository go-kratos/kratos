package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetPassportDetail(t *testing.T) {
	convey.Convey("passport", t, func() {
		var mid int64 = 27515586
		httpMock("GET", _passportURL).Reply(200).BodyString(`{"code":0}`)
		res, err := d.PassportDetail(context.TODO(), mid)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldBeNil)
	})
}
