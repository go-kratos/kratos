package dao

import (
	"context"
	"testing"

	"go-common/app/admin/main/vip/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	ctx = context.TODO()
)

func Test_GetMonth(t *testing.T) {
	Convey("Test_GetMonth", t, func() {
		res, err := d.GetMonth(context.Background(), 11)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}
func Test_MonthList(t *testing.T) {
	Convey("Test_MonthList", t, func() {
		res, err := d.MonthList(context.TODO())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func Test_MonthEdit(t *testing.T) {
	var (
		id     int64 = 31
		status int8  = 1
		op           = "test"
	)
	Convey("Test_MonthEdit", t, func() {
		res, err := d.MonthEdit(context.Background(), id, status, op)
		So(err, ShouldBeNil)
		So(res, ShouldBeGreaterThanOrEqualTo, 0)
	})
}
func Test_GetPrice(t *testing.T) {
	var id int64 = 60
	Convey("Test_GetPrice", t, func() {
		res, err := d.GetPrice(context.Background(), id)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func Test_MonthPriceList(t *testing.T) {
	Convey("Test_MonthPriceList", t, func() {
		var (
			err error
			eff int64
			res []*model.VipMonthPrice
		)
		ap := &model.VipMonthPrice{MonthID: 2, Money: 2.00}
		eff, err = d.PriceAdd(ctx, ap)
		So(err, ShouldBeNil)
		So(eff, ShouldEqual, 1)
		res, err = d.PriceList(context.TODO(), 2)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func Test_PriceEdit(t *testing.T) {
	var (
		vp = &model.VipMonthPrice{MonthID: 2, Money: 2.00}
	)
	Convey("Test_PriceEdit", t, func() {
		res, err := d.PriceEdit(context.Background(), vp)
		So(err, ShouldBeNil)
		So(res, ShouldBeGreaterThanOrEqualTo, 0)
	})
}
