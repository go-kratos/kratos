package coupon

import (
	"context"
	"flag"
	"fmt"
	"testing"

	"go-common/app/interface/main/account/conf"
	v1 "go-common/app/service/main/coupon/api"
	"go-common/app/service/main/coupon/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	flag.Set("conf", "../../cmd/account-interface-example.toml")
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	s = New(conf.Conf)
}

// go test  -test.v -test.run TestServiceAllowanceList
func TestServiceAllowanceList(t *testing.T) {
	Convey("TestServiceAllowanceList", t, func() {
		res, err := s.AllowanceList(context.TODO(), int64(1), int8(0))
		t.Logf("%v", res)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestCouponPage
// func TestCouponPage(t *testing.T) {
// 	Convey("TestCouponPage", t, func() {
// 		res, err := s.CouponPage(context.TODO(), 1, int8(0), 1, 10)
// 		t.Logf("%v", res)
// 		So(err, ShouldBeNil)
// 	})
// }

// // go test  -test.v -test.run TestCouponCartoonPage
// func TestCouponCartoonPage(t *testing.T) {
// 	Convey("TestCouponCartoonPage", t, func() {
// 		res, err := s.CouponCartoonPage(context.TODO(), 1, int8(0), 1, 10)
// 		t.Logf("%v", res)
// 		So(err, ShouldBeNil)
// 	})
// }

// go test  -test.v -test.run TestServiceCaptchaToken
func TestServiceCaptchaToken(t *testing.T) {
	Convey("TestServiceCaptchaToken", t, func() {
		res, err := s.CaptchaToken(context.Background(), &v1.CaptchaTokenReq{Ip: ""})
		fmt.Println("res:", res)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestServiceUseCouponCode
func TestServiceUseCouponCode(t *testing.T) {
	Convey("TestServiceUseCouponCode", t, func() {
		res, err := s.UseCouponCode(context.Background(), &model.ArgUseCouponCode{
			IP:     "",
			Token:  "927a6ea6e9d64e929beadfba6d2bd491",
			Code:   "sasazxcvfdsa",
			Verify: "e8z90",
			Mid:    1,
		})
		fmt.Println("res:", res)
		So(err, ShouldBeNil)
	})
}
