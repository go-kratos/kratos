package reply

import (
	"context"
	"go-common/library/database/sql"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestReplyNewReplyDao(t *testing.T) {
	convey.Convey("NewReplyDao", t, func(ctx convey.C) {
		var (
			db      = &sql.DB{}
			dbSlave = &sql.DB{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			dao := NewReplyDao(db, dbSlave)
			ctx.Convey("Then dao should not be nil.", func(ctx convey.C) {
				ctx.So(dao, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyReplyhit(t *testing.T) {
	convey.Convey("hit", t, func(ctx convey.C) {
		var (
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.Reply.hit(oid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyCountLike(t *testing.T) {
	convey.Convey("CountLike", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			tp  = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.Reply.CountLike(c, oid, tp)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyReplyGet(t *testing.T) {
	convey.Convey("Get", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			rpID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			r, err := d.Reply.Get(c, oid, rpID)
			ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(r, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyGetIDsByDialogAsc(t *testing.T) {
	convey.Convey("GetIDsByDialogAsc", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			oid      = int64(0)
			tp       = int8(0)
			root     = int64(0)
			dialog   = int64(0)
			maxFloor = int64(0)
			count    = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rpIDs, err := d.Reply.GetIDsByDialogAsc(c, oid, tp, root, dialog, maxFloor, count)
			ctx.Convey("Then err should be nil.rpIDs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rpIDs, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyGetDialogMinMaxFloor(t *testing.T) {
	convey.Convey("GetDialogMinMaxFloor", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			oid    = int64(0)
			tp     = int8(0)
			root   = int64(0)
			dialog = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			minFloor, maxFloor, err := d.Reply.GetDialogMinMaxFloor(c, oid, tp, root, dialog)
			ctx.Convey("Then err should be nil.minFloor,maxFloor should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(maxFloor, convey.ShouldNotBeNil)
				ctx.So(minFloor, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyGetIDsByDialogDesc(t *testing.T) {
	convey.Convey("GetIDsByDialogDesc", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			oid      = int64(0)
			tp       = int8(0)
			root     = int64(0)
			dialog   = int64(0)
			minFloor = int64(0)
			count    = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rpIDs, err := d.Reply.GetIDsByDialogDesc(c, oid, tp, root, dialog, minFloor, count)
			ctx.Convey("Then err should be nil.rpIDs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rpIDs, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyGetIDsByDialog(t *testing.T) {
	convey.Convey("GetIDsByDialog", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			oid    = int64(0)
			tp     = int8(0)
			root   = int64(0)
			dialog = int64(0)
			offset = int(0)
			count  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rpIDs, err := d.Reply.GetIDsByDialog(c, oid, tp, root, dialog, offset, count)
			ctx.Convey("Then err should be nil.rpIDs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rpIDs, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyGetByIds(t *testing.T) {
	convey.Convey("GetByIds", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			oid   = int64(0)
			tp    = int8(0)
			rpIds = []int64{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rpMap, err := d.Reply.GetByIds(c, oid, tp, rpIds)
			ctx.Convey("Then err should be nil.rpMap should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rpMap, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyGetIdsSortFloor(t *testing.T) {
	convey.Convey("GetIdsSortFloor", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			oid    = int64(0)
			tp     = int8(0)
			offset = int(0)
			count  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Reply.GetIdsSortFloor(c, oid, tp, offset, count)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyGetIdsSortCount(t *testing.T) {
	convey.Convey("GetIdsSortCount", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			oid    = int64(0)
			tp     = int8(0)
			offset = int(0)
			count  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Reply.GetIdsSortCount(c, oid, tp, offset, count)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyGetIdsSortLike(t *testing.T) {
	convey.Convey("GetIdsSortLike", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			oid    = int64(0)
			tp     = int8(0)
			offset = int(0)
			count  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Reply.GetIdsSortLike(c, oid, tp, offset, count)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyGetIdsByRoot(t *testing.T) {
	convey.Convey("GetIdsByRoot", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			oid    = int64(0)
			root   = int64(0)
			tp     = int8(0)
			offset = int(0)
			count  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Reply.GetIdsByRoot(c, oid, root, tp, offset, count)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyGetIDsByRootWithoutState(t *testing.T) {
	convey.Convey("GetIDsByRootWithoutState", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			oid    = int64(0)
			root   = int64(0)
			tp     = int8(0)
			offset = int(0)
			count  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Reply.GetIDsByRootWithoutState(c, oid, root, tp, offset, count)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyGetIDsByFloorOffset(t *testing.T) {
	convey.Convey("GetIDsByFloorOffset", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			oid   = int64(0)
			tp    = int8(0)
			start = int(0)
			end   = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Reply.GetIDsByFloorOffset(c, oid, tp, start, end)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplySetHateCount(t *testing.T) {
	convey.Convey("SetHateCount", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			oid   = int64(0)
			rpID  = int64(0)
			count = int32(0)
			now   = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.Reply.SetHateCount(c, oid, rpID, count, now)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplySetLikeCount(t *testing.T) {
	convey.Convey("SetLikeCount", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			oid   = int64(0)
			rpID  = int64(0)
			count = int32(0)
			now   = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.Reply.SetLikeCount(c, oid, rpID, count, now)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyAddFilteredReply(t *testing.T) {
	convey.Convey("AddFilteredReply", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			rpID    = int64(-1)
			oid     = int64(0)
			mid     = int64(0)
			tp      = int8(0)
			level   = int8(0)
			message = ""
			now     = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Reply.AddFilteredReply(c, rpID, oid, mid, tp, level, message, now)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		d.mysql.Exec(context.Background(), "DELETE FROM reply_filtered WHERE rpid=-1")
	})
}
