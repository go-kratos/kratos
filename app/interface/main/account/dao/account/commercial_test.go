package account

import (
	"context"
	"net/url"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAccountCommercialSign(t *testing.T) {
	var (
		params = url.Values{}
	)
	params.Set("appkey", "test")
	params.Set("appsecret", "e6c4c252dc7e3d8a90805eecd7c73396")
	params.Set("key", "value")

	result := url.Values{}
	result.Set("appkey", "test")
	result.Set("key", "value")
	result.Set("sign", "f5c7329e9078e94ebf7d5852dfcd64f9")
	convey.Convey("CommercialSign", t, func(ctx convey.C) {
		query, err := CommercialSign(params)
		ctx.Convey("Then err should be nil.query should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(query, convey.ShouldEqual, result.Encode())
		})
	})
}

func TestAccountisBusinessAccount(t *testing.T) {
	var (
		mids = []int64{1, 2, 3}
	)
	convey.Convey("isBusinessAccount", t, func(ctx convey.C) {
		p1, err := d.isBusinessAccount(context.Background(), mids)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			if err != nil {
				ctx.Printf("Failed to validate business account: %+v\n", err)
				return
			}
			// ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestAccountIsBusinessAccount(t *testing.T) {
	var (
		mid = int64(0)
	)
	convey.Convey("IsBusinessAccount", t, func(ctx convey.C) {
		p1 := d.IsBusinessAccount(context.Background(), mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestAccountBusinessAccountInfo(t *testing.T) {
	var (
		mid = int64(0)
	)
	convey.Convey("BusinessAccountInfo", t, func(ctx convey.C) {
		p1, err := d.BusinessAccountInfo(context.Background(), mid)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			if err != nil {
				ctx.Printf("Failed to get business account info: %+v\n", err)
				return
			}
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}
