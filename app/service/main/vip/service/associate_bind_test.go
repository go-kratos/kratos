package service

import (
	"fmt"
	"testing"

	"go-common/app/service/main/vip/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServiceaddOpenBind(t *testing.T) {
	Convey(" TestServiceaddOpenBind ", t, func() {
		err := s.addOpenBind(c, 1, "xxxx", 2, &model.OpenBindInfo{}, &model.OpenBindInfo{})
		So(err, ShouldBeNil)
	})
}

func TestServiceOpenBindByOutOpenID(t *testing.T) {
	Convey(" TestServiceOpenBindByOutOpenID ", t, func() {
		err := s.OpenBindByOutOpenID(c, &model.ArgBind{
			AppID:     32,
			OpenID:    "bdca8b71e7a6726885d40a395bf9ccd1",
			OutOpenID: "7a6726885d40a395bf9ccd2",
		})
		So(err, ShouldBeNil)
	})
}

func TestServiceOpenBindByMid(t *testing.T) {
	Convey(" TestServiceOpenBindByMid ", t, func() {
		err := s.OpenBindByMid(c, &model.ArgOpenBindByMid{
			AppID:     32,
			OutOpenID: "7a6726885d40a395bf9ccd3",
			Mid:       1,
		})
		So(err, ShouldBeNil)
	})
}

func TestServiceBindInfoByMid(t *testing.T) {
	Convey(" TestServiceBindInfoByMid ", t, func() {
		res, err := s.BindInfoByMid(c, &model.ArgBindInfo{
			AppID: 30,
			Mid:   1,
		})
		So(err, ShouldBeNil)
		fmt.Println("res", res.Account, res.Outer)
		So(res, ShouldNotBeNil)
	})
}
