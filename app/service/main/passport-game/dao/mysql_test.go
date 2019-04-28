package dao

import (
	"context"
	"go-common/app/service/main/passport-game/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaohit(t *testing.T) {
	var (
		mid = int64(1)
	)
	convey.Convey("hit", t, func(ctx convey.C) {
		p1 := hit(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoMemberInfo(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(1)
	)
	convey.Convey("MemberInfo", t, func(ctx convey.C) {
		res, err := d.MemberInfo(c, mid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestDaoApps(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("Apps", t, func(ctx convey.C) {
		res, err := d.Apps(c)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddToken(t *testing.T) {
	var (
		c  = context.TODO()
		no = &model.Perm{}
	)
	convey.Convey("AddToken", t, func(ctx convey.C) {
		affected, err := d.AddToken(c, no)
		ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(affected, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpdateToken(t *testing.T) {
	var (
		c  = context.TODO()
		no = &model.Perm{}
	)
	convey.Convey("UpdateToken", t, func(ctx convey.C) {
		affected, err := d.UpdateToken(c, no)
		ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affected, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoToken(t *testing.T) {
	var (
		c           = context.TODO()
		accessToken = "123456"
	)
	convey.Convey("Token", t, func(ctx convey.C) {
		res, err := d.Token(c, accessToken)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestDaoAsoAccount(t *testing.T) {
	var (
		c            = context.TODO()
		identify     = "123456"
		identifyHash = "654321"
	)
	convey.Convey("AsoAccount", t, func(ctx convey.C) {
		res, err := d.AsoAccount(c, identify, identifyHash)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestDaoAccountInfo(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(1)
	)
	convey.Convey("AccountInfo", t, func(ctx convey.C) {
		res, err := d.AccountInfo(c, mid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTokenFromOtherRegion(t *testing.T) {
	var (
		c           = context.TODO()
		accessToken = "123456"
	)
	convey.Convey("TokenFromOtherRegion", t, func(ctx convey.C) {
		res, err := d.TokenFromOtherRegion(c, accessToken)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}
