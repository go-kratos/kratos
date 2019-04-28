package data

import (
	"context"
	"encoding/binary"
	"go-common/app/interface/main/creative/model/data"
	"go-common/library/ecode"
	"reflect"
	"testing"

	hbase "go-common/library/database/hbase.v2"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
	"github.com/tsuna/gohbase/hrpc"
)

func TestDataUpFansAnalysisForApp(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		ty  = int(0)
		err error
		res *data.AppFan
	)
	convey.Convey("1", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			return nil, ecode.CreativeDataErr
		})
		defer guard.Unpatch()
		res, err = d.UpFansAnalysisForApp(c, mid, ty)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldEqual, ecode.CreativeDataErr)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
	convey.Convey("2", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			res := &hrpc.Result{}
			return res, nil
		})
		defer guard.Unpatch()
		res, err = d.UpFansAnalysisForApp(c, mid, ty)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldBeNil)
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
				Qualifier: []byte("all"),
				Value:     bs,
			})
			return res, nil
		})
		defer guard.Unpatch()
		res, err = d.UpFansAnalysisForApp(c, mid, ty)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
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
				Family:    []byte("t"),
				Qualifier: []byte("dr"),
				Value:     bs,
			})
			return res, nil
		})
		defer guard.Unpatch()
		res, err = d.UpFansAnalysisForApp(c, mid, ty)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDataViewerArea(t *testing.T) {
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
		_, err = d.ViewerArea(c, mid, dt)
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
		_, err = d.ViewerArea(c, mid, dt)
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
		_, err = d.ViewerArea(c, mid, dt)
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
				Family:    []byte("g"),
				Qualifier: []byte("male"),
				Value:     bs,
			})
			return res, nil
		})
		defer guard.Unpatch()
		_, err = d.ViewerArea(c, mid, dt)
		ctx.Convey("4", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
