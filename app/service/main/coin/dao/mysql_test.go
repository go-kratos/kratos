package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaomidHit(t *testing.T) {
	var (
		mid = int64(0)
	)
	convey.Convey("midHit", t, func(ctx convey.C) {
		p1 := d.midHit(mid)
		ctx.Convey("p1 should not be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoaidHit(t *testing.T) {
	var (
		aid = int64(0)
	)
	convey.Convey("aidHit", t, func(ctx convey.C) {
		p1 := d.aidHit(aid)
		ctx.Convey("p1 should not be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaohitCoinPeriod(t *testing.T) {
	var (
		c      = context.TODO()
		now, _ = time.Parse("2006-01-02 15:04:05", "2018-01-26 00:00:00")
	)
	convey.Convey("hitCoinPeriod", t, func(ctx convey.C) {
		id, err := d.hitCoinPeriod(c, now)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("id should not be nil", func(ctx convey.C) {
			ctx.So(id, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpdateCoinSettleBD(t *testing.T) {
	var (
		c        = context.TODO()
		aid      = int64(4051951)
		tp       = int64(1)
		expSub   = int64(20)
		describe = "UpdateCoinSettleBD"
		now      = time.Now()
	)
	convey.Convey("UpdateCoinSettleBD", t, func(ctx convey.C) {
		currentYear, currentMonth, _ := now.Date()
		curTime := now.Location()
		firstDay := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, curTime)
		now = firstDay.AddDate(0, 1, -1)
		affect, err := d.UpdateCoinSettleBD(c, aid, tp, expSub, describe, now)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("affect should not be nil", func(ctx convey.C) {
			ctx.So(affect, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoCoinList(t *testing.T) {
	var (
		c    = context.TODO()
		mid  = int64(4051951)
		tp   = int64(1)
		ts   = int64(time.Now().Unix())
		size = int64(1024)
	)
	convey.Convey("CoinList", t, func(ctx convey.C) {
		_, err := d.CoinList(c, mid, tp, ts, size)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		// ctx.Convey("rcs should not be nil", func(ctx convey.C) {
		// 	ctx.So(rcs, convey.ShouldNotBeEmpty)
		// })
	})
}

func TestDaoCoinsAddedByMid(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(4051951)
		aid = int64(12)
		tp  = int64(20)
	)
	convey.Convey("CoinsAddedByMid", t, func(ctx convey.C) {
		added, err := d.CoinsAddedByMid(c, mid, aid, tp)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("added should not be nil", func(ctx convey.C) {
			ctx.So(added, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddedCoins(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(4051951)
		upMid = int64(4051950)
	)
	convey.Convey("AddedCoins", t, func(ctx convey.C) {
		added, err := d.AddedCoins(c, mid, upMid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("added should not be nil", func(ctx convey.C) {
			ctx.So(added, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUserCoinsAdded(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(1)
	)
	convey.Convey("UserCoinsAdded", t, func(ctx convey.C) {
		addeds, err := d.UserCoinsAdded(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("addeds should not be nil", func(ctx convey.C) {
			ctx.So(addeds, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoInsertCoinArchive(t *testing.T) {
	var (
		c         = context.TODO()
		aid       = int64(1)
		tp        = int64(21)
		mid       = int64(4051951)
		timestamp = int64(time.Now().Unix())
		multiply  = int64(0)
	)
	convey.Convey("InsertCoinArchive", t, func(ctx convey.C) {
		err := d.InsertCoinArchive(c, aid, tp, mid, timestamp, multiply)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoInsertCoinMember(t *testing.T) {
	var (
		c         = context.TODO()
		aid       = int64(1)
		tp        = int64(2)
		mid       = int64(4051951)
		timestamp = int64(time.Now().Unix())
		multiply  = int64(21)
		upMid     = int64(1)
	)
	convey.Convey("InsertCoinMember", t, func(ctx convey.C) {
		err := d.InsertCoinMember(c, aid, tp, mid, timestamp, multiply, upMid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoUpdateItemCoinCount(t *testing.T) {
	var (
		c     = context.TODO()
		aid   = int64(0)
		tp    = int64(0)
		count = int64(0)
	)
	convey.Convey("UpdateItemCoinCount", t, func(ctx convey.C) {
		err := d.UpdateItemCoinCount(c, aid, tp, count)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoRawItemCoin(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(0)
		tp  = int64(0)
	)
	convey.Convey("RawItemCoin", t, func(ctx convey.C) {
		count, err := d.RawItemCoin(c, aid, tp)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("count should not be nil", func(ctx convey.C) {
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpdateCoinMemberCount(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(4051951)
		upMid = int64(1)
		count = int64(2)
	)
	convey.Convey("UpdateCoinMemberCount", t, func(ctx convey.C) {
		err := d.UpdateCoinMemberCount(c, mid, upMid, count)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoBeginTran(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("BeginTran", t, func(ctx convey.C) {
		no, err := d.BeginTran(c)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("no should not be nil", func(ctx convey.C) {
			ctx.So(no, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpdateItemCoins(t *testing.T) {
	var (
		c     = context.TODO()
		aid   = int64(1)
		tp    = int64(20)
		coins = int64(20)
		now   = time.Now()
	)
	convey.Convey("UpdateItemCoins", t, func(ctx convey.C) {
		affect, err := d.UpdateItemCoins(c, aid, tp, coins, now)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("affect should not be nil", func(ctx convey.C) {
			ctx.So(affect, convey.ShouldNotBeNil)
		})
	})
}
