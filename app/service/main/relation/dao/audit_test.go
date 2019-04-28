package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoPassportDetail(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
		ip  = ""
	)
	convey.Convey("PassportDetail", t, func(cv convey.C) {
		res, err := d.PassportDetail(c, mid, ip)
		cv.Convey("Then err should be nil.res should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(res, convey.ShouldNotBeNil)
		})
	})
}
