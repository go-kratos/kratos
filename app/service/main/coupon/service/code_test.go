package service

import (
	"fmt"
	"go-common/app/service/main/coupon/model"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// go test  -test.v -test.run TestAddCartoonCoupon
func TestCaptchaToken(t *testing.T) {
	Convey("TestCaptchaToken ", t, func() {
		res, err := s.CaptchaToken(c, "")
		fmt.Println("token:", res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

// go test  -test.v -test.run TestUseCouponCode
func TestUseCouponCode(t *testing.T) {
	Convey("TestUseCouponCode ", t, func() {
		coupon, err := s.UseCouponCode(c, &model.ArgUseCouponCode{
			Token:  "a15e6f81374b4c5bb591ec3a1eba7461",
			Code:   "sasazxcvfdsa",
			Verify: "69vrz",
			IP:     "",
			Mid:    1,
		})
		fmt.Println("coupon:", coupon)
		So(err, ShouldBeNil)
		So(coupon, ShouldNotBeNil)
	})
}
