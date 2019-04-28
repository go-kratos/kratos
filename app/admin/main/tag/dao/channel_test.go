package dao

import (
	"context"
	"go-common/app/admin/main/tag/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaosetChannels(t *testing.T) {
	var (
		channels = []*model.Channel{}
	)
	convey.Convey("setChannels", t, func(ctx convey.C) {
		p1 := setChannels(channels)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoChannel(t *testing.T) {
	var (
		c   = context.TODO()
		tid = int64(0)
	)
	convey.Convey("Channel", t, func(ctx convey.C) {
		res, err := d.Channel(c, tid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestDaoChannels(t *testing.T) {
	var (
		c    = context.TODO()
		tids = []int64{0}
	)
	convey.Convey("Channels", t, func(ctx convey.C) {
		res, err := d.Channels(c, tids)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldHaveLength, 0)
		})
	})
}

func TestDaoChannelMap(t *testing.T) {
	var (
		c    = context.TODO()
		tids = []int64{0}
	)
	convey.Convey("ChannelMap", t, func(ctx convey.C) {
		res, err := d.ChannelMap(c, tids)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldHaveLength, 0)
		})
	})
}

func TestDaoChannelAll(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("ChannelAll", t, func(ctx convey.C) {
		res, tids, err := d.ChannelAll(c)
		ctx.Convey("Then err should be nil.res,tids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(tids, convey.ShouldNotBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoChannelsByType(t *testing.T) {
	var (
		c  = context.TODO()
		tp = int64(0)
	)
	convey.Convey("ChannelsByType", t, func(ctx convey.C) {
		res, tids, err := d.ChannelsByType(c, tp)
		ctx.Convey("Then err should be nil.res,tids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(tids, convey.ShouldNotBeEmpty)
			ctx.So(res, convey.ShouldNotBeEmpty)
		})
	})
}

func TestDaoChanneList(t *testing.T) {
	var (
		c     = context.TODO()
		sqls  = []string{}
		order = "id"
		sort  = "DESC"
		start = int32(0)
		end   = int32(0)
	)
	convey.Convey("ChanneList", t, func(ctx convey.C) {
		res, ids, err := d.ChanneList(c, sqls, order, sort, start, end)
		ctx.Convey("Then err should be nil.res,ids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ids, convey.ShouldHaveLength, 0)
			ctx.So(res, convey.ShouldHaveLength, 0)
		})
	})
}

func TestDaoCountChanneList(t *testing.T) {
	var (
		c    = context.TODO()
		sqls = []string{}
	)
	convey.Convey("CountChanneList", t, func(ctx convey.C) {
		count, err := d.CountChanneList(c, sqls)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoCountChannelByType(t *testing.T) {
	var (
		c  = context.TODO()
		tp = int64(0)
	)
	convey.Convey("CountChannelByType", t, func(ctx convey.C) {
		count, err := d.CountChannelByType(c, tp)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoCountRecommendChannel(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("CountRecommendChannel", t, func(ctx convey.C) {
		count, err := d.CountRecommendChannel(c)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTxInsertChannel(t *testing.T) {
	var channel = &model.Channel{}
	convey.Convey("TxInsertChannel", t, func(ctx convey.C) {
		tx, err := d.BeginTran(context.TODO())
		if err != nil {
			return
		}
		id, err := d.TxInsertChannel(tx, channel)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(id, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
		tx.Rollback()
	})
}

func TestDaoUpdateChannel(t *testing.T) {
	var (
		c       = context.TODO()
		channel = &model.Channel{}
	)
	convey.Convey("UpdateChannel", t, func(ctx convey.C) {
		affect, err := d.UpdateChannel(c, channel)
		ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affect, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTxUpChannel(t *testing.T) {
	var channel = &model.Channel{}
	convey.Convey("TxUpChannel", t, func(ctx convey.C) {
		tx, err := d.BeginTran(context.TODO())
		if err != nil {
			return
		}
		affect, err := d.TxUpChannel(tx, channel)
		ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affect, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
		tx.Rollback()
	})
}

func TestDaoUpdateChannels(t *testing.T) {
	var (
		c        = context.TODO()
		channels = []*model.Channel{
			{
				ID:   22233,
				Name: "unit test",
			},
		}
	)
	convey.Convey("UpdateChannels", t, func(ctx convey.C) {
		affect, err := d.UpdateChannels(c, channels)
		ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affect, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}

func TestDaoRecommandChannel(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("RecommandChannel", t, func(ctx convey.C) {
		res, tids, err := d.RecommandChannel(c)
		ctx.Convey("Then err should be nil.res,tids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(tids, convey.ShouldNotBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTxUpChannelAttr(t *testing.T) {
	var channel = &model.Channel{
		ID:       1,
		Attr:     8,
		Operator: "ut",
	}
	convey.Convey("TxUpChannelAttr", t, func(ctx convey.C) {
		tx, err := d.BeginTran(context.TODO())
		if err != nil {
			return
		}
		id, err := d.TxUpChannelAttr(tx, channel)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(id, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
		tx.Rollback()
	})
}
