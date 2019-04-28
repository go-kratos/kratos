package dao

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go-common/app/admin/main/vip/model"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDaoAddVipPriceConfig(t *testing.T) {
	Convey("Test_AddVipPriceConfig", t, func() {
		v := &model.VipPriceConfig{
			Plat:     1,
			PdName:   "ddd",
			PdID:     "sss",
			SuitType: 1,
			Month:    1,
			SubType:  2,
			OPrice:   1,
			Selected: 1,
			Remark:   "",
			Status:   0,
			Operator: "sss",
			OpID:     1,
		}
		err := d.AddVipPriceConfig(context.TODO(), v)
		So(err, ShouldBeNil)
	})
	Convey("Test_AddVipPriceConfig add build 200-300", t, func() {
		v := &model.VipPriceConfig{
			Plat:       10,
			PdName:     "ddd",
			PdID:       "sss",
			SuitType:   1,
			Month:      1,
			SubType:    1,
			OPrice:     1,
			Selected:   1,
			Remark:     "",
			Status:     0,
			Operator:   "sss",
			OpID:       1,
			StartBuild: 200,
			EndBuild:   300,
		}
		err := d.AddVipPriceConfig(context.TODO(), v)
		So(err, ShouldBeNil)
	})
	Convey("Test_AddVipPriceConfig add build 0-0", t, func() {
		v := &model.VipPriceConfig{
			Plat:       10,
			PdName:     "ddd",
			PdID:       "sss",
			SuitType:   1,
			Month:      1,
			SubType:    1,
			OPrice:     1,
			Selected:   1,
			Remark:     "",
			Status:     0,
			Operator:   "sss",
			OpID:       1,
			StartBuild: 200,
			EndBuild:   300,
		}
		err := d.AddVipPriceConfig(context.TODO(), v)
		So(err, ShouldBeNil)
	})
}

func TestDaoAddVipDPriceConfig(t *testing.T) {
	Convey("Test_AddVipPriceConfig", t, func() {
		v := &model.VipDPriceConfig{
			DisID:    1,
			ID:       1,
			PdID:     "sss",
			DPrice:   1,
			STime:    xtime.Time(time.Now().Unix()),
			Remark:   "",
			Operator: "sss",
			OpID:     1,
		}
		err := d.AddVipDPriceConfig(context.TODO(), v)
		So(err, ShouldBeNil)
	})
}

func TestUpVipPriceConfig(t *testing.T) {
	Convey("TestUpVipPriceConfig", t, func() {
		err := d.UpVipPriceConfig(context.TODO(), &model.VipPriceConfig{
			ID:         230,
			Plat:       10,
			PdName:     "ddd",
			PdID:       "sss",
			SuitType:   1,
			Month:      1,
			SubType:    1,
			OPrice:     1,
			Selected:   1,
			Remark:     "",
			Status:     0,
			Operator:   "sss",
			OpID:       1,
			StartBuild: 250,
			EndBuild:   300,
		})
		So(err, ShouldBeNil)
	})
}

func TestDaoUpVipPriceConfig(t *testing.T) {
	Convey("Test_UpVipPriceConfig", t, func() {
		v := &model.VipPriceConfig{
			ID:       1,
			Plat:     1,
			PdName:   "ddd",
			PdID:     "sss",
			SuitType: 1,
			Month:    1,
			SubType:  2,
			OPrice:   1,
			Selected: 1,
			Remark:   "",
			Status:   0,
			Operator: "sss",
			OpID:     1,
		}
		err := d.AddVipPriceConfig(context.TODO(), v)
		So(err, ShouldBeNil)
	})
}

func TestDaoUpVipDPriceConfig(t *testing.T) {
	Convey("Test_UpVipDPriceConfig", t, func() {
		v := &model.VipDPriceConfig{
			DisID:    1,
			ID:       1,
			PdID:     "sss",
			DPrice:   1,
			STime:    xtime.Time(time.Now().Unix()),
			Remark:   "",
			Operator: "sss",
			OpID:     1,
		}
		err := d.UpVipDPriceConfig(context.TODO(), v)
		So(err, ShouldBeNil)
	})
}

func TestDaoDelVipPriceConfig(t *testing.T) {
	Convey("Test_DelVipPriceConfig", t, func() {
		err := d.DelVipPriceConfig(context.TODO(), &model.ArgVipPriceID{ID: 1})
		So(err, ShouldBeNil)
	})
}
func TestDaoDelVipDPriceConfig(t *testing.T) {
	Convey("Test_DelVipDPriceConfig", t, func() {
		err := d.DelVipDPriceConfig(context.TODO(), &model.ArgVipDPriceID{DisID: 1})
		So(err, ShouldBeNil)
	})
}

func TestDaoVipPriceConfigUQCheck(t *testing.T) {
	Convey("Test_VipPriceConfigUQCheck", t, func() {
		res, err := d.VipPriceConfigUQCheck(context.TODO(), &model.ArgAddOrUpVipPrice{Plat: 1, Month: 1, SubType: 1, SuitType: 0})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
	Convey("Test_VipPriceConfigUQCheck count 0-0", t, func() {
		count, err := d.VipPriceConfigUQCheck(context.TODO(), &model.ArgAddOrUpVipPrice{Plat: 10, Month: 1, SubType: 1, SuitType: 1})
		fmt.Println("count:", count)
		So(err, ShouldBeNil)
		So(count, ShouldNotBeNil)
	})
	Convey("Test_VipPriceConfigUQCheck count 220-0", t, func() {
		count, err := d.VipPriceConfigUQCheck(context.TODO(), &model.ArgAddOrUpVipPrice{Plat: 10, Month: 1, SubType: 1, SuitType: 1, StartBuild: 220, EndBuild: 0})
		fmt.Println("count:", count)
		So(err, ShouldBeNil)
		So(count, ShouldNotBeNil)
	})
	Convey("Test_VipPriceConfigUQCheck count 120-180", t, func() {
		count, err := d.VipPriceConfigUQCheck(context.TODO(), &model.ArgAddOrUpVipPrice{Plat: 10, Month: 1, SubType: 1, SuitType: 1, StartBuild: 120, EndBuild: 180})
		fmt.Println("count:", count)
		So(err, ShouldBeNil)
		So(count, ShouldNotBeNil)
	})
	Convey("Test_VipPriceConfigUQCheck count 0-180", t, func() {
		count, err := d.VipPriceConfigUQCheck(context.TODO(), &model.ArgAddOrUpVipPrice{Plat: 10, Month: 1, SubType: 1, SuitType: 1, StartBuild: 0, EndBuild: 180})
		fmt.Println("count:", count)
		So(err, ShouldBeNil)
		So(count, ShouldNotBeNil)
	})
	Convey("Test_VipPriceConfigUQCheck count 0-600", t, func() {
		count, err := d.VipPriceConfigUQCheck(context.TODO(), &model.ArgAddOrUpVipPrice{Plat: 10, Month: 1, SubType: 1, SuitType: 1, StartBuild: 0, EndBuild: 600})
		fmt.Println("count:", count)
		So(err, ShouldBeNil)
		So(count, ShouldNotBeNil)
	})
}

func TestDaoVipPriceConfigs(t *testing.T) {
	Convey("Test_VipPriceConfigs", t, func() {
		count, err := d.VipPriceConfigs(context.TODO())
		So(err, ShouldBeNil)
		So(count, ShouldNotBeEmpty)
	})
}

func TestDaoVipPriceConfigID(t *testing.T) {
	Convey("Test_VipPriceConfigID", t, func() {
		ids, err := d.VipPriceConfigID(context.TODO(), &model.ArgVipPriceID{ID: 1})
		So(err, ShouldBeNil)
		So(ids, ShouldNotBeEmpty)
	})
}

func TestDaoVipDPriceConfigs(t *testing.T) {
	Convey("Test_VipDPriceConfigs", t, func() {
		count, err := d.VipDPriceConfigs(context.TODO(), &model.ArgVipPriceID{ID: 1})
		So(err, ShouldBeNil)
		So(count, ShouldNotBeEmpty)
	})
}

func TestDaoVipDPriceConfigID(t *testing.T) {
	Convey("Test_VipDPriceConfigID", t, func() {
		res, err := d.VipDPriceConfigID(context.TODO(), &model.ArgVipDPriceID{DisID: 1})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestDaoVipDPriceConfigUQTime(t *testing.T) {
	Convey("Test_VipDPriceConfigUQTime", t, func() {
		arg := &model.ArgAddOrUpVipDPrice{
			ID:    992,
			STime: xtime.Time(time.Now().Unix()),
			ETime: xtime.Time(time.Now().Unix()),
		}
		_, err := d.VipDPriceConfigUQTime(context.TODO(), arg)
		So(err, ShouldBeNil)
	})
}

func TestDaoVipPriceDiscountConfigs(t *testing.T) {
	Convey("Test_VipPriceDiscountConfigs", t, func() {
		res, err := d.VipPriceDiscountConfigs(context.TODO())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestDaoVipMaxPriceDiscount(t *testing.T) {
	Convey("Test_VipMaxPriceDiscount", t, func() {
		res, err := d.VipMaxPriceDiscount(context.TODO(), &model.ArgAddOrUpVipPrice{ID: 1})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestDaoCountVipPriceConfigByPlat(t *testing.T) {
	Convey("TestDaoCountVipPriceConfigByPlat", t, func() {
		res, err := d.CountVipPriceConfigByPlat(context.TODO(), 1)
		fmt.Println(res)
		So(err, ShouldBeNil)
		So(res, ShouldBeGreaterThanOrEqualTo, 0)
	})
}
