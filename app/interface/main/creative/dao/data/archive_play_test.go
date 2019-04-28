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

func TestDatasourceOtherMerge(t *testing.T) {
	var (
		v = ""
	)
	convey.Convey("sourceOtherMerge", t, func(ctx convey.C) {
		p1 := sourceOtherMerge(v)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDatareverseString(t *testing.T) {
	var (
		s = ""
	)
	convey.Convey("reverseString", t, func(ctx convey.C) {
		p1 := reverseString(s)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDatafansRowKey(t *testing.T) {
	var (
		id = int64(0)
		ty = int(0)
	)
	convey.Convey("fansRowKey", t, func(ctx convey.C) {
		p1 := fansRowKey(id, ty)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDataplaySourceKey(t *testing.T) {
	var (
		id = int64(0)
	)
	convey.Convey("playSourceKey", t, func(ctx convey.C) {
		p1 := playSourceKey(id)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDataarcPlayKey(t *testing.T) {
	var (
		id = int64(0)
	)
	convey.Convey("arcPlayKey", t, func(ctx convey.C) {
		p1 := arcPlayKey(id)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDataarcQueryKey(t *testing.T) {
	var (
		id = int64(1)
		dt = ""
		cp = int(0)
	)
	convey.Convey("arcQueryKey", t, func(ctx convey.C) {
		p1 := arcQueryKey(id, dt, cp)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldEqual, p1)
		})
	})
}

func TestDataupFansMedalRowKey(t *testing.T) {
	var (
		id = int64(1)
		ty = int(1)
	)
	convey.Convey("upFansMedalRowKey", t, func(ctx convey.C) {
		p1 := upFansMedalRowKey(id, ty)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldEqual, p1)
		})
	})
}

func TestDatabyteToInt32(t *testing.T) {
	var (
		b = []byte{255, 255, 255, 249}
	)
	convey.Convey("byteToInt32", t, func(ctx convey.C) {
		p1 := byteToInt32(b)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldEqual, -7)
		})
	})
}

func TestDataUpFansAnalysisForWeb(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		ty  = int(0)
		err error
		res *data.WebFan
	)
	convey.Convey("1", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.hbase), "GetStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ ...func(hrpc.Call) error) (*hrpc.Result, error) {
			return nil, ecode.CreativeDataErr
		})
		defer guard.Unpatch()
		res, err = d.UpFansAnalysisForWeb(c, mid, ty)
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
		res, err = d.UpFansAnalysisForWeb(c, mid, ty)
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
		res, err = d.UpFansAnalysisForWeb(c, mid, ty)
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
		res, err = d.UpFansAnalysisForWeb(c, mid, ty)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDataUpPlaySourceAnalysis(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
	)
	convey.Convey("UpPlaySourceAnalysis", t, func(ctx convey.C) {
		res, err := d.UpPlaySourceAnalysis(c, mid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldEqual, ecode.CreativeDataErr)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestDataUpArcPlayAnalysis(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(0)
	)
	convey.Convey("UpArcPlayAnalysis", t, func(ctx convey.C) {
		res, err := d.UpArcPlayAnalysis(c, aid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldEqual, ecode.CreativeDataErr)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestDataUpArcQuery(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		dt  = ""
		cp  = int(0)
	)
	convey.Convey("UpArcQuery", t, func(ctx convey.C) {
		res, err := d.UpArcQuery(c, mid, dt, cp)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldEqual, ecode.CreativeDataErr)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestDataUpFansMedal(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
	)
	convey.Convey("UpFansMedal", t, func(ctx convey.C) {
		res, err := d.UpFansMedal(c, mid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldEqual, ecode.CreativeDataErr)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}
