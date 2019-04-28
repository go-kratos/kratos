package dao

import (
	"context"
	"go-common/app/job/main/credit/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaovoteIndexKey(t *testing.T) {
	convey.Convey("voteIndexKey", t, func(ctx convey.C) {
		var (
			cid   = int64(1)
			otype = int8(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := voteIndexKey(cid, otype)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaocaseIndexKey(t *testing.T) {
	convey.Convey("caseIndexKey", t, func(ctx convey.C) {
		var (
			cid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := caseIndexKey(cid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoblockIndexKey(t *testing.T) {
	convey.Convey("blockIndexKey", t, func(ctx convey.C) {
		var (
			otype = int64(1)
			btype = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := blockIndexKey(otype, btype)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelCaseIdx(t *testing.T) {
	convey.Convey("DelCaseIdx", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			cid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelCaseIdx(c, cid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelBlockedInfoIdx(t *testing.T) {
	convey.Convey("DelBlockedInfoIdx", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			bl = &model.BlockedInfo{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelBlockedInfoIdx(c, bl)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddBlockInfoIdx(t *testing.T) {
	convey.Convey("AddBlockInfoIdx", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			bl = &model.BlockedInfo{MTime: "2018-11-14 00:00:00"}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddBlockInfoIdx(c, bl)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelVoteIdx(t *testing.T) {
	convey.Convey("DelVoteIdx", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			cid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelVoteIdx(c, cid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetGrantCase(t *testing.T) {
	convey.Convey("SetGrantCase", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mcases = make(map[int64]*model.SimCase)
		)
		mcases[1] = &model.SimCase{}
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetGrantCase(c, mcases)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelGrantCase(t *testing.T) {
	convey.Convey("DelGrantCase", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			cids = []int64{1, 2}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelGrantCase(c, cids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoTotalGrantCase(t *testing.T) {
	convey.Convey("TotalGrantCase", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.TotalGrantCase(c)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGrantCases(t *testing.T) {
	convey.Convey("GrantCases", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.GrantCases(c)
			ctx.Convey("Then err should be nil.cids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				// ctx.So(cids, convey.ShouldNotBeNil)
			})
		})
	})
}
