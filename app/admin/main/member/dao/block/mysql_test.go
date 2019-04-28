package block

import (
	"context"
	model "go-common/app/admin/main/member/model/block"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestBlockhistoryIdx(t *testing.T) {
	convey.Convey("historyIdx", t, func(ctx convey.C) {
		var (
			mid = int64(46333)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := historyIdx(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestBlockUser(t *testing.T) {
	convey.Convey("User", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(46333)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			user, err := d.User(c, mid)
			ctx.Convey("Then err should be nil.user should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(user, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestBlockUsers(t *testing.T) {
	convey.Convey("Users", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{46333}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			users, err := d.Users(c, mids)
			ctx.Convey("Then err should be nil.users should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(users, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestBlockTxUpdateUser(t *testing.T) {
	convey.Convey("TxUpdateUser", t, func(ctx convey.C) {
		var (
			c                        = context.Background()
			tx, _                    = d.BeginTX(c)
			mid                      = int64(46333)
			status model.BlockStatus = model.BlockStatusFalse
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.TxUpdateUser(c, tx, mid, status)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestBlockUserDetails(t *testing.T) {
	convey.Convey("UserDetails", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{46333}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			users, err := d.UserDetails(c, mids)
			ctx.Convey("Then err should be nil.users should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(users, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestBlockUpdateAddBlockCount(t *testing.T) {
	convey.Convey("UpdateAddBlockCount", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(46333)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.UpdateAddBlockCount(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestBlockHistory(t *testing.T) {
	convey.Convey("History", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(46333)
			start = int(0)
			limit = int(10)
			desc  bool
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			history, err := d.History(c, mid, start, limit, desc)
			ctx.Convey("Then err should be nil.history should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(history, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestBlockHistoryCount(t *testing.T) {
	convey.Convey("HistoryCount", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(46333)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			total, err := d.HistoryCount(c, mid)
			ctx.Convey("Then err should be nil.total should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestBlockTxInsertHistory(t *testing.T) {
	convey.Convey("TxInsertHistory", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTX(c)
			h     = &model.DBHistory{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.TxInsertHistory(c, tx, h)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestBlockintsToStrs(t *testing.T) {
	convey.Convey("intsToStrs", t, func(ctx convey.C) {
		var (
			ints = []int64{46333}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			strs := intsToStrs(ints)
			ctx.Convey("Then strs should not be nil.", func(ctx convey.C) {
				ctx.So(strs, convey.ShouldNotBeNil)
			})
		})
	})
}
