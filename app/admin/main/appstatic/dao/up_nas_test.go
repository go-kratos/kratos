package dao

import (
	"context"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSign(t *testing.T) {
	var (
		params = url.Values{}
		nas    = d.c.Nas
	)
	params.Set("appkey", nas.Key)
	params.Set("appsecret", nas.Secret)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	convey.Convey("Sign", t, func(ctx convey.C) {
		query, err := Sign(params)
		ctx.Convey("Then err should be nil.query should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(query, convey.ShouldNotBeNil)
		})
	})
}

func TestDaogetSign(t *testing.T) {
	var nas = d.c.Nas
	convey.Convey("getSign", t, func(ctx convey.C) {
		uri, err := getSign(nas)
		ctx.Convey("Then err should be nil.uri should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(uri, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUploadNas(t *testing.T) {
	var (
		c        = context.Background()
		fileName = "test.txt"
		data     = []byte("test123")
		nas      = d.c.Nas
	)
	convey.Convey("UploadNas", t, func(ctx convey.C) {
		location, err := d.UploadNas(c, fileName, data, nas)
		ctx.Convey("Then err should be nil.location should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(location, convey.ShouldNotBeNil)
		})
	})
}
