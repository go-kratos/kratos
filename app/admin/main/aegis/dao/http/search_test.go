package http

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"go-common/app/admin/main/aegis/model"
)

var arg = &model.SearchParams{
	BusinessID: 1,
	//CtimeFrom:  "2018-01-01 01:01:01",
	//CtimeTo:    "2019-11-11 11:11:11",
	OID:     []string{"206137981869783258", "206137981871880409", "206129056929839317", "206117679558326449"},
	KeyWord: "分享图片",
	FlowID:  -12345,
	Extra1:  "1,2,3,4,5",
	Mid:     -12345,
	State:   0,
}

func TestHttpResourceES(t *testing.T) {
	convey.Convey("ResourceES", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ResourceES(c, arg)
			t.Logf("res(%+v)", res)
			ctx.Convey("Then err should be nil.sres should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestUpsertES(t *testing.T) {
	convey.Convey("UpsertES", t, func(ctx convey.C) {
		rsc := []*model.UpsertItem{
			{
				ID:    47,
				State: -3,
			},
			nil,
			{
				State: 0,
			},
		}
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			httpMock("POST", d.c.Host.Manager+_upsertES).Reply(200).JSON(`{"code":0,"message":"msg"}`)
			err := d.UpsertES(context.Background(), rsc)
			ctx.Convey("No return values", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
