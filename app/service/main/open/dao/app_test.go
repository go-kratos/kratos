package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSecret(t *testing.T) {
	convey.Convey("Secret", t, func() {
		res, err := d.Secret(context.TODO(), "19e7f5b7d8ad531b")
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}
