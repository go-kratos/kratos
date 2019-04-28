package block

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	model "go-common/app/job/main/member/model/block"
	xsql "go-common/library/database/sql"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

func TestBlockhistoryIdx(t *testing.T) {
	convey.Convey("historyIdx", t, func(convCtx convey.C) {
		var (
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := historyIdx(mid)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestBlockUserStatusList(t *testing.T) {
	convey.Convey("UserStatusList", t, func(convCtx convey.C) {
		var (
			c       = context.Background()
			status  model.BlockStatus
			startID = int64(0)
			limit   = int(10)
		)
		rows, _ := d.db.Query(c, _userStatusList, status, startID, limit)
		convCtx.Convey("UserStatusList success", func(convCtx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(rows), "Scan", func(_ *xsql.Rows, _ ...interface{}) error {
				return nil
			})
			defer guard.Unpatch()
			maxID, mids, err := d.UserStatusList(c, status, startID, limit)
			convCtx.Convey("Then err should be nil.maxID,mids should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(mids, convey.ShouldNotBeNil)
				convCtx.So(maxID, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestBlockUserLastHistory(t *testing.T) {
	convey.Convey("UserLastHistory", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		convey.Convey("UserLastHistory success", func() {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.db.QueryRow(context.TODO(), fmt.Sprintf(_userLastHistory, historyIdx(mid)), 1)), "Scan", func(_ *xsql.Row, _ ...interface{}) error {
				return nil
			})
			defer monkey.UnpatchAll()
			his, err := d.UserLastHistory(c, mid)
			convCtx.So(err, convey.ShouldBeNil)
			convCtx.So(his, convey.ShouldNotBeNil)
		})
		convey.Convey("UserLastHistory err", func() {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.db.QueryRow(context.TODO(), fmt.Sprintf(_userLastHistory, historyIdx(mid)), 1)), "Scan", func(_ *xsql.Row, _ ...interface{}) error {
				return fmt.Errorf("row.Scan error")
			})
			defer monkey.UnpatchAll()
			his, err := d.UserLastHistory(c, mid)
			convCtx.So(err, convey.ShouldNotBeNil)
			convCtx.So(his, convey.ShouldNotBeNil)
		})
		convey.Convey("UserLastHistory no record", func() {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.db.QueryRow(context.TODO(), fmt.Sprintf(_userLastHistory, historyIdx(mid)), 1)), "Scan", func(_ *xsql.Row, _ ...interface{}) error {
				return sql.ErrNoRows
			})
			defer monkey.UnpatchAll()
			his, err := d.UserLastHistory(c, mid)
			convCtx.So(err, convey.ShouldBeNil)
			convCtx.So(his, convey.ShouldBeNil)
		})
	})
}

func TestBlockUserExtra(t *testing.T) {
	convey.Convey("UserExtra", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		convey.Convey("UserExtra success", func() {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.db.QueryRow(context.TODO(), _userExtra, 1)), "Scan", func(_ *xsql.Row, _ ...interface{}) error {
				return nil
			})
			defer monkey.UnpatchAll()
			ex, err := d.UserExtra(c, mid)
			convCtx.So(err, convey.ShouldBeNil)
			convCtx.So(ex, convey.ShouldNotBeNil)
		})
		convey.Convey("UserExtra err", func() {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.db.QueryRow(context.TODO(), _userExtra, 1)), "Scan", func(_ *xsql.Row, _ ...interface{}) error {
				return fmt.Errorf("row.Scan error")
			})
			defer monkey.UnpatchAll()
			ex, err := d.UserExtra(c, mid)
			convCtx.So(err, convey.ShouldNotBeNil)
			convCtx.So(ex, convey.ShouldNotBeNil)
		})
		convey.Convey("UserExtra no record", func() {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.db.QueryRow(context.TODO(), _userExtra, 1)), "Scan", func(_ *xsql.Row, _ ...interface{}) error {
				return sql.ErrNoRows
			})
			defer monkey.UnpatchAll()
			ex, err := d.UserExtra(c, mid)
			convCtx.So(err, convey.ShouldBeNil)
			convCtx.So(ex, convey.ShouldBeNil)
		})

	})
}

func TestBlockTxUpsertUser(t *testing.T) {
	convey.Convey("TxUpsertUser", t, func(convCtx convey.C) {
		var (
			c      = context.Background()
			tx, _  = d.BeginTX(c)
			mid    = int64(0)
			status model.BlockStatus
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			count, err := d.TxUpsertUser(c, tx, mid, status)
			convCtx.Convey("Then err should be nil.count should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(count, convey.ShouldNotBeNil)
			})
		})
		convCtx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestBlockInsertExtra(t *testing.T) {
	convey.Convey("InsertExtra", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			ex = &model.DBExtra{}
		)
		convCtx.Convey("InsertExtra success", func(convCtx convey.C) {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Exec", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, nil
			})
			defer monkey.UnpatchAll()
			err := d.InsertExtra(c, ex)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})

		convCtx.Convey("InsertExtra err", func(convCtx convey.C) {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Exec", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, fmt.Errorf("insert err")
			})
			defer monkey.UnpatchAll()
			err := d.InsertExtra(c, ex)
			convCtx.Convey("Then err should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestBlockTxUpsertExtra(t *testing.T) {
	convey.Convey("TxUpsertExtra", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTX(c)
			ex    = &model.DBExtra{}
		)

		convCtx.Convey("TxUpsertExtra success", func(convCtx convey.C) {
			monkey.PatchInstanceMethod(reflect.TypeOf(tx), "Exec", func(_ *xsql.Tx, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, nil
			})
			defer monkey.UnpatchAll()
			err := d.TxUpsertExtra(c, tx, ex)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
		convCtx.Convey("TxUpsertExtra err", func(convCtx convey.C) {
			monkey.PatchInstanceMethod(reflect.TypeOf(tx), "Exec", func(_ *xsql.Tx, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, fmt.Errorf("insert err")
			})
			defer monkey.UnpatchAll()
			err := d.TxUpsertExtra(c, tx, ex)
			convCtx.Convey("Then err should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldNotBeNil)
			})
		})

		convCtx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestBlockTxInsertHistory(t *testing.T) {
	convey.Convey("TxInsertHistory", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTX(c)
			h     = &model.DBHistory{}
		)
		convCtx.Convey("TxUpsertExtra success", func(convCtx convey.C) {
			monkey.PatchInstanceMethod(reflect.TypeOf(tx), "Exec", func(_ *xsql.Tx, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, nil
			})
			defer monkey.UnpatchAll()
			err := d.TxInsertHistory(c, tx, h)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})

		convCtx.Convey("TxUpsertExtra err", func(convCtx convey.C) {
			monkey.PatchInstanceMethod(reflect.TypeOf(tx), "Exec", func(_ *xsql.Tx, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, fmt.Errorf("insert err")
			})
			defer monkey.UnpatchAll()
			err := d.TxInsertHistory(c, tx, h)
			convCtx.Convey("Then err should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldNotBeNil)
			})
		})

		convCtx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestBlockUpsertAddBlockCount(t *testing.T) {
	convey.Convey("UpsertAddBlockCount", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		convCtx.Convey("UpsertAddBlockCount success", func(convCtx convey.C) {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Exec", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, nil
			})
			defer monkey.UnpatchAll()
			err := d.UpsertAddBlockCount(c, mid)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
		convCtx.Convey("InsertExtra err", func(convCtx convey.C) {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Exec", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, fmt.Errorf("insert err")
			})
			defer monkey.UnpatchAll()
			err := d.UpsertAddBlockCount(c, mid)
			convCtx.Convey("Then err should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldNotBeNil)
			})
		})

	})
}
