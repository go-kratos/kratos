package vip

import (
	"context"
	"fmt"
	"go-common/app/interface/main/account/model"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// go test  -test.v -test.run TestServiceOpenIDByAuthCode
func TestServiceOpenIDByAuthCode(t *testing.T) {
	Convey("TestServiceOpenIDByAuthCode", t, func() {
		res, err := s.OpenIDByAuthCode(context.TODO(), &model.ArgAuthCode{})
		fmt.Println("res", res)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestServiceOpenAuthCallBack
func TestServiceOpenAuthCallBack(t *testing.T) {
	Convey("TestServiceOpenAuthCallBack", t, func() {
		err := s.OpenAuthCallBack(context.TODO(), &model.ArgOpenAuthCallBack{})
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestServiceBilibiliPrizeGrant
func TestServiceBilibiliPrizeGrant(t *testing.T) {
	Convey("TestServiceBilibiliPrizeGrant", t, func() {
		res, err := s.BilibiliPrizeGrant(context.TODO(), &model.ArgBilibiliPrizeGrant{
			PrizeKey: "coupon_ele_1",
			UniqueNo: "1x",
			OpenID:   "e11303e8c26268a6cbdc2dc7fce55199",
			AppID:    32,
		})
		fmt.Println("res:", res)
		So(err, ShouldBeNil)
	})
}

func TestServiceOpenBindByOutOpenID(t *testing.T) {
	Convey("TestServiceOpenBindByOutOpenID", t, func() {
		err := s.OpenBindByOutOpenID(context.TODO(), &model.ArgBind{
			AppID:     32,
			OutOpenID: "o8f999ad5d724b4a2ljbp7cm",
			OpenID:    "e11303e8c26268a6cbdc2dc7fce55199",
		})
		So(err, ShouldBeNil)
	})
}

func TestServiceElemeOAuthURI(t *testing.T) {
	Convey("TestServiceElemeOAuthURI", t, func() {
		url := s.ElemeOAuthURI(context.TODO(), "state")
		fmt.Println("url-------", url)
		So(url, ShouldNotBeNil)
	})
}
