package dao

import (
	"context"
	"go-common/app/job/main/passport/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"net/url"
	"reflect"
	"testing"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

func TestDao_SetToken(t *testing.T) {
	convey.Convey("SetToken", t, func(ctx convey.C) {
		token := &model.Token{
			Mid:   88888970,
			Token: "foo",
		}
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SetToken(context.TODO(), token)
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDao_DelCache(t *testing.T) {
	convey.Convey("DelCache", t, func(ctx convey.C) {
		token := "foo"
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelCache(context.TODO(), token)
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDao_NotifyGame(t *testing.T) {
	convey.Convey("NotifyGame", t, func(ctx convey.C) {
		var (
			mid    = &model.AccessInfo{}
			action = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			mock := monkey.PatchInstanceMethod(reflect.TypeOf(d.gameClient), "Get", func(d *bm.Client, _ context.Context, _, _ string, _ url.Values, _ interface{}) error {
				return nil
			})
			defer mock.Unpatch()
			err := d.NotifyGame(mid, action)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
			mock2 := monkey.PatchInstanceMethod(reflect.TypeOf(d.gameClient), "Get", func(d *bm.Client, _ context.Context, _, _ string, _ url.Values, _ interface{}) error {
				return ecode.Int(500)
			})
			defer mock2.Unpatch()
			err = d.NotifyGame(mid, action)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})

}
