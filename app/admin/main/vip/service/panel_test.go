package service

import (
	"context"
	"fmt"
	"go-common/app/admin/main/vip/model"
	xtime "go-common/library/time"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServiceVipPriceConfigs(t *testing.T) {
	Convey("Test_VipPriceConfigs", t, func() {
		av := &model.ArgVipPrice{
			Plat:     1,
			Month:    1,
			SubType:  1,
			SuitType: 1,
		}
		res, err := s.VipPriceConfigs(context.TODO(), av)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestServiceVipPriceConfigID(t *testing.T) {
	Convey("Test_VipPriceConfigID", t, func() {
		res, err := s.VipPriceConfigID(context.TODO(), &model.ArgVipPriceID{ID: 1})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestServiceAddVipPriceConfig(t *testing.T) {
	Convey("Test_AddVipPriceConfig", t, func() {
		aavpc := &model.ArgAddOrUpVipPrice{
			Plat:     1,
			PdName:   "xxx",
			PdID:     "xxx",
			Month:    1,
			SubType:  1,
			SuitType: 0,
			OPrice:   1,
		}
		err := s.AddVipPriceConfig(context.TODO(), aavpc)
		So(err, ShouldBeNil)
	})
}

func TestServiceUpVipPriceConfig(t *testing.T) {
	Convey("Test_UpVipPriceConfig", t, func() {
		aavpc := &model.ArgAddOrUpVipPrice{
			ID:         235,
			Plat:       1,
			PdName:     "xxx",
			PdID:       "xxx",
			Month:      1,
			SubType:    1,
			SuitType:   1,
			OPrice:     1,
			StartBuild: 0,
			EndBuild:   103,
		}
		err := s.UpVipPriceConfig(context.TODO(), aavpc)
		So(err, ShouldBeNil)
	})
}

func TestServiceDelVipPriceConfig(t *testing.T) {
	advp := &model.ArgVipPriceID{ID: 1}
	Convey("Test_DelVipPriceConfig", t, func() {
		err := s.DelVipPriceConfig(context.TODO(), advp)
		So(err, ShouldBeNil)
	})
}

func TestServiceVipDPriceConfigs(t *testing.T) {
	Convey("Test_VipDPriceConfigs", t, func() {
		av := &model.ArgVipPriceID{
			ID: 155,
		}
		res, err := s.VipDPriceConfigs(context.TODO(), av)
		fmt.Println("res", res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestServiceVipDPriceConfigID(t *testing.T) {
	Convey("Test_VipDPriceConfigID", t, func() {
		res, err := s.VipDPriceConfigID(context.TODO(), &model.ArgVipDPriceID{DisID: 10})
		fmt.Println("res", res.FirstPrice)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestServiceAddVipDPriceConfig(t *testing.T) {
	Convey("Test_AddVipDPriceConfig", t, func() {
		aavpc := &model.ArgAddOrUpVipDPrice{
			DisID:      1,
			ID:         155,
			PdID:       "tv.danmaku.bilianimexAuto3VIP",
			DPrice:     15,
			STime:      xtime.Time(time.Now().Unix()) + 100000000000,
			ETime:      xtime.Time(time.Now().Unix()) + 200000000000,
			Remark:     "test2",
			Operator:   "admin",
			OpID:       11,
			FirstPrice: 13,
		}
		err := s.AddVipDPriceConfig(context.TODO(), aavpc)
		So(err, ShouldBeNil)
	})
}

func TestServiceUpVipDPriceConfig(t *testing.T) {
	Convey("Test_UpVipDPriceConfig", t, func() {
		aavpc := &model.ArgAddOrUpVipDPrice{
			DisID:      11,
			ID:         155,
			PdID:       "tv.danmaku.bilianimexAuto3VIP",
			DPrice:     1.1,
			STime:      xtime.Time(time.Now().Unix()) + 100000000000,
			ETime:      xtime.Time(time.Now().Unix()) + 200000000000,
			Remark:     "test2",
			Operator:   "admin",
			OpID:       11,
			FirstPrice: 11,
		}
		err := s.UpVipDPriceConfig(context.TODO(), aavpc)
		So(err, ShouldBeNil)
	})
}

func TestServiceDelVipDPriceConfig(t *testing.T) {
	advp := &model.ArgVipDPriceID{DisID: 1}
	Convey("Test_DelVipDPriceConfig", t, func() {
		err := s.DelVipDPriceConfig(context.TODO(), advp)
		So(err, ShouldBeNil)
	})
}
