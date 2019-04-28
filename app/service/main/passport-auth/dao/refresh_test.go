package dao

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoRefresh(t *testing.T) {
	var (
		c             = context.TODO()
		rk, _         = hex.DecodeString("5f1813d287eb4238f2eadf225bda5d84")
		rkNotExist, _ = hex.DecodeString("5f1813d287eb4238f2eadf225bda5d81")
		ct, _         = time.Parse("01/02/2006", "08/27/2018")
	)
	convey.Convey("Refresh", t, func(ctx convey.C) {
		res, err := d.Refresh(c, rk, ct)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
		res2, err2 := d.Refresh(c, rkNotExist, ct)
		ctx.Convey("Then err should be nil.res should be nil.", func(ctx convey.C) {
			ctx.So(err2, convey.ShouldBeNil)
			ctx.So(res2, convey.ShouldBeNil)
		})
	})
}

func TestDaoformatRefreshSuffix(t *testing.T) {
	var (
		no = time.Now()
	)
	convey.Convey("formatRefreshSuffix", t, func(ctx convey.C) {
		p1 := formatRefreshSuffix(no)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoformatByDate(t *testing.T) {
	var (
		year  = int(2018)
		month = int(8)
	)
	convey.Convey("formatByDate", t, func(ctx convey.C) {
		p1 := formatByDate(year, month)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}
