package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/service/main/point/model"
	xtime "go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoBeginTran(t *testing.T) {
	convey.Convey("BeginTran", t, func(ctx convey.C) {
		var c = context.Background()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.BeginTran(c)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoPointInfo(t *testing.T) {
	convey.Convey("PointInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.PointInfo(c, mid)
			ctx.Convey("Then err should be nil.pi should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoTxPointInfo(t *testing.T) {
	convey.Convey("TxPointInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tx, err := d.BeginTran(c)
			ctx.So(err, convey.ShouldBeNil)
			_, err = d.TxPointInfo(c, tx, mid)
			ctx.Convey("Then err should be nil.pi should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoPointHistory(t *testing.T) {
	convey.Convey("PointHistory", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(1)
			cursor = int(1)
			ps     = int(2)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.PointHistory(c, mid, cursor, ps)
			ctx.Convey("Then err should be nil.phs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoPointHistoryCount(t *testing.T) {
	convey.Convey("PointHistoryCount", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.PointHistoryCount(c, mid)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpdatePointInfo(t *testing.T) {
	convey.Convey("UpdatePointInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			pi  = &model.PointInfo{}
			ver = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tx, err := d.BeginTran(c)
			ctx.So(err, convey.ShouldBeNil)
			a, err := d.UpdatePointInfo(c, tx, pi, ver)
			ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(a, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertPoint(t *testing.T) {
	convey.Convey("InsertPoint", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			pi = &model.PointInfo{
				Mid: time.Now().Unix(),
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tx, err := d.BeginTran(c)
			ctx.So(err, convey.ShouldBeNil)
			a, err := d.InsertPoint(c, tx, pi)
			ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(a, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertPointHistory(t *testing.T) {
	convey.Convey("InsertPointHistory", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			ph = &model.PointHistory{
				Mid:        1,
				ChangeTime: xtime.Time(time.Now().Unix()),
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tx, err := d.BeginTran(c)
			ctx.So(err, convey.ShouldBeNil)
			a, err := d.InsertPointHistory(c, tx, ph)
			ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(a, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSelPointHistory(t *testing.T) {
	convey.Convey("SelPointHistory", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			mid       = int64(0)
			startDate xtime.Time
			endDate   xtime.Time
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.SelPointHistory(c, mid, startDate, endDate)
			ctx.Convey("Then err should be nil.phs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoExistPointOrder(t *testing.T) {
	convey.Convey("ExistPointOrder", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			orID = "1"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			id, err := d.ExistPointOrder(c, orID)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAllPointConfig(t *testing.T) {
	convey.Convey("AllPointConfig", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.AllPointConfig(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoOldPointHistory(t *testing.T) {
	convey.Convey("OldPointHistory", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(1)
			start = int(0)
			ps    = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.OldPointHistory(c, mid, start, ps)
			ctx.Convey("Then err should be nil.phs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
