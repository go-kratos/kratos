package material

import (
	"context"
	"fmt"
	xsql "go-common/library/database/sql"
	"reflect"
	"testing"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

func TestFilterCategoryBind(t *testing.T) {
	var (
		c  = context.TODO()
		tp = int8(7)
	)
	convey.Convey("CategoryBind", t, func(ctx convey.C) {
		res, err := d.CategoryBind(c, tp)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestMaterialSubtitles(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("Subtitles", t, func(ctx convey.C) {
		res, err := d.Basic(c)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestMaterialFilters(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("Filters", t, func(ctx convey.C) {
		res, resMap, err := d.Filters(c)
		ctx.Convey("Then err should be nil.res,resMap should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(resMap, convey.ShouldNotBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestMaterialVstickers(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("Vstickers", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Query", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (*xsql.Rows, error) {
			return nil, fmt.Errorf("db.Query error")
		})
		defer guard.Unpatch()
		res, resMap, err := d.Vstickers(c)
		ctx.Convey("Then err should be nil.res,resMap should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(len(resMap), convey.ShouldEqual, 0)
			ctx.So(len(res), convey.ShouldEqual, 0)
		})
	})
}

func TestCooperates(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("Cooperates", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Query", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (*xsql.Rows, error) {
			return nil, fmt.Errorf("db.Query error")
		})
		defer guard.Unpatch()
		res, resMap, err := d.Cooperates(c)
		ctx.Convey("Then err should be nil.res,resMap should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(len(resMap), convey.ShouldEqual, 0)
			ctx.So(len(res), convey.ShouldEqual, 0)
		})
	})
}

func TestBasic(t *testing.T) {
	var (
		err error
		c   = context.TODO()
		res = make(map[string]interface{})
	)
	convey.Convey("Basic", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Query", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (*xsql.Rows, error) {
			return nil, fmt.Errorf("db.Query error")
		})
		defer guard.Unpatch()
		res, err = d.Basic(c)
		ctx.Convey("Then err should be nil.res,resMap should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(len(res), convey.ShouldEqual, 0)
		})
	})
}
