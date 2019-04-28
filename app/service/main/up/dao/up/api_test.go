package up

import (
	"context"
	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
	httpx "go-common/library/net/http/blademaster"
	"net/url"
	"reflect"
	"testing"
)

const (
	uid = 12345
)

func TestUpPic(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(uid)
		ip  = ""
	)
	convey.Convey("Pic", t, func(ctx convey.C) {
		type result struct {
			Code int `json:"code"`
			Data Pic `json:"data"`
		}
		res := new(result)
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.client), "Get", func(_ *httpx.Client, _ context.Context, _, _ string, _ url.Values, _ interface{}) error {
			res.Code = 0
			res.Data = Pic{Has: 1}
			return nil
		})
		defer guard.Unpatch()
		has, err := d.Pic(c, mid, ip)
		ctx.Convey("Then err should be nil.has should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(has, convey.ShouldNotBeNil)
		})
	})
}

func TestUpBlink(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(uid)
		ip  = ""
	)
	convey.Convey("Blink", t, func(ctx convey.C) {
		type result struct {
			Code int   `json:"code"`
			Data Blink `json:"data"`
		}
		res := new(result)
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.client), "Get", func(_ *httpx.Client, _ context.Context, _, _ string, _ url.Values, _ interface{}) error {
			res.Code = 0
			res.Data = Blink{Has: 1}
			return nil
		})
		defer guard.Unpatch()
		has, err := d.Blink(c, mid, ip)
		ctx.Convey("Then err should be nil.has should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(has, convey.ShouldNotBeNil)
		})
	})
}

func TestUpIsAuthor(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(uid)
		ip  = ""
	)
	convey.Convey("IsAuthor", t, func(ctx convey.C) {
		isArt, err := d.IsAuthor(c, mid, ip)
		ctx.Convey("Then err should be nil.isArt should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(isArt, convey.ShouldNotBeNil)
		})
	})
}
