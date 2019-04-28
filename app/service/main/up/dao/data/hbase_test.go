package data

import (
	"context"
	"github.com/bouk/monkey"
	"github.com/tsuna/gohbase/hrpc"
	"go-common/library/database/hbase.v2"
	"reflect"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDatahbaseMd5Key(t *testing.T) {
	var (
		aid = int64(0)
	)
	convey.Convey("hbaseMd5Key", t, func(ctx convey.C) {
		p1 := hbaseMd5Key(aid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDataBaseUpStat(t *testing.T) {
	var (
		c    = context.TODO()
		mid  = int64(0)
		date = ""
	)
	// Hbase never ok
	convey.Convey("BaseUpStat", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "Get", func(_ *hbase.Client, _ context.Context, _ []byte, _ []byte, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			cells := make([]*hrpc.Cell, 5)
			for i := range cells {
				cell := new(hrpc.Cell)
				cell.Value = []byte("test")
				cells[i] = cell
			}
			res := &hrpc.Result{
				Cells: cells,
			}
			return res, nil
		})
		defer guard.Unpatch()
		stat, err := d.BaseUpStat(c, mid, date)
		ctx.Convey("Then err should be nil.stat should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(stat, convey.ShouldNotBeNil)
		})
	})
}
