package dao

import (
	"context"
	"go-common/app/service/main/tv/internal/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoNewYstClient(t *testing.T) {
	convey.Convey("NewYstClient", t, func(ctx convey.C) {
		var (
			cli = &bm.Client{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := NewYstClient(cli)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoNewRequest(t *testing.T) {
	convey.Convey("NewRequest", t, func(ctx convey.C) {
		var (
			cli  = NewYstClient(&bm.Client{})
			uri  = ""
			data = interface{}(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			request, err := cli.NewRequest(uri, data)
			ctx.Convey("Then err should be nil.request should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(request, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoresultCode2err(t *testing.T) {
	convey.Convey("resultCode2err", t, func(ctx convey.C) {
		var (
			code = "999"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := resultCode2err(code)
			ctx.Convey("Then err should be TVIPYstSystemErr.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, ecode.TVIPYstSystemErr)
			})
		})
	})
}

func TestDaoYstOrderState(t *testing.T) {
	convey.Convey("YstOrderState", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			req = &model.YstOrderStateReq{
				SeqNo: "08265187190109103054",
				TraceNo: "CA20190109102306315819491155",
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.YstOrderState(c, req)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}