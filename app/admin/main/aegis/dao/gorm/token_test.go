package gorm

import (
	"go-common/app/admin/main/aegis/model/net"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoTokens(t *testing.T) {
	var (
		ids = []int64{}
	)
	convey.Convey("Tokens", t, func(ctx convey.C) {
		no, err := d.Tokens(cntx, ids)
		ctx.Convey("Then err should be nil.no should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(no, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTokenList(t *testing.T) {
	convey.Convey("TokenList", t, func(ctx convey.C) {
		result, err := d.TokenList(cntx, []int64{1}, []int64{}, "1", true)
		ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(result, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTokenByID(t *testing.T) {
	convey.Convey("TokenByID", t, func(ctx convey.C) {
		d.TokenByID(cntx, 1)
		ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
		})
	})
}

func TestDaoTokenListWithPager(t *testing.T) {
	var (
		pm = &net.ListTokenParam{}
	)
	convey.Convey("TokenListWithPager", t, func(ctx convey.C) {
		result, err := d.TokenListWithPager(cntx, pm)
		ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(result, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTokenByUnique(t *testing.T) {
	convey.Convey("TokenByUnique", t, func(ctx convey.C) {
		_, err := d.TokenByUnique(cntx, 0, "", 0, "")
		ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoTokenBinds(t *testing.T) {
	convey.Convey("TokenBinds", t, func(ctx convey.C) {
		result, err := d.TokenBinds(cntx, []int64{})
		ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(result, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTokenBindByElement(t *testing.T) {
	convey.Convey("TokenBindByElement", t, func(ctx convey.C) {
		result, err := d.TokenBindByElement(cntx, []int64{}, []int8{}, true)
		ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(result, convey.ShouldNotBeNil)
		})
	})
}
