package dao

import (
	"context"
	"testing"

	"go-common/app/admin/main/manager/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestMaxBid(t *testing.T) {
	convey.Convey("maxBid", t, func() {
		_, err := d.maxBid(context.Background())
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestUpdateBusiness(t *testing.T) {
	convey.Convey("UpdateBusiness", t, func() {
		p := &model.Business{
			ID:   1,
			Name: "test",
			Flow: 1,
		}
		err := d.UpdateBusiness(context.Background(), p)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestBusinessChilds(t *testing.T) {
	convey.Convey("BusinessChilds", t, func() {
		_, err := d.BusinessChilds(context.Background(), 1)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestMaxRidByBid(t *testing.T) {
	convey.Convey("MaxRidByBid", t, func() {
		_, err := d.MaxRidByBid(context.Background(), 1)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestParentBusiness(t *testing.T) {
	convey.Convey("ParentBusiness", t, func() {
		_, err := d.ParentBusiness(context.Background(), int64(1))
		convey.So(err, convey.ShouldBeNil)
	})
}
