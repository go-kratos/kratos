package reply

import (
	"context"
	model "go-common/app/interface/main/reply/model/reply"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestReplyChildrenIDsOfRootReply(t *testing.T) {
	convey.Convey("ChildrenIDsOfRootReply", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			oid    = int64(0)
			rootID = int64(0)
			tp     = int8(0)
			offset = int(0)
			limit  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.Reply.ChildrenIDsOfRootReply(c, oid, rootID, tp, offset, limit)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyCacheKeyRootReplyIDs(t *testing.T) {
	convey.Convey("CacheKeyRootReplyIDs", t, func(ctx convey.C) {
		var (
			oid  = int64(0)
			tp   = int8(0)
			sort = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.Redis.CacheKeyRootReplyIDs(oid, tp, sort)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyParentChildrenReplyIDMap(t *testing.T) {
	convey.Convey("ParentChildrenReplyIDMap", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			parentIDs = []int64{}
			start     = int(0)
			end       = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			parentChildrenMap, missedIDs, err := d.Redis.ParentChildrenReplyIDMap(c, parentIDs, start, end)
			ctx.Convey("Then err should be nil.parentChildrenMap,missedIDs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(missedIDs, convey.ShouldBeNil)
				ctx.So(parentChildrenMap, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyRangeChildrenReplyIDs(t *testing.T) {
	convey.Convey("RangeChildrenReplyIDs", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			keys  = []string{}
			start = int(0)
			end   = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			arrOfChildrenReplyIDs, missedKeys, err := d.Redis.RangeChildrenReplyIDs(c, keys, start, end)
			ctx.Convey("Then err should be nil.arrOfChildrenReplyIDs,missedKeys should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(missedKeys, convey.ShouldBeNil)
				ctx.So(arrOfChildrenReplyIDs, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyRangeChildrenIDByCursorScore(t *testing.T) {
	convey.Convey("RangeChildrenIDByCursorScore", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			key    = ""
			cursor = &model.Cursor{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.Redis.RangeChildrenIDByCursorScore(c, key, cursor)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyRangeRootIDByCursorScore(t *testing.T) {
	convey.Convey("RangeRootIDByCursorScore", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			key    = ""
			cursor = &model.Cursor{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, p2, err := d.Redis.RangeRootIDByCursorScore(c, key, cursor)
			ctx.Convey("Then err should be nil.p1,p2 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p2, convey.ShouldNotBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyRangeRootReplyIDs(t *testing.T) {
	convey.Convey("RangeRootReplyIDs", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			key   = ""
			start = int(0)
			end   = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.Redis.RangeRootReplyIDs(c, key, start, end)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyExpireCache(t *testing.T) {
	convey.Convey("ExpireCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			key = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.Redis.ExpireCache(c, key)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplygenChildrenKeyParentIDMap(t *testing.T) {
	convey.Convey("genChildrenKeyParentIDMap", t, func(ctx convey.C) {
		var (
			ids = []int64{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := genChildrenKeyParentIDMap(ids)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplygenChildrenKeyByRootReplyIDs(t *testing.T) {
	convey.Convey("genChildrenKeyByRootReplyIDs", t, func(ctx convey.C) {
		var (
			ids = []int64{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := genChildrenKeyByRootReplyIDs(ids)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplygenNewChildrenKeyByRootReplyID(t *testing.T) {
	convey.Convey("GenNewChildrenKeyByRootReplyID", t, func(ctx convey.C) {
		var (
			id = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := GenNewChildrenKeyByRootReplyID(id)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplygenReplyKeyByID(t *testing.T) {
	convey.Convey("genReplyKeyByID", t, func(ctx convey.C) {
		var (
			id = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := genReplyKeyByID(id)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplycontains(t *testing.T) {
	convey.Convey("contains", t, func(ctx convey.C) {
		var (
			arr = []string{}
			b   = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := contains(arr, b)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyGetReplyByIDs(t *testing.T) {
	convey.Convey("GetReplyByIDs", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, p2, err := d.Mc.GetReplyByIDs(c, ids)
			ctx.Convey("Then err should be nil.p1,p2 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p2, convey.ShouldNotBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyChildrenIDSortByFloorCursor(t *testing.T) {
	convey.Convey("ChildrenIDSortByFloorCursor", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			oid    = int64(0)
			tp     = int8(0)
			rootID = int64(0)
			cursor = &model.Cursor{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.Reply.ChildrenIDSortByFloorCursor(c, oid, tp, rootID, cursor)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyRootIDSortByFloorCursor(t *testing.T) {
	convey.Convey("RootIDSortByFloorCursor", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			oid    = int64(0)
			tp     = int8(0)
			cursor = &model.Cursor{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.Reply.RootIDSortByFloorCursor(c, oid, tp, cursor)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
