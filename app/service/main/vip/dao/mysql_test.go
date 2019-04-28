package dao

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"go-common/app/service/main/vip/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoOldInsertVipBcoinSalary(t *testing.T) {
	var (
		c = context.TODO()
		r = &model.VipBcoinSalary{}
	)
	r.Amount = 10
	r.Status = 0
	r.Mid = 1
	r.Memo = "年费会员发放B币"
	r.GiveNowStatus = 0
	convey.Convey("OldInsertVipBcoinSalary,should return true where err == nil", t, func(ctx convey.C) {
		err := d.OldInsertVipBcoinSalary(c, r)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSelAllConfig(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("SelAllConfig, should return true where err == nil and res not empty", t, func(ctx convey.C) {
		res, err := d.SelAllConfig(c)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("res should not be nil", func(ctx convey.C) {
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSelVipAppInfo(t *testing.T) {
	var (
		c  = context.TODO()
		no = int(0)
	)
	convey.Convey("SelVipAppInfo", t, func(ctx convey.C) {
		_, err := d.SelVipAppInfo(c, no)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoOldSelLastBcoin(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(1)
	)
	convey.Convey("OldSelLastBcoin,should return true where err == nil and r != nil", t, func(ctx convey.C) {
		r, err := d.OldSelLastBcoin(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("r should not be nil", func(ctx convey.C) {
			ctx.So(r, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSelBusiness(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(4)
	)
	convey.Convey("SelBusiness,should return true where err == nil and r != nil", t, func(ctx convey.C) {
		_, err := d.SelBusiness(c, id)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSelResourcePool(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(5)
	)
	convey.Convey("SelResourcePool", t, func(ctx convey.C) {
		_, err := d.SelResourcePool(c, id)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoOldSelVipUserInfo(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(1)
	)
	convey.Convey("OldSelVipUserInfo", t, func(ctx convey.C) {
		_, err := d.OldSelVipUserInfo(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSelVipResourceBatch(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(5)
	)
	convey.Convey("SelVipResourceBatch,should return true where err == nil and r != nil", t, func(ctx convey.C) {
		_, err := d.SelVipResourceBatch(c, id)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoOldSelVipChangeHistory(t *testing.T) {
	var (
		c                = context.TODO()
		relationID       = "ad4bb9b8f5d9d4a7123459780"
		batchID    int64 = 1
	)
	convey.Convey("OldSelVipChangeHistory", t, func(ctx convey.C) {
		_, err := d.OldVipchangeHistory(c, relationID, batchID)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoOldUpdateVipUserInfo(t *testing.T) {
	var (
		c = context.TODO()
		r = &model.VipUserInfo{Mid: rand.Int63()}
	)
	convey.Convey("OldUpdateVipUserInfo", t, func(ctx convey.C) {
		tx, err := d.StartTx(c)
		ctx.Convey("Tx Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		defer tx.Commit()
		a, err := d.OldUpdateVipUserInfo(c, tx, r)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("a should not be nil", func(ctx convey.C) {
			ctx.So(a, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpdateBatchCount3(t *testing.T) {
	var (
		c = context.TODO()

		r   = &model.VipResourceBatch{ID: 1}
		ver = int64(0)
	)
	convey.Convey("UpdateBatchCount", t, func(ctx convey.C) {
		tx, err := d.StartTx(c)
		ctx.So(err, convey.ShouldBeNil)
		defer tx.Commit()
		a, err := d.UpdateBatchCount(c, tx, r, ver)
		ctx.So(err, convey.ShouldBeNil)
		ctx.So(a, convey.ShouldBeGreaterThan, 0)

	})
}

func TestDaoOldInsertVipUserInfo(t *testing.T) {
	var (
		c = context.TODO()
		r = &model.VipUserInfo{Mid: time.Now().Unix()}
	)
	convey.Convey("OldInsertVipUserInfo", t, func(ctx convey.C) {
		tx, err := d.OldStartTx(c)
		ctx.So(err, convey.ShouldBeNil)
		err = d.OldInsertVipUserInfo(c, tx, r)
		ctx.So(err, convey.ShouldBeNil)
	})
}

func TestDaoOldInsertVipChangeHistory(t *testing.T) {
	var (
		c = context.TODO()
		r = &model.VipChangeHistory{}
	)
	convey.Convey("OldInsertVipChangeHistory", t, func(ctx convey.C) {
		tx, err := d.OldStartTx(c)
		ctx.Convey("Tx Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		id, err := d.OldInsertVipChangeHistory(c, tx, r)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("id should not be nil", func(ctx convey.C) {
			ctx.So(id, convey.ShouldNotBeNil)
		})
	})
}

func TestDao_OldTxDelBcoinSalary(t *testing.T) {
	convey.Convey("old tx del bcoin salary", t, func() {
		tx, _ := d.OldStartTx(context.TODO())
		err := d.OldTxDelBcoinSalary(tx, 106138791, time.Now().AddDate(0, 1, 0))
		tx.Commit()
		convey.So(err, convey.ShouldBeNil)
	})
}
