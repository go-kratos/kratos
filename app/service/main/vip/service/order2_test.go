package service

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	memrpc "go-common/app/service/main/member/api/gorpc"
	memmdl "go-common/app/service/main/member/model"
	"go-common/app/service/main/vip/model"
	"go-common/library/ecode"

	"github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
)

//  go test  -test.v -test.run TestServiceCreateOrder2
func TestServiceCreateOrder2(t *testing.T) {
	Convey("TestServiceCreateOrder2 ", t, func() {
		arg := &model.ArgCreateOrder2{
			Mid:       1,
			Month:     3,
			Platform:  "ios",
			MobiApp:   "iphone",
			Device:    "phone",
			AppID:     1,
			AppSubID:  "11",
			OrderType: 0,
		}
		// ios
		r, _, err := s.CreateOrder2(c, arg)
		t.Logf("ios r(%+v)", r)
		t.Logf("ios PayParam(%+v)", r.PayParam)
		So(err, ShouldBeNil)
		// ios auto
		arg.OrderType = 1
		arg.Month = 1
		r, _, err = s.CreateOrder2(c, arg)
		t.Logf("ios auto r(%+v)", r)
		t.Logf("ios auto PayParam(%+v)", r.PayParam)
		So(err, ShouldBeNil)
		// pc
		arg.Platform = ""
		arg.MobiApp = ""
		arg.Device = ""
		arg.OrderType = 0
		arg.ReturnURL = "http://www.bilibili.com/"
		r, _, err = s.CreateOrder2(c, arg)
		t.Logf("pc r(%+v)", r)
		t.Logf("pc PayParam(%+v)", r.PayParam)
		So(err, ShouldBeNil)
	})
}

//  go test  -test.v -test.run TestServiceCreateQrCodeOrder
func TestServiceCreateQrCodeOrder(t *testing.T) {
	Convey("TestServiceCreateQrCodeOrder ", t, func() {
		arg := &model.ArgCreateOrder2{
			Mid:         1,
			Month:       12,
			AppID:       1,
			AppSubID:    "11",
			OrderType:   0,
			CouponToken: "",
		}
		// ios
		r, err := s.CreateQrCodeOrder(c, arg)
		t.Logf("ios r(%+v)", r)
		t.Logf("ios PayQrCodeResp(%+v)", r.PayQrCodeResp)
		So(err, ShouldBeNil)
	})
}

//  go test  -test.v -test.run TestServiceOrderID
func TestServiceOrderID(t *testing.T) {
	Convey("TestServiceOrderID ", t, func() {
		// ios
		for i := 0; i < 100; i++ {
			t.Logf("order no (%s)", s.orderID())
		}
	})
}

//  go test  -test.v -test.run TestSignPayPlatform
func TestSignPayPlatform(t *testing.T) {
	Convey("TestSignPayPlatform ", t, func() {
		pp := make(map[string]interface{})
		pp["customerId"] = 10001
		pp["deviceType"] = 3
		sign := s.signPayPlatform(pp, s.dao.PaySign)
		t.Logf("order no (%s)", sign)
		r := &model.PayParam{
			CustomerID: 10001,
			DeviceType: 3,
		}
		str, _ := json.Marshal(r)
		pstr, _ := json.Marshal(pp)
		t.Logf("order no (%s) (%s)", str, pstr)
		So(string(str) == string(pstr), ShouldBeTrue)
	})
}

func TestServiceCreatePayParams(t *testing.T) {
	Convey("TestServiceCreatePayParams ", t, func(ctx C) {
		var (
			o = &model.PayOrder{
				OrderType: model.AutoRenew,
				Mid:       1,
			}
			p    = &model.VipPanelInfo{}
			a    = &model.ArgCreateOrder2{}
			plat int
		)
		ctx.Convey("When everything goes positive base ok", func(ctx C) {
			r, err := s.createPayParams(c, o, p, a, plat)
			fmt.Println("r[displayAccount]", r["displayAccount"])
			So(err, ShouldBeNil)
			So(r["displayAccount"] != o.Mid, ShouldBeTrue)
		})
		ctx.Convey("When everything goes positive", func(ctx C) {
			monkey.PatchInstanceMethod(reflect.TypeOf(s.memRPC), "Base", func(_ *memrpc.Service, _ context.Context, _ *memmdl.ArgMemberMid) (*memmdl.BaseInfo, error) {
				return nil, ecode.MemberNotExist
			})
			r, err := s.createPayParams(c, o, p, a, plat)
			So(err, ShouldBeNil)
			So(r["displayAccount"] == o.Mid, ShouldBeTrue)
		})
		ctx.Convey("When everything gose positive base ok PlatfromANDROIDI", func(ctx C) {
			r, err := s.createPayParams(c, o, p, a, model.PlatfromANDROIDI)
			fmt.Println("r", r)
			So(err, ShouldBeNil)
			So(r["serviceType"] == model.ServiceTypeInternational, ShouldBeTrue)
		})
		ctx.Reset(func() {
			monkey.UnpatchAll()
		})
	})
}
