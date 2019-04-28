package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/vip/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoOpenInfoOpenID(t *testing.T) {
	convey.Convey("OpenInfoOpenID", t, func(convCtx convey.C) {
		var (
			c      = context.Background()
			openID = ""
			appID  = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.RawOpenInfoByOpenID(c, openID, appID)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoOpenInfoByMid(t *testing.T) {
	convey.Convey("OpenInfoByMid", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(0)
			appID = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.OpenInfoByMid(c, mid, appID)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoByOutOpenID(t *testing.T) {
	convey.Convey("ByOutOpenID", t, func(convCtx convey.C) {
		var (
			c         = context.Background()
			outOpenID = ""
			appID     = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.ByOutOpenID(c, outOpenID, appID)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRawBindInfoByMid(t *testing.T) {
	convey.Convey("RawBindInfoByMid", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(0)
			appID = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.RawBindInfoByMid(c, mid, appID)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxAddBind(t *testing.T) {
	convey.Convey("TxAddBind", t, func(convCtx convey.C) {
		var (
			tx, _ = d.BeginTran(context.Background())
			arg   = &model.OpenBindInfo{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.TxAddBind(tx, arg)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
			tx.Commit()
		})
	})
}

func TestDaoTxUpdateBindByOutOpenID(t *testing.T) {
	convey.Convey("TxUpdateBindByOutOpenID", t, func(convCtx convey.C) {
		var (
			tx, _ = d.BeginTran(context.Background())
			arg   = &model.OpenBindInfo{
				AppID:     32,
				OutOpenID: "test-2333",
			}
		)
		convCtx.Convey("When TxAddBind", func(convCtx convey.C) {
			err := d.TxAddBind(tx, arg)
			convCtx.So(err, convey.ShouldBeNil)
			tx.Commit()
		})
		convCtx.Convey("When TxUpdateBindByOutOpenID", func(convCtx convey.C) {
			err := d.TxUpdateBindByOutOpenID(tx, arg)
			convCtx.So(err, convey.ShouldBeNil)
			tx.Commit()
		})
		convCtx.Convey("When TxDeleteBindInfo", func(convCtx convey.C) {
			err := d.TxDeleteBindInfo(tx, arg.ID)
			convCtx.So(err, convey.ShouldBeNil)
			tx.Commit()
		})
	})
}

func TestDaoUpdateBindState(t *testing.T) {
	convey.Convey("UpdateBindState", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			arg = &model.OpenBindInfo{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.UpdateBindState(c, arg)
			convCtx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAddOpenInfo(t *testing.T) {
	convey.Convey("AddOpenInfo", t, func(convCtx convey.C) {
		var (
			c = context.Background()
			a = &model.OpenInfo{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.AddOpenInfo(c, a)
			convCtx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoBindInfoByOutOpenIDAndMid(t *testing.T) {
	convey.Convey("BindInfoByOutOpenIDAndMid", t, func(convCtx convey.C) {
		var (
			c         = context.Background()
			mid       = int64(0)
			outOpenID = ""
			appID     = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.BindInfoByOutOpenIDAndMid(c, mid, outOpenID, appID)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
