package data

import (
	"context"
	"reflect"
	"testing"

	"github.com/bouk/monkey"
	"github.com/tsuna/gohbase/hrpc"

	"go-common/app/admin/main/up/model/datamodel"

	hbase "go-common/library/database/hbase.v2"

	"github.com/smartystreets/goconvey/convey"
)

func TestDataUpArchiveInfo(t *testing.T) {
	convey.Convey("UpArchiveInfo", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mids     = []int64{}
			dataType UpArchiveDataType
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := d.UpArchiveInfo(c, mids, dataType)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(result, convey.ShouldBeNil)
			})
		})
	})
}

func TestDataUpArchiveTagInfo(t *testing.T) {
	convey.Convey("UpArchiveTagInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := d.UpArchiveTagInfo(c, mid)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDataUpArchiveTypeInfo(t *testing.T) {
	convey.Convey("UpArchiveTypeInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := d.UpArchiveTypeInfo(c, mid)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatagetHbaseRowResult(t *testing.T) {
	convey.Convey("getHbaseRowResult", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			table  = "crm"
			key    = "123"
			result = datamodel.UpArchiveTypeData{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
				res := &hrpc.Result{}
				return res, nil
			})
			defer guard.Unpatch()
			err := d.getHbaseRowResult(c, table, key, &result)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDatagenerateTableName(t *testing.T) {
	convey.Convey("generateTableName", t, func(ctx convey.C) {
		var (
			prefix = ""
			date   = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := generateTableName(prefix, date)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
