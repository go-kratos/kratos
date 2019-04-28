package dao

import (
	"context"
	"go-common/app/interface/main/credit/model"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddBlockedInfo(t *testing.T) {
	convey.Convey("AddBlockedInfo", t, func(convCtx convey.C) {
		var (
			c = context.Background()
			r = &model.BlockedInfo{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.AddBlockedInfo(c, r)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoTxAddBlockedInfo(t *testing.T) {
	convey.Convey("TxAddBlockedInfo", t, func(convCtx convey.C) {
		var (
			tx, _ = d.BeginTran(context.Background())
			rs    = []*model.BlockedInfo{}
			r     = &model.BlockedInfo{Uname: "test", UID: 1024}
		)
		rs = append(rs, r)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.TxAddBlockedInfo(tx, rs)
			if err == nil {
				if err = tx.Commit(); err != nil {
					tx.Rollback()
				}
			} else {
				tx.Rollback()
			}
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoBlockedCount(t *testing.T) {
	convey.Convey("BlockedCount", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			count, err := d.BlockedCount(c, mid)
			convCtx.Convey("Then err should be nil.count should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBlockedNumUser(t *testing.T) {
	convey.Convey("BlockedNumUser", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			count, err := d.BlockedNumUser(c, mid)
			convCtx.Convey("Then err should be nil.count should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBLKHistoryCount(t *testing.T) {
	convey.Convey("BLKHistoryCount", t, func(convCtx convey.C) {
		var (
			c      = context.Background()
			ArgHis = &model.ArgHistory{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			count, err := d.BLKHistoryCount(c, ArgHis)
			convCtx.Convey("Then err should be nil.count should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBlockTotalTime(t *testing.T) {
	convey.Convey("BlockTotalTime", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			ts  = time.Now()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			total, err := d.BlockTotalTime(c, mid, ts)
			convCtx.Convey("Then err should be nil.total should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(total, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBlockedUserList(t *testing.T) {
	convey.Convey("BlockedUserList", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.BlockedUserList(c, mid)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBlockedList(t *testing.T) {
	convey.Convey("BlockedList", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			otype = int8(0)
			btype = int8(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.BlockedList(c, otype, btype)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBLKHistorys(t *testing.T) {
	convey.Convey("BLKHistorys", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			ah = &model.ArgHistory{MID: 0}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.BLKHistorys(c, ah)
			convCtx.Convey("Then err should be nil.res should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoBlockedInfoByID(t *testing.T) {
	convey.Convey("BlockedInfoByID", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			id = int64(234)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			r, err := d.BlockedInfoByID(c, id)
			convCtx.Convey("Then err should be nil.r should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(r, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBlockedInfoIDs(t *testing.T) {
	convey.Convey("BlockedInfoIDs", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{1, 234, 27515668}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.BlockedInfoIDs(c, ids)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBlockedInfos(t *testing.T) {
	convey.Convey("BlockedInfos", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{243, 629}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.BlockedInfos(c, ids)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
