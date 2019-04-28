package dao

import (
	"context"
	"go-common/app/admin/main/tag/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaosetChannelRules(t *testing.T) {
	var channelRules = []*model.ChannelRule{}
	convey.Convey("setChannelRules", t, func(ctx convey.C) {
		p1 := setChannelRules(channelRules)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoChannelRule(t *testing.T) {
	var (
		c   = context.TODO()
		tid = int64(0)
	)
	convey.Convey("ChannelRule", t, func(ctx convey.C) {
		res, tids, err := d.ChannelRule(c, tid)
		ctx.Convey("Then err should be nil.res,tids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(tids, convey.ShouldHaveLength, 0)
			ctx.So(res, convey.ShouldHaveLength, 0)
		})
	})
}

func TestDaoChannelRules(t *testing.T) {
	var (
		c    = context.TODO()
		tids = []int64{0}
	)
	convey.Convey("ChannelRules", t, func(ctx convey.C) {
		res, err := d.ChannelRules(c, tids)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldHaveLength, 0)
		})
	})
}

func TestDaoTxUpChannelRuleState(t *testing.T) {
	var (
		tid   = int64(0)
		state = int32(0)
		uname = ""
	)
	convey.Convey("TxUpChannelRuleState", t, func(ctx convey.C) {
		tx, err := d.BeginTran(context.TODO())
		if err != nil {
			return
		}
		affect, err := d.TxUpChannelRuleState(tx, tid, state, uname)
		ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affect, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
		tx.Rollback()
	})
}

func TestDaoTxUpdateChannelRules(t *testing.T) {
	var channelRules = []*model.ChannelRule{
		{
			Tid:    1,
			InRule: "2",
		},
	}
	convey.Convey("TxUpdateChannelRules", t, func(ctx convey.C) {
		tx, err := d.BeginTran(context.TODO())
		if err != nil {
			return
		}
		affect, err := d.TxUpdateChannelRules(tx, channelRules)
		ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affect, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
		tx.Rollback()
	})
}
