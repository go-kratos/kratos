package reply

import (
	"context"
	"go-common/app/job/main/reply/model/reply"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestReplyTxInsert(t *testing.T) {
	convey.Convey("TxInsert", t, func(ctx convey.C) {
		var (
			tx, _ = d.mysql.Begin(context.Background())
			r     = &reply.Reply{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.Reply.TxInsert(tx, r)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		tx.Rollback()
	})
}

func TestReplyTxIncrCount2(t *testing.T) {
	convey.Convey("TxIncrCount", t, func(ctx convey.C) {
		var (
			tx, _ = d.mysql.Begin(context.Background())
			oid   = int64(0)
			rpID  = int64(0)
			now   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.Reply.TxIncrCount(tx, oid, rpID, now)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		tx.Rollback()
	})
}

func TestReplyTxIncrFCount2(t *testing.T) {
	convey.Convey("TxIncrFCount", t, func(ctx convey.C) {
		var (
			tx, _ = d.mysql.Begin(context.Background())
			oid   = int64(0)
			rpID  = int64(0)
			now   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.Reply.TxIncrFCount(tx, oid, rpID, now)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		tx.Rollback()
	})
}

func TestReplyTxIncrRCount2(t *testing.T) {
	convey.Convey("TxIncrRCount", t, func(ctx convey.C) {
		var (
			tx, _ = d.mysql.Begin(context.Background())
			oid   = int64(0)
			rpID  = int64(0)
			now   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.Reply.TxIncrRCount(tx, oid, rpID, now)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		tx.Rollback()
	})
}

func TestReplyTxDecrCount2(t *testing.T) {
	convey.Convey("TxDecrCount", t, func(ctx convey.C) {
		var (
			tx, _ = d.mysql.Begin(context.Background())
			oid   = int64(0)
			rpID  = int64(0)
			now   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.Reply.TxDecrCount(tx, oid, rpID, now)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		tx.Rollback()
	})
}

func TestReplyGetForUpdate(t *testing.T) {
	convey.Convey("GetForUpdate", t, func(ctx convey.C) {
		var (
			tx, _ = d.mysql.Begin(context.Background())
			oid   = int64(0)
			rpID  = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			r, err := d.Reply.GetForUpdate(tx, oid, rpID)
			ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(r, convey.ShouldBeNil)
			})
		})
		tx.Rollback()
	})
}

func TestReplyTxUpState(t *testing.T) {
	convey.Convey("TxUpState", t, func(ctx convey.C) {
		var (
			tx, _ = d.mysql.Begin(context.Background())
			oid   = int64(0)
			rpID  = int64(0)
			state = int8(0)
			now   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.Reply.TxUpState(tx, oid, rpID, state, now)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		tx.Rollback()
	})
}

func TestReplyUpState(t *testing.T) {
	convey.Convey("UpState", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			oid   = int64(0)
			rpID  = int64(0)
			state = int8(0)
			now   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.Reply.UpState(c, oid, rpID, state, now)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyTxUpAttr2(t *testing.T) {
	convey.Convey("TxUpAttr", t, func(ctx convey.C) {
		var (
			tx, _ = d.mysql.Begin(context.Background())
			oid   = int64(0)
			rpID  = int64(0)
			attr  = uint32(0)
			now   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.Reply.TxUpAttr(tx, oid, rpID, attr, now)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		tx.Rollback()
	})
}

func TestReplyUpLike(t *testing.T) {
	convey.Convey("UpLike", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			rpID = int64(0)
			like = int(0)
			hate = int(0)
			now  = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.Reply.UpLike(c, oid, rpID, like, hate, now)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyGet3(t *testing.T) {
	convey.Convey("Get", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			rpID = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			r, err := d.Reply.Get(c, oid, rpID)
			ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(r, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyGetTop2(t *testing.T) {
	convey.Convey("GetTop", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			tp  = int8(0)
			bit = uint32(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			r, err := d.Reply.GetTop(c, oid, tp, bit)
			ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(r, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyGetByDialog(t *testing.T) {
	convey.Convey("GetByDialog", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			oid    = int64(0)
			typ    = int8(0)
			root   = int64(0)
			dialog = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rps, err := d.Reply.GetByDialog(c, oid, typ, root, dialog)
			ctx.Convey("Then err should be nil.rps should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rps, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyGetAllInSlice(t *testing.T) {
	convey.Convey("GetAllInSlice", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			oid      = int64(0)
			typ      = int8(0)
			maxFloor = int(0)
			shard    = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rs, err := d.Reply.GetAllInSlice(c, oid, typ, maxFloor, shard)
			ctx.Convey("Then err should be nil.rs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(rs, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyGetByFloorLimit(t *testing.T) {
	convey.Convey("GetByFloorLimit", t, func(ctx convey.C) {
		var (
			oid   = int64(0)
			typ   = int8(0)
			floor = int(0)
			limit = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rs, err := d.Reply.GetByFloorLimit(context.Background(), oid, typ, floor, limit)
			ctx.Convey("Then err should be nil.rs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rs, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplygetByFloorLimit(t *testing.T) {
	convey.Convey("getByFloorLimit", t, func(ctx convey.C) {
		var (
			oid   = int64(0)
			typ   = int8(0)
			floor = int(0)
			limit = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rs, err := d.Reply.getByFloorLimit(context.Background(), oid, typ, floor, limit)
			ctx.Convey("Then err should be nil.rs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rs, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyGetByLikeLimit(t *testing.T) {
	convey.Convey("GetByLikeLimit", t, func(ctx convey.C) {
		var (
			oid   = int64(0)
			typ   = int8(0)
			limit = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rs, err := d.Reply.GetByLikeLimit(context.Background(), oid, typ, limit)
			ctx.Convey("Then err should be nil.rs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rs, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyGetByCountLimit(t *testing.T) {
	convey.Convey("GetByCountLimit", t, func(ctx convey.C) {
		var (
			oid   = int64(0)
			typ   = int8(0)
			limit = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rs, err := d.Reply.GetByCountLimit(context.Background(), oid, typ, limit)
			ctx.Convey("Then err should be nil.rs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rs, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyGetAllByFloor(t *testing.T) {
	convey.Convey("GetAllByFloor", t, func(ctx convey.C) {
		var (
			oid   = int64(0)
			typ   = int8(0)
			start = int(0)
			end   = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rs, err := d.Reply.GetAllByFloor(context.Background(), oid, typ, start, end)
			ctx.Convey("Then err should be nil.rs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rs, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyGetAll(t *testing.T) {
	convey.Convey("GetAll", t, func(ctx convey.C) {
		var (
			oid = int64(0)
			tp  = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rs, err := d.Reply.GetAll(context.Background(), oid, tp)
			ctx.Convey("Then err should be nil.rs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rs, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyGetAllByRoot(t *testing.T) {
	convey.Convey("GetAllByRoot", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			rpID = int64(0)
			tp   = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rs, err := d.Reply.GetAllByRoot(c, oid, rpID, tp)
			ctx.Convey("Then err should be nil.rs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rs, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyGetsByRoot(t *testing.T) {
	convey.Convey("GetsByRoot", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			oid   = int64(0)
			rpID  = int64(0)
			tp    = int8(0)
			state = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rs, err := d.Reply.GetsByRoot(c, oid, rpID, tp, state)
			ctx.Convey("Then err should be nil.rs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rs, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyFixDialogGetRepliesByRoot(t *testing.T) {
	convey.Convey("FixDialogGetRepliesByRoot", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			oid    = int64(0)
			tp     = int8(0)
			rootID = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rps, err := d.Reply.FixDialogGetRepliesByRoot(c, oid, tp, rootID)
			ctx.Convey("Then err should be nil.rps should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rps, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyFixDialogSetDialogBatch(t *testing.T) {
	convey.Convey("FixDialogSetDialogBatch", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			oid    = int64(0)
			dialog = int64(0)
			rpIDs  = []int64{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			d.Reply.FixDialogSetDialogBatch(c, oid, dialog, rpIDs)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}
