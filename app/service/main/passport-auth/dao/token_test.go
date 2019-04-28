package dao

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoToken(t *testing.T) {
	var (
		c                = context.TODO()
		token, _         = hex.DecodeString("baa3443180f346db244780ba6d0c6f6c")
		tokenNotExist, _ = hex.DecodeString("baa3443180f346db244780ba6d0c6f61")
		ct, _            = time.Parse("01/02/2006", "10/27/2018")
	)
	convey.Convey("Token", t, func(ctx convey.C) {
		res, err := d.Token(c, token, ct)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
		res2, err2 := d.Token(c, tokenNotExist, ct)
		ctx.Convey("Then err should be nil.res should be nil.", func(ctx convey.C) {
			ctx.So(err2, convey.ShouldBeNil)
			ctx.So(res2, convey.ShouldBeNil)
		})
	})
}

func TestDaoformatSuffix(t *testing.T) {
	var (
		no = time.Now()
	)
	convey.Convey("formatSuffix", t, func(ctx convey.C) {
		p1 := formatSuffix(no)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}
