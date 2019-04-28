package service

import (
	"fmt"
	"testing"

	v1 "go-common/app/service/main/vip/api"
	"go-common/app/service/main/vip/model"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_loadVipPriceConfig(t *testing.T) {
	Convey("loadVipPriceConfig test ", t, func() {
		s.loadVipPriceConfig()
	})
}
func Test_VipUserPanel(t *testing.T) {
	Convey("VipUserPanel test ", t, func() {
		mvps, err := s.VipUserPanel(c, 1, 1, 1, 0)
		So(err, ShouldBeNil)
		So(mvps, ShouldNotBeNil)
	})
}

//  go test  -test.v -test.run TestServiceVipPanelExplain

func TestServiceVipPanelExplain(t *testing.T) {
	Convey("TestServiceVipPanelExplain test ", t, func() {
		res, err := s.VipPanelExplain(c, &model.ArgPanelExplain{
			Mid: 9,
		})
		t.Logf("res(%+v)", res)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestVipUserPanel5
func TestVipUserPanel5(t *testing.T) {
	Convey("VipUserPanel5 test ", t, func() {
		mvps, err := s.VipUserPanelV5(c, &model.ArgPanel{
			Mid:       150,
			Platform:  "ios",
			Device:    "phone",
			PanelType: "normal",
			MobiApp:   "iphone",
		})
		fmt.Println("vps", mvps.Vps)
		fmt.Println("vps dprice", mvps.Vps[0].DPrice)
		So(err, ShouldBeNil)
		So(mvps, ShouldNotBeNil)
	})
}

func TestPriceByProductID(t *testing.T) {
	Convey("PriceByProductID test ", t, func() {
		mvps, err := s.PriceByProductID(c, "tv.danmaku.bilianimexAuto1VIP")
		fmt.Println("PriceByProductID", mvps)
		fmt.Println("PriceByProductID", mvps.DPrice)
		So(err, ShouldBeNil)
		So(mvps, ShouldNotBeNil)
	})
}

func TestMakeVipConfig(t *testing.T) {
	Convey("TestMakeVipConfig test ", t, func() {
		vmap := map[int8]map[int16]*model.VipPriceConfig{}
		vpc := &model.VipPriceConfig{
			SuitType: 1,
			SubType:  0,
			Month:    1,
		}
		s.makeVipConfig(vmap, vpc, vpc.SuitType, 600)
		So(len(vmap[0]) == 1, ShouldBeTrue)
		vpc = &model.VipPriceConfig{
			SuitType:   1,
			SubType:    0,
			Month:      2,
			StartBuild: 100,
			EndBuild:   300,
		}
		s.makeVipConfig(vmap, vpc, vpc.SuitType, 600)
		So(len(vmap[0]) == 1, ShouldBeTrue)
		s.makeVipConfig(vmap, vpc, vpc.SuitType, 200)
		So(len(vmap[0]) == 2, ShouldBeTrue)
	})
}

func TestServiceVipUserPanelV9(t *testing.T) {
	Convey("VipUserPanelV9 ", t, func() {
		res, err := s.VipUserPanelV9(c, &v1.VipUserPanelReq{
			Mid:       1,
			Platform:  "pc",
			PanelType: "normal",
		})
		fmt.Println("res", res.Coupon)
		So(err, ShouldBeNil)
	})
}
