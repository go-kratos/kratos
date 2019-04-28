package block

import (
	"context"
	"testing"
	"time"

	model "go-common/app/service/main/member/model/block"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaohistoryIdx(t *testing.T) {
	convey.Convey("historyIdx", t, func(ctx convey.C) {
		var (
			mid = int64(46333)
		)
		ctx.Convey("When everything right", func(ctx convey.C) {
			shard := historyIdx(mid)
			ctx.Convey("Then shard should equal mid % 10.", func(ctx convey.C) {
				ctx.So(shard, convey.ShouldEqual, 3)
			})
		})
	})
}

func TestDaoUser(t *testing.T) {
	convey.Convey("User", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(46333)
		)
		ctx.Convey("When everything right", func(ctx convey.C) {
			user, err := d.User(c, mid)
			ctx.Convey("Then err should be nil.user should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(user, convey.ShouldNotBeNil)
				ctx.So(user.MID, convey.ShouldEqual, mid)
			})
		})
	})
}

func TestDaoTxInsertHistory(t *testing.T) {
	convey.Convey("TxInsertHistory", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.db.Begin(c)
			his   = &model.DBHistory{
				MID:       46333,
				AdminID:   2333,
				AdminName: "ut",
				Source:    model.BlockSourceBlackHouse,
				Area:      model.BlockAreaAlbum,
				Reason:    "ut reason",
				Comment:   "ut comment",
				Action:    model.BlockActionLimit,
				StartTime: time.Now(),
				Duration:  60,
				Notify:    false,
				CTime:     time.Now(),
				MTime:     time.Now(),
			}
		)
		ctx.Convey("When everything right", func(ctx convey.C) {
			err := d.TxInsertHistory(c, tx, his)
			ctx.So(err, convey.ShouldBeNil)
			err = tx.Commit()
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.Convey("When get UserLastHistory", func(ctx convey.C) {
					his2, err := d.UserLastHistory(c, his.MID)
					ctx.Convey("Then err should be nil.his2 should resemble his.", func(ctx convey.C) {
						ctx.So(err, convey.ShouldBeNil)
						ctx.So(his2.MID, convey.ShouldEqual, his.MID)
						ctx.So(his2.Action, convey.ShouldEqual, his.Action)
						ctx.So(his2.StartTime.Unix(), convey.ShouldEqual, his.StartTime.Unix())
						ctx.So(his2.Duration, convey.ShouldEqual, his.Duration)
					})
				})
			})
		})
	})
}

func TestDaoTxUpdateUser(t *testing.T) {
	var (
		c      = context.Background()
		tx, _  = d.db.Begin(c)
		mid    = int64(46333)
		status = model.BlockStatusFalse
	)
	convey.Convey("TxUpdateUser", t, func(ctx convey.C) {
		ctx.Convey("When everything right", func(ctx convey.C) {
			err := d.TxUpdateUser(c, tx, mid, status)
			ctx.So(err, convey.ShouldBeNil)
			err = tx.Commit()
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUserStatusMapWithMIDs(t *testing.T) {
	convey.Convey("UserStatusMapWithMIDs", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			status = model.BlockStatusFalse
			mids   = []int64{46333, 2, 35858}
		)
		ctx.Convey("When everything right", func(ctx convey.C) {
			midMap, err := d.UserStatusMapWithMIDs(c, status, mids)
			ctx.Convey("Then err should be nil.midMap should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(midMap, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateAddBlockCount(t *testing.T) {
	convey.Convey("UpdateAddBlockCount", t, func(ctx convey.C) {
		var (
			c   = context.TODO()
			mid = int64(46333)
		)
		ctx.Convey("When everything right", func(ctx convey.C) {
			err := d.UpdateAddBlockCount(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestUserDetails(t *testing.T) {
	convey.Convey("get user details", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(46333)
		)
		ctx.Convey("When get user details from db", func(ctx convey.C) {
			_, err := d.UserDetails(c, []int64{mid})
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
