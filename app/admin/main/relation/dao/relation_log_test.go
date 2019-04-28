package dao

import (
	"context"

	"fmt"
	"reflect"
	"testing"
	"time"

	"go-common/library/database/hbase.v2"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
	"github.com/tsuna/gohbase/hrpc"
)

func TestDaoreverse(t *testing.T) {
	convey.Convey("reverse", t, func(convCtx convey.C) {
		var (
			s = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := reverse(s)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaorpad(t *testing.T) {
	convey.Convey("rpad", t, func(convCtx convey.C) {
		var (
			s = ""
			l = int(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := rpad(s, l)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaologKey(t *testing.T) {
	convey.Convey("logKey", t, func(convCtx convey.C) {
		var (
			mid = int64(0)
			fid = int64(0)
			ts  = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := logKey(mid, fid, ts)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRelationLogs(t *testing.T) {
	convey.Convey("RelationLogs", t, func(convCtx convey.C) {
		var (
			ctx  = context.Background()
			mid  = int64(0)
			fid  = int64(0)
			from = time.Now()
			to   = time.Now()
		)
		convCtx.Convey("RelationLogs failed", func(convCtx convey.C) {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "ScanRangeStr", func(_ *hbase.Client, _ context.Context, _ string, _ string, _ string, _ ...func(hrpc.Call) error) (hrpc.Scanner, error) {
				return nil, fmt.Errorf("hbase scan err")
			})
			defer monkey.UnpatchAll()
			p1, err := d.RelationLogs(ctx, mid, fid, from, to)
			convCtx.Convey("Then err should be nil.p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldNotBeNil)
				convCtx.So(p1, convey.ShouldBeNil)
			})
		})

	})
}
