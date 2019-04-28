package vip

import (
	"context"
	"fmt"
	"testing"

	"go-common/app/interface/main/account/model"

	. "github.com/smartystreets/goconvey/convey"
)

// go test  -test.v -test.run TestServiceVipPanelV2
func TestServiceVipPanelV2(t *testing.T) {
	Convey("TestServiceVipPanelV2", t, func() {
		res, err := s.VipPanelV2(context.TODO(), &model.ArgVipPanel{
			Platform: "ios",
			MobiApp:  "mobi_app",
			Device:   "phone",
			Build:    1000000,
			Mid:      476853,
		})
		t.Logf("res %+v", res)
		t.Logf("vps %+v", res.Vps)
		t.Logf("Privileges %+v", res.Privileges)
		t.Logf("TipInfo %+v", res.TipInfo)
		t.Logf("UserInfo %+v", res.UserInfo)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestServiceVipPanelV8
func TestServiceVipPanelV8(t *testing.T) {
	Convey("TestServiceVipPanelV8", t, func() {
		res, err := s.VipPanelV8(context.Background(), &model.ArgVipPanel{
			Platform: "android",
			MobiApp:  "android",
			Device:   "android",
		})
		fmt.Println("res,err", res, err)
		fmt.Println("res:", res)
		fmt.Println("res.AssociateVips:", res.AssociateVips)
		fmt.Println("res.Vps:", res.Vps)
		So(err, ShouldBeNil)
	})
}
