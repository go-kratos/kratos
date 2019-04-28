package data

import (
	"context"
	"reflect"
	"testing"

	"go-common/library/ecode"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
	"github.com/tsuna/gohbase/hrpc"
	hbase "go-common/library/database/hbase.v2"
)

func TestArtThirtyDay(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(27515256)
		tp  = byte(1)
	)
	convey.Convey("ArtThirtyDay1", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			return nil, ecode.CreativeDataErr
		})
		defer guard.Unpatch()
		_, err := d.ArtThirtyDay(c, mid, tp)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
	convey.Convey("ArtThirtyDay2", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			res := &hrpc.Result{}
			return res, nil
		})
		defer guard.Unpatch()
		_, err := d.ArtThirtyDay(c, mid, tp)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestArtRank(t *testing.T) {
	var (
		c    = context.TODO()
		mid  = int64(27515256)
		tp   = byte(1)
		date = ""
	)
	convey.Convey("TestArtRank1", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			return nil, ecode.CreativeDataErr
		})
		defer guard.Unpatch()
		_, err := d.ArtRank(c, mid, tp, date)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
	convey.Convey("TestArtRank2", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			res := &hrpc.Result{}
			return res, nil
		})
		defer guard.Unpatch()
		_, err := d.ArtRank(c, mid, tp, date)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestReadAnalysis(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(27515256)
	)

	convey.Convey("TestReadAnalysis1", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			return nil, ecode.CreativeDataErr
		})
		defer guard.Unpatch()
		_, err := d.ReadAnalysis(c, mid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})

	})
	convey.Convey("TestReadAnalysis2", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			res := &hrpc.Result{}
			return res, nil
		})
		defer guard.Unpatch()
		_, err := d.ReadAnalysis(c, mid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
