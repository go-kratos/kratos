package vip

import (
	"context"
	"testing"

	vipmol "go-common/app/service/main/vip/model"
	"go-common/library/ecode"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_CodeVerify(t *testing.T) {
	Convey("code verify", t, func() {
		_, err := d.CodeVerify(context.TODO())
		So(err, ShouldBeNil)
	})
}

//  go test  -test.v -test.run TestDaoCancelUseCoupon
func TestDaoCancelUseCoupon(t *testing.T) {
	Convey("TestDaoCancelUseCoupon", t, func() {
		err := d.CancelUseCoupon(context.TODO(), &vipmol.ArgCancelUseCoupon{
			CouponToken: "672889783020180721180426",
			Mid:         1,
		})
		So(err == ecode.CouPonTokenNotFoundErr, ShouldBeTrue)
	})
}
