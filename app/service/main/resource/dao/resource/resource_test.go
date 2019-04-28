package resource

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	xsql "go-common/library/database/sql"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

func TestResourceResources(t *testing.T) {
	convey.Convey("Resources", t, func(ctx convey.C) {
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			res, err := d.Resources(context.Background())
			ctx.Convey("Error should be nil, res should not be empty", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeEmpty)
			})
		})
		ctx.Convey("When db.Query gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Query", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (*xsql.Rows, error) {
				return nil, fmt.Errorf("db.Query error")
			})
			defer guard.Unpatch()
			_, err := d.Resources(context.Background())
			ctx.Convey("Error should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

// Set d.close() to get reversal case
func TestResourceAssignment(t *testing.T) {
	convey.Convey("Assignment", t, func(ctx convey.C) {
		convey.Convey("When everything is correct,", func(ctx convey.C) {
			asgs, err := d.Assignment(context.Background())
			ctx.Convey("Error should be nil, asgs should not be nil(No Data)", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(asgs, convey.ShouldBeNil)
			})
		})
		convey.Convey("When set db closed", WithReopenDB(func(d *Dao) {
			d.Close()
			_, err := d.Assignment(context.Background())
			convey.Convey("Error should not be nil", func(ctx convey.C) {
				convey.So(err, convey.ShouldNotBeNil)
			})
		}))
	})
}

func TestResourceAssignmentNew(t *testing.T) {
	convey.Convey("AssignmentNew", t, func(ctx convey.C) {
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			_, err := d.AssignmentNew(context.Background())
			ctx.Convey("Error should be nil, asgs should not be empty(No Data)", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestResourceCategoryAssignment(t *testing.T) {
	convey.Convey("CategoryAssignment", t, func(ctx convey.C) {
		_, err := d.CategoryAssignment(context.Background())
		ctx.Convey("Error should be nil, res should not be empty(No Data)", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("When db.Query gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Query", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (*xsql.Rows, error) {
				return nil, fmt.Errorf("db.Query error")
			})
			defer guard.Unpatch()
			_, err := d.CategoryAssignment(context.Background())
			ctx.Convey("Error should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestResourceDefaultBanner(t *testing.T) {
	convey.Convey("DefaultBanner", t, func(ctx convey.C) {
		res, err := d.DefaultBanner(context.Background())
		ctx.Convey("Error should be nil, res should not be empty", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeEmpty)
		})
	})
}

func TestResourceIndexIcon(t *testing.T) {
	convey.Convey("IndexIcon", t, func(ctx convey.C) {
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			_, err := d.IndexIcon(context.Background())
			ctx.Convey("Error should be nil, res should not be empty(No Data)", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Convey("When db.Query gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Query", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (*xsql.Rows, error) {
				return nil, fmt.Errorf("db.Query error")
			})
			defer guard.Unpatch()
			_, err := d.IndexIcon(context.Background())
			ctx.Convey("Error should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestResourcePlayerIcon(t *testing.T) {
	convey.Convey("PlayerIcon", t, func(ctx convey.C) {
		res, err := d.PlayerIcon(context.Background())
		ctx.Convey("Error should be nil, res should not be empty", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeEmpty)
		})
	})
}

func TestResourceCmtbox(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("Cmtbox", t, func(ctx convey.C) {
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			res, err := d.Cmtbox(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When db.Query gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Query", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (*xsql.Rows, error) {
				return nil, fmt.Errorf("db.Query error")
			})
			defer guard.Unpatch()
			_, err := d.Cmtbox(context.Background())
			ctx.Convey("Error should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestResourceTxOffLine(t *testing.T) {
	var (
		tx, _ = d.BeginTran(context.TODO())
		id    = int(0)
	)
	convey.Convey("TxOffLine", t, func(ctx convey.C) {
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			row, err := d.TxOffLine(tx, id)
			ctx.Convey("Then err should be nil.row should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(row, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When tx.Exec gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(tx), "Exec", func(_ *xsql.Tx, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, fmt.Errorf("tx.Exec error")
			})
			defer guard.Unpatch()
			_, err := d.TxOffLine(tx, id)
			ctx.Convey("Then err should be not nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Rollback()
		})
	})
}

func TestResourceTxFreeApply(t *testing.T) {
	var (
		tx, _ = d.BeginTran(context.TODO())
		ids   = []string{"0", "1"}
	)
	convey.Convey("TxFreeApply", t, func(ctx convey.C) {
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			row, err := d.TxFreeApply(tx, ids)
			ctx.Convey("Then err should be nil.row should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(row, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When tx.Exec gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(tx), "Exec", func(_ *xsql.Tx, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, fmt.Errorf("tx.Exec error")
			})
			defer guard.Unpatch()
			_, err := d.TxFreeApply(tx, ids)
			ctx.Convey("Then err should be not nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Rollback()
		})
	})
}

func TestResourceTxInResourceLogger(t *testing.T) {
	var (
		tx, _   = d.BeginTran(context.TODO())
		module  = ""
		content = ""
		oid     = int(0)
	)
	convey.Convey("TxInResourceLogger", t, func(ctx convey.C) {
		ctx.Convey("When everyting is correct", func(ctx convey.C) {
			row, err := d.TxInResourceLogger(tx, module, content, oid)
			ctx.Convey("Then err should be nil.row should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(row, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When tx.Exec gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(tx), "Exec", func(_ *xsql.Tx, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, fmt.Errorf("tx.Exec error")
			})
			defer guard.Unpatch()
			_, err := d.TxInResourceLogger(tx, module, content, oid)
			ctx.Convey("Then err should be not nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Rollback()
		})
	})
}
