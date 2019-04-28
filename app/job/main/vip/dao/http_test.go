package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDaoAutoRenewPay(t *testing.T) {
	Convey("TestDaoAutoRenewPay", t, func() {
		res, err := d.AutoRenewPay(context.Background(), 1234)
		t.Logf("%+v,%+v", res, err)
		So(res, ShouldNotBeNil)
	})
}
