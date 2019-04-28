package datadao

import (
	"context"
	"encoding/json"
	"github.com/bouk/monkey"
	"go-common/library/net/http/blademaster"
	"net/url"
	"reflect"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDatadaoHTTPDataHandle(t *testing.T) {
	var (
		c      = context.Background()
		params url.Values
		key    = "archives"
	)
	convey.Convey("HTTPDataHandle", t, func(ctx convey.C) {
		type result struct {
			Code    int             `json:"code"`
			Data    json.RawMessage `json:"data"`
			Message string          `json:"message"`
		}
		res := new(result)
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.bmClient), "Get", func(_ *blademaster.Client, _ context.Context, _, _ string, _ url.Values, _ interface{}) error {
			res.Code = 0
			res.Data = json.RawMessage("play")
			res.Message = "success"
			return nil
		})
		defer guard.Unpatch()
		_, err := d.HTTPDataHandle(c, params, key)
		ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDatadaogetURI(t *testing.T) {
	convey.Convey("getURI", t, func(ctx convey.C) {
		var (
			key = "archives"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			uri, err := d.getURI(key)
			ctx.Convey("Then err should be nil.uri should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(uri, convey.ShouldNotBeNil)
			})
		})
	})
}
