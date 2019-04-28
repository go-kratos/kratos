package data

import (
	"context"
	"go-common/library/ecode"
	"reflect"
	"testing"

	hbase "go-common/library/database/hbase.v2"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
	"github.com/tsuna/gohbase/hrpc"
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

func TestDataVideoQuitPoints(t *testing.T) {
	var (
		c   = context.TODO()
		cid = int64(1)
	)
	convey.Convey("VideoQuitPoints1", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			res := &hrpc.Result{
				Cells: make([]*hrpc.Cell, 0),
			}
			res.Cells = append(res.Cells, &hrpc.Cell{
				Family:    HBaseFamilyPlat,
				Qualifier: HBaseColumnShare,
			})
			return res, nil
		})
		defer guard.Unpatch()
		_, err := d.VideoQuitPoints(c, cid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
	convey.Convey("VideoQuitPoints2", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			return nil, ecode.CreativeDataErr
		})
		defer guard.Unpatch()
		_, err := d.VideoQuitPoints(c, cid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestDataArchiveStat(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(1)
		err error
	)
	convey.Convey("1", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			return nil, ecode.CreativeDataErr
		})
		defer guard.Unpatch()
		_, err = d.ArchiveStat(c, aid)
		ctx.Convey("Then err should be nil.stat should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldEqual, ecode.CreativeDataErr)
		})
	})
	convey.Convey("2", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			res := &hrpc.Result{}
			return res, nil
		})
		defer guard.Unpatch()
		_, err = d.ArchiveStat(c, aid)
		ctx.Convey("Then err should be nil.stat should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
	convey.Convey("3", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			return nil, nil
		})
		defer guard.Unpatch()
		_, err = d.ArchiveStat(c, aid)
		ctx.Convey("Then err should be nil.stat should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
	convey.Convey("4", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			res := &hrpc.Result{
				Cells: make([]*hrpc.Cell, 0),
			}
			res.Cells = append(res.Cells, &hrpc.Cell{
				Family:    HBaseFamilyPlat,
				Qualifier: HBaseColumnWebPC,
			})
			return res, nil
		})
		defer guard.Unpatch()
		_, err = d.ArchiveStat(c, aid)
		ctx.Convey("Then err should be nil.stat should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
	convey.Convey("5", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			res := &hrpc.Result{
				Cells: make([]*hrpc.Cell, 0),
			}
			res.Cells = append(res.Cells, &hrpc.Cell{
				Family:    HBaseFamilyPlat,
				Qualifier: HBaseColumnWebH5,
			})
			return res, nil
		})
		defer guard.Unpatch()
		_, err = d.ArchiveStat(c, aid)
		ctx.Convey("Then err should be nil.stat should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
	convey.Convey("6", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			res := &hrpc.Result{
				Cells: make([]*hrpc.Cell, 0),
			}
			res.Cells = append(res.Cells, &hrpc.Cell{
				Family:    HBaseFamilyPlat,
				Qualifier: HBaseColumnShare,
			})
			return res, nil
		})
		defer guard.Unpatch()
		_, err = d.ArchiveStat(c, aid)
		ctx.Convey("Then err should be nil.stat should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDataArchiveArea(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(1)
	)
	convey.Convey("ArchiveArea1", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			return nil, ecode.CreativeDataErr
		})
		defer guard.Unpatch()
		_, err := d.ArchiveArea(c, aid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldEqual, ecode.CreativeDataErr)
		})
	})
	convey.Convey("ArchiveArea2", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			res := &hrpc.Result{
				Cells: make([]*hrpc.Cell, 0),
			}
			res.Cells = append(res.Cells, &hrpc.Cell{
				Family:    HBaseFamilyPlat,
				Qualifier: HBaseColumnShare,
			})
			return res, nil
		})
		defer guard.Unpatch()
		_, err := d.ArchiveArea(c, aid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDataBaseUpStat(t *testing.T) {
	var (
		c    = context.TODO()
		mid  = int64(0)
		date = ""
	)
	convey.Convey("BaseUpStat", t, func(ctx convey.C) {
		_, err := d.BaseUpStat(c, mid, date)
		ctx.Convey("Then err should be nil.stat should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldEqual, ecode.CreativeDataErr)
		})
	})
}

func TestDataUpArchiveStatQuery(t *testing.T) {
	var (
		c    = context.TODO()
		mid  = int64(1)
		date = ""
	)
	convey.Convey("UpArchiveStatQuery1", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			return nil, ecode.CreativeDataErr
		})
		defer guard.Unpatch()
		_, err := d.UpArchiveStatQuery(c, mid, date)
		ctx.Convey("Then err should be nil.stat should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
	convey.Convey("UpArchiveStatQuery2", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			res := &hrpc.Result{}
			return res, nil
		})
		defer guard.Unpatch()
		_, err := d.UpArchiveStatQuery(c, mid, date)
		ctx.Convey("Then err should be nil.stat should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
