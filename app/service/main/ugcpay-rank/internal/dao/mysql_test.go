package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoBeginTran(t *testing.T) {
	convey.Convey("BeginTran", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			tx, err := d.BeginTran(c)
			ctx.Convey("Then err should be nil.tx should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(tx, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRawElecUPRankList(t *testing.T) {
	convey.Convey("RawElecUPRankList", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			upMID = int64(15555180)
			ver   = int64(0)
			limit = int(10)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			list, err := d.RawElecUPRankList(c, upMID, ver, limit)
			ctx.Convey("Then err should be nil.list should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(list, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRawElecUPRank(t *testing.T) {
	convey.Convey("RawElecUPRank", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			upMID  = int64(15555180)
			ver    = int64(0)
			payMID = int64(14137123)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			data, err := d.RawElecUPRank(c, upMID, ver, payMID)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRawCountElecUPRank(t *testing.T) {
	convey.Convey("RawCountElecUPRank", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			upMID = int64(15555180)
			ver   = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			count, err := d.RawCountElecUPRank(c, upMID, ver)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRawElecAVRankList(t *testing.T) {
	convey.Convey("RawElecAVRankList", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			avID  = int64(4052629)
			ver   = int64(0)
			limit = int(10)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			list, err := d.RawElecAVRankList(c, avID, ver, limit)
			ctx.Convey("Then err should be nil.list should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(list, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRawElecAVRank(t *testing.T) {
	convey.Convey("RawElecAVRank", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			avID   = int64(4052629)
			ver    = int64(0)
			payMID = int64(4780461)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			data, err := d.RawElecAVRank(c, avID, ver, payMID)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRawCountElecAVRank(t *testing.T) {
	convey.Convey("RawCountElecAVRank", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			avID = int64(4052629)
			ver  = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			count, err := d.RawCountElecAVRank(c, avID, ver)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRawElecUPMessages(t *testing.T) {
	convey.Convey("RawElecUPMessages", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			payMIDs = []int64{27515316}
			upMID   = int64(27515241)
			ver     = int64(201705)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			dataMap, err := d.RawElecUPMessages(c, payMIDs, upMID, ver)
			ctx.Convey("Then err should be nil.dataMap should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(dataMap, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRawElecAVMessagesByVer(t *testing.T) {
	convey.Convey("RawElecAVMessages", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			payMIDs = []int64{27515316}
			avID    = int64(123)
			ver     = int64(201705)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			dataMap, err := d.RawElecAVMessagesByVer(c, payMIDs, avID, ver)
			ctx.Convey("Then err should be nil.dataMap should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(dataMap, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRawElecAVMessages(t *testing.T) {
	convey.Convey("RawElecAVMessages", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			payMIDs = []int64{27515316}
			avID    = int64(123)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			dataMap, err := d.RawElecAVMessages(c, payMIDs, avID)
			ctx.Convey("Then err should be nil.dataMap should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(dataMap, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRawElecUPUserRank(t *testing.T) {
	convey.Convey("RawElecUPUserRank", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			upMID     = int64(15555180)
			ver       = int64(0)
			payAmount = int64(200)
			mtime     = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rank, err := d.RawElecUPUserRank(c, upMID, ver, payAmount, mtime)
			ctx.Convey("Then err should be nil.rank should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rank, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRawElecAVUserRank(t *testing.T) {
	convey.Convey("RawElecAVUserRank", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			avID      = int64(4052629)
			ver       = int64(0)
			payAmount = int64(200)
			mtime     = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rank, err := d.RawElecAVUserRank(c, avID, ver, payAmount, mtime)
			ctx.Convey("Then err should be nil.rank should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rank, convey.ShouldNotBeNil)
			})
		})
	})
}
func TestDaoRawCountUPTotalElec(t *testing.T) {
	convey.Convey("RawCountUPTotalElec", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			upMID = int64(46333)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			count, err := d.RawCountUPTotalElec(c, upMID)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}
func TestDaoRawElecUserSettings(t *testing.T) {
	convey.Convey("RawElecUserSetting", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = 0
			limit = 10
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			value, maxID, err := d.RawElecUserSettings(c, id, limit)
			ctx.Convey("Then err should be nil.value should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(maxID, convey.ShouldBeGreaterThan, 0)
				ctx.So(value, convey.ShouldNotBeNil)
				ctx.So(value, convey.ShouldHaveLength, 10)
			})
		})
	})
}
