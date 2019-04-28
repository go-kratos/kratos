package data

import (
	"context"
	"encoding/binary"
	"reflect"
	"testing"

	"go-common/library/ecode"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
	"github.com/tsuna/gohbase/hrpc"
	hbase "go-common/library/database/hbase.v2"
)

func TestDataViewerBase(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		dt  = "dt"
		err error
	)
	convey.Convey("1", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			return nil, ecode.CreativeDataErr
		})
		defer guard.Unpatch()
		_, err = d.ViewerBase(c, mid, dt)
		ctx.Convey("1", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
	convey.Convey("2", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			res := &hrpc.Result{}
			return res, nil
		})
		defer guard.Unpatch()
		_, err = d.ViewerBase(c, mid, dt)
		ctx.Convey("2", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
	convey.Convey("3", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			bs := make([]byte, 4)
			binary.LittleEndian.PutUint32(bs, 123)
			res := &hrpc.Result{
				Cells: make([]*hrpc.Cell, 0),
			}
			res.Cells = append(res.Cells, &hrpc.Cell{
				Family:    []byte("f"),
				Qualifier: []byte("male"),
				Value:     bs,
			})
			return res, nil
		})
		defer guard.Unpatch()
		_, err = d.ViewerBase(c, mid, dt)
		ctx.Convey("3", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
	convey.Convey("4", t, func(ctx convey.C) {
		g1 := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			bs := make([]byte, 4)
			binary.LittleEndian.PutUint32(bs, 123)
			res := &hrpc.Result{
				Cells: make([]*hrpc.Cell, 0),
			}
			res.Cells = append(res.Cells, &hrpc.Cell{
				Family:    []byte("g"),
				Qualifier: []byte("female"),
				Value:     bs,
			})
			return res, nil
		})
		defer g1.Unpatch()
		_, err = d.ViewerBase(c, mid, dt)
		ctx.Convey("41", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		g2 := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			bs := make([]byte, 4)
			binary.LittleEndian.PutUint32(bs, 123)
			res := &hrpc.Result{
				Cells: make([]*hrpc.Cell, 0),
			}
			res.Cells = append(res.Cells, &hrpc.Cell{
				Family:    []byte("g"),
				Qualifier: []byte("male"),
				Value:     bs,
			})
			return res, nil
		})
		defer g2.Unpatch()
		_, err = d.ViewerBase(c, mid, dt)
		ctx.Convey("42", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDataViewerTrend(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		dt  = "dt"
		err error
	)
	convey.Convey("1", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			return nil, ecode.CreativeDataErr
		})
		defer guard.Unpatch()
		_, err = d.ViewerTrend(c, mid, dt)
		ctx.Convey("1", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
	convey.Convey("2", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			res := &hrpc.Result{}
			return res, nil
		})
		defer guard.Unpatch()
		_, err = d.ViewerTrend(c, mid, dt)
		ctx.Convey("2", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
	convey.Convey("3", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			bs := make([]byte, 4)
			binary.LittleEndian.PutUint32(bs, 123)
			res := &hrpc.Result{
				Cells: make([]*hrpc.Cell, 0),
			}
			res.Cells = append(res.Cells, &hrpc.Cell{
				Family:    []byte("fs"),
				Qualifier: []byte("male"),
				Value:     bs,
			})
			return res, nil
		})
		defer guard.Unpatch()
		_, err = d.ViewerTrend(c, mid, dt)
		ctx.Convey("3", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
	convey.Convey("4", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			bs := make([]byte, 4)
			binary.LittleEndian.PutUint32(bs, 123)
			res := &hrpc.Result{
				Cells: make([]*hrpc.Cell, 0),
			}
			res.Cells = append(res.Cells, &hrpc.Cell{
				Family:    []byte("gs"),
				Qualifier: []byte("male"),
				Value:     bs,
			})
			return res, nil
		})
		defer guard.Unpatch()
		_, err = d.ViewerTrend(c, mid, dt)
		ctx.Convey("4", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDataRelationFansDay(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
	)
	convey.Convey("RelationFansDay", t, func(ctx convey.C) {
		_, err := d.RelationFansDay(c, mid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestDataRelationFansHistory(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(0)
		month = ""
	)
	convey.Convey("RelationFansHistory", t, func(ctx convey.C) {
		_, err := d.RelationFansHistory(c, mid, month)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestDataRelationFansMonth(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
	)
	convey.Convey("RelationFansMonth", t, func(ctx convey.C) {
		_, err := d.RelationFansMonth(c, mid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestDataViewerActionHour(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		dt  = "dt"
		err error
	)
	convey.Convey("1", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			return nil, ecode.CreativeDataErr
		})
		defer guard.Unpatch()
		_, err = d.ViewerActionHour(c, mid, dt)
		ctx.Convey("1", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
	convey.Convey("2", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			res := &hrpc.Result{}
			return res, nil
		})
		defer guard.Unpatch()
		_, err = d.ViewerActionHour(c, mid, dt)
		ctx.Convey("2", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
	convey.Convey("3", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			bs := make([]byte, 4)
			binary.LittleEndian.PutUint32(bs, 123)
			res := &hrpc.Result{
				Cells: make([]*hrpc.Cell, 0),
			}
			res.Cells = append(res.Cells, &hrpc.Cell{
				Family:    []byte("fs"),
				Qualifier: []byte("male"),
				Value:     bs,
			})
			return res, nil
		})
		defer guard.Unpatch()
		_, err = d.ViewerActionHour(c, mid, dt)
		ctx.Convey("3", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
	convey.Convey("4", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			bs := make([]byte, 4)
			binary.LittleEndian.PutUint32(bs, 123)
			res := &hrpc.Result{
				Cells: make([]*hrpc.Cell, 0),
			}
			res.Cells = append(res.Cells, &hrpc.Cell{
				Family:    []byte("gs"),
				Qualifier: []byte("male"),
				Value:     bs,
			})
			return res, nil
		})
		defer guard.Unpatch()
		_, err = d.ViewerActionHour(c, mid, dt)
		ctx.Convey("4", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDataUpIncr(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		ty  = int8(2)
		now = ""
		err error
	)
	convey.Convey("1", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			return nil, ecode.CreativeDataErr
		})
		defer guard.Unpatch()
		_, err = d.UpIncr(c, mid, ty, now)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
	convey.Convey("2", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			res := &hrpc.Result{}
			return res, nil
		})
		defer guard.Unpatch()
		_, err = d.UpIncr(c, mid, ty, now)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
	convey.Convey("3", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			bs := make([]byte, 4)
			binary.LittleEndian.PutUint32(bs, 123)
			res := &hrpc.Result{
				Cells: make([]*hrpc.Cell, 0),
			}
			res.Cells = append(res.Cells, &hrpc.Cell{
				Family:    []byte("u"),
				Qualifier: []byte("male"),
				Value:     bs,
			})
			return res, nil
		})
		defer guard.Unpatch()
		_, err = d.UpIncr(c, mid, ty, now)
		ctx.Convey("3", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDataThirtyDayArchive(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		ty  = int8(2)
		err error
	)
	convey.Convey("1", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			return nil, ecode.CreativeDataErr
		})
		defer guard.Unpatch()
		_, err = d.ThirtyDayArchive(c, mid, ty)
		ctx.Convey("1", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
	convey.Convey("2", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			res := &hrpc.Result{}
			return res, nil
		})
		defer guard.Unpatch()
		_, err = d.ThirtyDayArchive(c, mid, ty)
		ctx.Convey("2", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
	convey.Convey("3", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			res := &hrpc.Result{
				Cells: make([]*hrpc.Cell, 0),
			}
			res.Cells = append(res.Cells, &hrpc.Cell{
				Family:    []byte("u"),
				Qualifier: []byte("20181111"),
				Value:     []byte("200"),
			})
			return res, nil
		})
		defer guard.Unpatch()
		_, err = d.ThirtyDayArchive(c, mid, ty)
		ctx.Convey("3", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDataparseKeyValue(t *testing.T) {
	var (
		k = ""
		v = ""
	)
	convey.Convey("parseKeyValue", t, func(ctx convey.C) {
		_, _, err := parseKeyValue(k, v)
		ctx.Convey("Then err should be nil.timestamp,value should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}
