package dao

import (
	"context"
	"fmt"
	"go-common/app/service/main/spy/model"
	"net/url"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaohmacsha1(t *testing.T) {
	convey.Convey("hmacsha1", t, func(ctx convey.C) {
		var (
			key  = ""
			text = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			h := d.hmacsha1(key, text)
			ctx.Convey("Then h should not be nil.", func(ctx convey.C) {
				ctx.So(h, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaomakeURL(t *testing.T) {
	convey.Convey("makeURL", t, func(ctx convey.C) {
		var (
			method    = "1"
			action    = "1"
			region    = "1"
			secretID  = "1"
			secretKey = "1"
			charset   = "1"
			URL       = "1"
		)
		args := url.Values{}
		args.Set("accountType", fmt.Sprintf("%d", model.AccountType))
		args.Set("uid", fmt.Sprintf("%d", 1))
		args.Set("phoneNumber", "13262609601")
		args.Set("registerTime", fmt.Sprintf("%d", time.Now().Unix()))
		args.Set("registerIp", "127.0.0.1")
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			req := d.makeURL(method, action, region, secretID, secretKey, args, charset, URL)
			ctx.Convey("Then req should not be nil.", func(ctx convey.C) {
				ctx.So(req, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaomakeQueryString(t *testing.T) {
	convey.Convey("makeQueryString", t, func(ctx convey.C) {
		var (
			v url.Values
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			str := d.makeQueryString(v)
			ctx.Convey("Then str should not be nil.", func(ctx convey.C) {
				ctx.So(str, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRegisterProtection(t *testing.T) {
	convey.Convey("RegisterProtection", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			ip = "127.0.0.1"
		)
		args := url.Values{}
		args.Set("accountType", fmt.Sprintf("%d", model.AccountType))
		args.Set("uid", fmt.Sprintf("%d", 1))
		args.Set("phoneNumber", "13262609601")
		args.Set("registerTime", fmt.Sprintf("%d", time.Now().Unix()))
		args.Set("registerIp", ip)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			level, err := d.RegisterProtection(c, args, ip)
			ctx.Convey("Then err should be nil.level should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(level, convey.ShouldNotBeNil)
			})
		})
	})
}
