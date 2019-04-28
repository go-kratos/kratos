package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/admin/main/tag/model"
	xtime "go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaosetChannelSynonyms(t *testing.T) {
	var (
		synonyms = []*model.ChannelSynonym{
			{
				PTid:     1,
				Tid:      2,
				Alias:    "",
				Rank:     0,
				Operator: "ut",
				CTime:    xtime.Time(time.Now().Unix()),
				MTime:    xtime.Time(time.Now().Unix()),
			},
		}
	)
	convey.Convey("setChannelSynonyms", t, func(ctx convey.C) {
		p1 := setChannelSynonyms(synonyms)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoChannelSynonymMap(t *testing.T) {
	var (
		c    = context.TODO()
		ptid = int64(-1)
	)
	convey.Convey("ChannelSynonymMap", t, func(ctx convey.C) {
		res, tids, err := d.ChannelSynonymMap(c, ptid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldHaveLength, 0)
			ctx.So(tids, convey.ShouldHaveLength, 0)
		})
	})
}

func TestDaoTxUpStateChannelSynonym(t *testing.T) {
	var (
		ptid  = int64(0)
		state = int32(0)
		uname = "ut"
	)
	convey.Convey("TxUpStateChannelSynonym", t, func(ctx convey.C) {
		tx, err := d.BeginTran(context.TODO())
		if err != nil {
			return
		}
		affect, err := d.TxUpStateChannelSynonym(tx, ptid, state, uname)
		ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affect, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
		tx.Rollback()
	})
}

func TestDaoTxUpdateChannelSynonyms(t *testing.T) {
	var (
		synonyms = []*model.ChannelSynonym{{
			PTid:     1,
			Tid:      2,
			Alias:    "",
			Rank:     0,
			Operator: "ut",
			CTime:    xtime.Time(time.Now().Unix()),
			MTime:    xtime.Time(time.Now().Unix()),
		}}
	)
	convey.Convey("TxUpdateChannelSynonyms", t, func(ctx convey.C) {
		tx, err := d.BeginTran(context.TODO())
		if err != nil {
			return
		}
		affect, err := d.TxUpdateChannelSynonyms(tx, synonyms)
		ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affect, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
		tx.Rollback()
	})
}
