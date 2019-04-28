package service

import (
	"fmt"
	"net"
	"testing"

	"go-common/app/service/main/vip/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServiceBilibiliVipGrant(t *testing.T) {
	Convey(" TestServiceBilibiliVipGrant ", t, func() {
		err := s.BilibiliVipGrant(c, &model.ArgBilibiliVipGrant{
			OpenID:     "bdca8b71e7a6726885d40a395bf9ccd1",
			OutOpenID:  "oe1822de904ceed07bpvCUOU",
			OutOrderNO: "15",
			Duration:   12,
			AppID:      32,
		})
		So(err, ShouldBeNil)
	})
}

func TestServiceCreateAssociateOrder(t *testing.T) {
	Convey(" TestServiceCreateAssociateOrder ", t, func() {
		r, err := s.CreateAssociateOrder(c, &model.ArgCreateOrder2{
			AppID:     32,
			Mid:       1,
			Platform:  "android",
			Device:    "android",
			MobiApp:   "android_i",
			PanelType: "ele",
			Month:     1,
			IP:        net.ParseIP("121.46.231.66"),
		})
		fmt.Println("r:", r)
		So(err, ShouldBeNil)
	})
}

func TestServiceEleVipGrant(t *testing.T) {
	Convey(" TestServiceEleVipGrant ", t, func() {
		err := s.EleVipGrant(c, &model.ArgEleVipGrant{
			OrderNO: "22333",
		})
		So(err, ShouldBeNil)
	})
}
