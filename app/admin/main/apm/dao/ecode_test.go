package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetCodes(t *testing.T) {
	convey.Convey("GetCodes", t, func() {
		res, err := d.GetCodes(context.Background(), "0", "20000")
		t.Logf("res:%+v", res)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}
