package archive

import (
	"context"
	"database/sql"
	"fmt"
	"go-common/app/service/main/videoup/model/archive"
	xsql "go-common/library/database/sql"
	"go-common/library/time"
	"reflect"
	"testing"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

func TestArchiveArcReport(t *testing.T) {
	var (
		c   = context.Background()
		aid = int64(2333)
		mid = int64(23333)
	)
	convey.Convey("ArcReport", t, func(ctx convey.C) {
		_, err := d.ArcReport(c, aid, mid)
		ctx.Convey("Then err should be nil.aa should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
func TestTxAddRelation(t *testing.T) {
	var (
		c     = context.Background()
		err   error
		tx, _ = d.BeginTran(c)
		v     = &archive.Video{
			Aid:   int64(10110817),
			Cid:   int64(10134702),
			Title: "iamtitle",
			Index: 1,
		}
	)
	convey.Convey("TestTxAddRelation", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(tx),
			"Exec",
			func(_ *xsql.Tx, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, fmt.Errorf("tx.Exec error")
			})
		defer guard.Unpatch()
		_, err = d.TxAddRelation(tx, v)
		ctx.Convey("TestArchivePOIAdd.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestTxUpForbid(t *testing.T) {
	var (
		c     = context.Background()
		err   error
		aid   = int64(10110817)
		fid   = int64(1)
		tx, _ = d.BeginTran(c)
	)
	convey.Convey("TestTxUpForbid", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(tx),
			"Exec",
			func(_ *xsql.Tx, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, fmt.Errorf("tx.Exec error")
			})
		defer guard.Unpatch()
		_, err = d.TxUpForbid(tx, aid, fid)
		ctx.Convey("TxUpForbid.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}
func TestTxUpForbidAttr(t *testing.T) {
	var (
		c     = context.Background()
		err   error
		tx, _ = d.BeginTran(c)
		af    = &archive.ForbidAttr{}
	)
	convey.Convey("TestTxUpForbidAttr", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(tx),
			"Exec",
			func(_ *xsql.Tx, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, fmt.Errorf("tx.Exec error")
			})
		defer guard.Unpatch()
		_, err = d.TxUpForbidAttr(tx, af)
		ctx.Convey("TxUpForbidAttr.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestTxAddDelay(t *testing.T) {
	var (
		c     = context.Background()
		err   error
		tx, _ = d.BeginTran(c)
		aid   = int64(2333)
		mid   = int64(23333)
		state = int8(1)
		tp    = int8(2)
		dtime time.Time
	)
	convey.Convey("TestTxAddDelay", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(tx),
			"Exec",
			func(_ *xsql.Tx, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, fmt.Errorf("tx.Exec error")
			})
		defer guard.Unpatch()
		_, err = d.TxAddDelay(tx, mid, aid, state, tp, dtime)
		ctx.Convey("TxAddDelay.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}
