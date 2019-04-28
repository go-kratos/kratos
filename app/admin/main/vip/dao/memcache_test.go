package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestVipInfoKey(t *testing.T) {
	convey.Convey("vipInfoKey", t, func() {
		p1 := vipInfoKey(0)
		convey.So(p1, convey.ShouldNotBeNil)
	})
}

func TestDaoSetSelCode(t *testing.T) {
	convey.Convey("SetSelCode", t, func() {
		linkmap := map[int64]int64{1: 1}
		err := d.SetSelCode(context.TODO(), "testSelCode", linkmap)
		convey.So(err, convey.ShouldBeNil)
	})
	convey.Convey("GetSelCode", t, func() {
		_, err := d.GetSelCode(context.TODO(), "testSelCode")
		convey.So(err, convey.ShouldBeNil)
	})
	convey.Convey("DelSelCode", t, func() {
		err := d.DelSelCode(context.TODO(), "testSelCode")
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaoPingMC(t *testing.T) {
	convey.Convey("PingMC", t, func() {
		err := d.PingMC(context.TODO())
		convey.So(err, convey.ShouldBeNil)
	})
}
