package dao

import (
	"testing"
	"time"

	"go-common/app/service/main/ugcpay-rank/internal/conf"
	"go-common/app/service/main/ugcpay-rank/internal/model"
	xtime "go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoLCLoadElecAVRankFromNil(t *testing.T) {
	convey.Convey("loadElecAVRank", t, func(ctx convey.C) {
		var (
			avID = int64(233)
			ver  = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rank, err := d.LCLoadElecAVRank(avID, ver)
			ctx.Convey("Then err should be nil.rank should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rank, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoLCStoreElecAVRank(t *testing.T) {
	convey.Convey("storeElecAVRank", t, func(ctx convey.C) {
		var (
			avID = int64(233)
			ver  = int64(0)
			rank = &model.RankElecAVProto{
				AVID: 233,
			}
		)
		conf.Conf.LocalCache.ElecAVRankTTL = xtime.Duration(time.Second)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.LCStoreElecAVRank(avID, ver, rank)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoLCLoadElecAVRank(t *testing.T) {
	convey.Convey("loadElecAVRank", t, func(ctx convey.C) {
		var (
			avID = int64(233)
			ver  = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rank, err := d.LCLoadElecAVRank(avID, ver)
			ctx.Convey("Then err should be nil.rank should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rank, convey.ShouldNotBeNil)
				ctx.So(rank.AVID, convey.ShouldEqual, avID)
			})
		})
	})
}

func TestDaoLCLoadElecAVRankAfterTTL(t *testing.T) {
	convey.Convey("loadElecAVRankAfterTTL", t, func(ctx convey.C) {
		var (
			avID = int64(233)
			ver  = int64(0)
		)
		<-time.After(time.Duration(conf.Conf.LocalCache.ElecAVRankTTL))
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rank, err := d.LCLoadElecAVRank(avID, ver)
			ctx.Convey("Then err should be nil.rank should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rank, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoLCLoadElecUPRankFromNil(t *testing.T) {
	convey.Convey("loadElecUPRank", t, func(ctx convey.C) {
		var (
			upMID = int64(4633)
			ver   = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rank, err := d.LCLoadElecUPRank(upMID, ver)
			ctx.Convey("Then err should be nil.rank should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rank, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoLCStoreElecUPRank(t *testing.T) {
	convey.Convey("storeElecUPRank", t, func(ctx convey.C) {
		var (
			upMID = int64(4633)
			ver   = int64(0)
			rank  = &model.RankElecUPProto{
				Count: 100,
			}
		)
		conf.Conf.LocalCache.ElecUPRankTTL = xtime.Duration(time.Second)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.LCStoreElecUPRank(upMID, ver, rank)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoLCLoadElecUPRank(t *testing.T) {
	convey.Convey("loadElecUPRank", t, func(ctx convey.C) {
		var (
			upMID = int64(4633)
			ver   = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rank, err := d.LCLoadElecUPRank(upMID, ver)
			ctx.Convey("Then err should be nil.rank should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rank, convey.ShouldNotBeNil)
				ctx.So(rank.Count, convey.ShouldEqual, 100)
			})
		})
	})
}

func TestDaoLCLoadElecUPRankAfterTTL(t *testing.T) {
	convey.Convey("loadElecUPRankAfterTTL", t, func(ctx convey.C) {
		var (
			upMID = int64(4633)
			ver   = int64(0)
		)
		<-time.After(time.Duration(conf.Conf.LocalCache.ElecUPRankTTL))
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rank, err := d.LCLoadElecUPRank(upMID, ver)
			ctx.Convey("Then err should be nil.rank should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rank, convey.ShouldBeNil)
			})
		})
	})
}
