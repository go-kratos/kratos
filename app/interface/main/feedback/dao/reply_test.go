package dao

import (
	"context"
	xsql "database/sql"
	"fmt"
	"reflect"
	"testing"

	"go-common/app/interface/main/feedback/model"
	"go-common/library/database/sql"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

func TestDaoTxAddReply(t *testing.T) {
	convey.Convey("TxAddReply", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.Background())
			r     = &model.Reply{}
		)
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			id, err := d.TxAddReply(tx, r)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When tx.Exec gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(tx), "Exec",
				func(_ *sql.Tx, _ string, _ ...interface{}) (xsql.Result, error) {
					return nil, fmt.Errorf("tx.Exec Error")
				})
			defer guard.Unpatch()
			_, err := d.TxAddReply(tx, r)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Rollback()
		})
	})
}

func TestDaoAddReply(t *testing.T) {
	convey.Convey("AddReply", t, func(ctx convey.C) {
		var (
			c = context.TODO()
			r = &model.Reply{}
		)
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			id, err := d.AddReply(c, r)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When d.inReply.Exec gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.inReply), "Exec",
				func(_ *sql.Stmt, _ context.Context, _ ...interface{}) (xsql.Result, error) {
					return nil, fmt.Errorf("d.inReply.Exec Error")
				})
			defer guard.Unpatch()
			_, err := d.AddReply(c, r)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoReplys(t *testing.T) {
	convey.Convey("Replys", t, func(ctx convey.C) {
		var (
			c      = context.TODO()
			ssnID  = int64(3131)
			offset = int(1)
			limit  = int(10)
		)
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			rs, err := d.Replys(c, ssnID, offset, limit)
			ctx.Convey("Then err should be nil.rs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.SkipSo(rs, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When d.selReply.Query gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.selReply), "Query",
				func(_ *sql.Stmt, _ context.Context, _ ...interface{}) (*sql.Rows, error) {
					return nil, fmt.Errorf("d.selReply.Query Error")
				})
			defer guard.Unpatch()
			_, err := d.Replys(c, ssnID, offset, limit)
			ctx.Convey("Then err should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoWebReplys(t *testing.T) {
	convey.Convey("WebReplys", t, func(ctx convey.C) {
		var (
			c     = context.TODO()
			ssnID = int64(3131)
			mid   = int64(1313)
		)
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			rs, err := d.WebReplys(c, ssnID, mid)
			ctx.Convey("Then err should be nil.rs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.SkipSo(rs, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When d.selReplyBySid.Query gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.selReplyBySid), "Query",
				func(_ *sql.Stmt, _ context.Context, _ ...interface{}) (*sql.Rows, error) {
					return nil, fmt.Errorf("d.selReplyBySid.Query Error")
				})
			defer guard.Unpatch()
			_, err := d.WebReplys(c, ssnID, mid)
			ctx.Convey("Then err should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoReplysByMid(t *testing.T) {
	convey.Convey("ReplysByMid", t, func(ctx convey.C) {
		var (
			c      = context.TODO()
			mid    = int64(11424224)
			offset = int(1)
			limit  = int(10)
		)
		ctx.Convey("ReplysByMid", func(ctx convey.C) {
			rs, err := d.ReplysByMid(c, mid, offset, limit)
			ctx.Convey("Then err should be nil.rs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.SkipSo(rs, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When d.selReplyByMid.Query gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.selReplyByMid), "Query",
				func(_ *sql.Stmt, _ context.Context, _ ...interface{}) (*sql.Rows, error) {
					return nil, fmt.Errorf("d.selReplyByMid.Query Error")
				})
			defer guard.Unpatch()
			_, err := d.ReplysByMid(c, mid, offset, limit)
			ctx.Convey("Then err should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
