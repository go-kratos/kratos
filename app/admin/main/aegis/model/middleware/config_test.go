package middleware

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"go-common/app/admin/main/aegis/model"
	"go-common/app/admin/main/aegis/model/resource"
)

var (
	d = &Aggregate{
		Hitn:  "extra1",
		Hitv:  "2,3,4,5,6,7",
		Mapn:  "extra1",
		Mapv:  "2",
		Order: 1,
	}
	ds = MiddleAggregate{
		Cfg:    []*Aggregate{d},
		Encode: true,
	}
)

func TestMiddleware_AggregateProcess(t *testing.T) {
	var (
		data = &model.AuditInfo{
			Resource: &resource.Res{Extra1: 6},
		}
		encode bool = true
		data2       = &model.SearchParams{
			Extra1: "2",
		}
	)
	convey.Convey("AggregateProcess", t, func(ctx convey.C) {
		d.Process(data, encode)   //将extra1=6替换成extra1=2
		d.Process(data2, !encode) //将extra1=2替换成extra1=2,3,4,5,6,7
		ctx.Convey("extra1 equal", func(ctx convey.C) {
			ctx.So(fmt.Sprintf("%d", data.Resource.Extra1), convey.ShouldEqual, d.Mapv)
			ctx.So(data2.Extra1, convey.ShouldEqual, d.Hitv)
		})
	})
}

func TestMiddlewaregetFieldByName(t *testing.T) {
	var (
		v = reflect.ValueOf(&model.AuditInfo{
			Resource: &resource.Res{Extra1: 6},
		})
		name = "extra1"
	)
	convey.Convey("getFieldByName", t, func(ctx convey.C) {
		res, ok := getFieldByName(v, name)
		ctx.Convey("Then res,ok should not be nil.", func(ctx convey.C) {
			ctx.So(ok, convey.ShouldEqual, true)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestMiddlewareProcess(t *testing.T) {
	var (
		data = &model.AuditInfo{
			Resource: &resource.Res{Extra1: 6},
		}
	)
	convey.Convey("Process", t, func(ctx convey.C) {
		ds.Encode = true
		ds.Process(data)
		ctx.Convey("No return values", func(ctx convey.C) {
			//将extra1=6替换成extra1=2
			ctx.So(fmt.Sprintf("%d", data.Resource.Extra1), convey.ShouldEqual, d.Mapv)
		})
	})
}

func TestMiddlewareLen(t *testing.T) {
	convey.Convey("Len", t, func(ctx convey.C) {
		p1 := AggregateArr(ds.Cfg).Len()
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldEqual, 1)
		})
	})
}

func TestMiddlewareLess(t *testing.T) {
	var (
		i = int(0)
		j = int(0)
	)
	convey.Convey("Less", t, func(ctx convey.C) {
		p1 := AggregateArr(ds.Cfg).Less(i, j)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldEqual, false)
		})
	})
}

func TestMiddlewareSwap(t *testing.T) {
	var (
		i = int(0)
		j = int(0)
	)
	convey.Convey("Swap", t, func(ctx convey.C) {
		AggregateArr(ds.Cfg).Swap(i, j)
		ctx.Convey("No return values", func(ctx convey.C) {
		})
	})
}
