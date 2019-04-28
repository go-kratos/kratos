package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/vip/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoCountAssociateOrder(t *testing.T) {
	convey.Convey("CountAssociateOrder", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(0)
			appID = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			count, err := d.CountAssociateOrder(c, mid, appID)
			convCtx.Convey("Then err should be nil.count should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCountAssociateGrants(t *testing.T) {
	convey.Convey("CountAssociateGrants", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(0)
			appID = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			count, err := d.CountAssociateGrants(c, mid, appID)
			convCtx.Convey("Then err should be nil.count should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCountGrantOrderByOutTradeNo(t *testing.T) {
	convey.Convey("CountGrantOrderByOutTradeNo", t, func(convCtx convey.C) {
		var (
			c          = context.Background()
			outTradeNo = ""
			appID      = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			count, err := d.CountGrantOrderByOutTradeNo(c, outTradeNo, appID)
			convCtx.Convey("Then err should be nil.count should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCountAssociateByMidAndMonths(t *testing.T) {
	convey.Convey("CountAssociateByMidAndMonths", t, func(convCtx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(0)
			appID  = int64(0)
			months = int32(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			count, err := d.CountAssociateByMidAndMonths(c, mid, appID, months)
			convCtx.Convey("Then err should be nil.count should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxInsertAssociateGrantOrder(t *testing.T) {
	convey.Convey("TxInsertAssociateGrantOrder", t, func(convCtx convey.C) {
		var (
			oa = &model.VipOrderAssociateGrant{}
			c  = context.Background()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			oldtx, err := d.OldStartTx(c)
			convCtx.So(err, convey.ShouldBeNil)
			aff, err := d.TxInsertAssociateGrantOrder(oldtx, oa)
			convCtx.Convey("Then err should be nil.aff should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(aff, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxUpdateAssociateGrantState(t *testing.T) {
	convey.Convey("TxUpdateAssociateGrantState", t, func(convCtx convey.C) {
		var (
			oa = &model.VipOrderAssociateGrant{}
			c  = context.Background()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			oldtx, err := d.OldStartTx(c)
			convCtx.So(err, convey.ShouldBeNil)
			aff, err := d.TxUpdateAssociateGrantState(oldtx, oa)
			convCtx.Convey("Then err should be nil.aff should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(aff, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAssociateGrantCountInfo(t *testing.T) {
	convey.Convey("AssociateGrantCountInfo", t, func(convCtx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(0)
			appID  = int64(0)
			months = int32(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.AssociateGrantCountInfo(c, mid, appID, months)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddAssociateGrantCount(t *testing.T) {
	convey.Convey("AddAssociateGrantCount", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			arg = &model.VipAssociateGrantCount{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.AddAssociateGrantCount(c, arg)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpdateAssociateGrantCount(t *testing.T) {
	convey.Convey("UpdateAssociateGrantCount", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			arg = &model.VipAssociateGrantCount{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.UpdateAssociateGrantCount(c, arg)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
