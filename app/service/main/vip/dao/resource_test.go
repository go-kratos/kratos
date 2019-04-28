package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/service/main/vip/model"
	xtime "go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSelNewResourcePool(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(0)
	)
	convey.Convey("SelNewResourcePool", t, func(ctx convey.C) {
		_, err := d.SelNewResourcePool(c, id)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSelNewBusiness(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(0)
	)
	convey.Convey("SelNewBusiness", t, func(ctx convey.C) {
		_, err := d.SelNewBusiness(c, id)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSelNewBusinessByAppkey(t *testing.T) {
	var (
		c      = context.TODO()
		appkey = ""
	)
	convey.Convey("SelNewBusinessByAppkey", t, func(ctx convey.C) {
		_, err := d.SelNewBusinessByAppkey(c, appkey)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSelCode(t *testing.T) {
	var (
		c       = context.TODO()
		codeStr = ""
	)
	convey.Convey("SelCode", t, func(ctx convey.C) {
		code, err := d.SelCode(c, codeStr)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("code should not be nil", func(ctx convey.C) {
			ctx.So(code, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSelCodes(t *testing.T) {
	var (
		c     = context.TODO()
		codes []string
	)
	convey.Convey("SelCodes", t, func(ctx convey.C) {
		_, err := d.SelCodes(c, codes)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSelBatchCode(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(0)
	)
	convey.Convey("SelBatchCode", t, func(ctx convey.C) {
		_, err := d.SelBatchCode(c, id)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSelBatchCodes(t *testing.T) {
	var (
		c   = context.TODO()
		ids []int64
	)
	convey.Convey("SelBatchCodes", t, func(ctx convey.C) {
		_, err := d.SelBatchCodes(c, ids)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSelBatchCodesByBisID(t *testing.T) {
	var (
		c     = context.TODO()
		bisID = int64(0)
	)
	convey.Convey("SelBatchCodesByBisID", t, func(ctx convey.C) {
		_, err := d.SelBatchCodesByBisID(c, bisID)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSelCodeOpened(t *testing.T) {
	var (
		c      = context.TODO()
		bisIDs = []int64{1}
		arg    = &model.ArgCodeOpened{}
	)
	convey.Convey("SelCodeOpened", t, func(ctx convey.C) {

		_, err := d.SelCodeOpened(c, bisIDs, arg)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoTxUpdateCode(t *testing.T) {
	var (
		c       = context.TODO()
		id      = int64(0)
		mid     = int64(0)
		useTime = xtime.Time(time.Now().Unix())
	)
	convey.Convey("TxUpdateCode", t, func(ctx convey.C) {
		tx, err := d.StartTx(c)
		ctx.Convey("Tx Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		defer tx.Commit()
		eff, err := d.TxUpdateCode(tx, id, mid, useTime)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("eff should not be nil", func(ctx convey.C) {
			ctx.So(eff, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTxUpdateCodeStatus(t *testing.T) {
	var (
		c      = context.TODO()
		id     = int64(0)
		status = int8(0)
	)
	convey.Convey("TxUpdateCodeStatus", t, func(ctx convey.C) {
		tx, err := d.StartTx(c)
		ctx.Convey("Tx Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		defer tx.Commit()
		eff, err := d.TxUpdateCodeStatus(tx, id, status)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("eff should not be nil", func(ctx convey.C) {
			ctx.So(eff, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTxUpdateBatchCode(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(0)
		sc = int32(0)
	)
	convey.Convey("TxUpdateBatchCode", t, func(ctx convey.C) {
		tx, err := d.StartTx(c)
		ctx.Convey("Tx Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		defer tx.Commit()
		eff, err := d.TxUpdateBatchCode(tx, id, sc)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("eff should not be nil", func(ctx convey.C) {
			ctx.So(eff, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSelCodesByBMid(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
	)
	convey.Convey("SelCodesByBMid", t, func(ctx convey.C) {
		cs, err := d.SelCodesByBMid(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("cs should not be nil", func(ctx convey.C) {
			ctx.So(cs, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSelActives(t *testing.T) {
	var (
		c         = context.TODO()
		relations []string
	)
	convey.Convey("SelActives", t, func(ctx convey.C) {
		_, err := d.SelActives(c, relations)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSelBatchCount(t *testing.T) {
	var (
		c                 = context.TODO()
		batchCodeID int64 = 1
		mid         int64 = 1001
	)
	convey.Convey("sel batch count", t, func(ctx convey.C) {
		_, err := d.SelBatchCount(c, batchCodeID, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
