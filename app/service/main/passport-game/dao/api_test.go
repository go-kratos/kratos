package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoMyInfo(t *testing.T) {
	var (
		c         = context.TODO()
		accessKey = "123456"
	)
	convey.Convey("MyInfo", t, func(ctx convey.C) {
		accountInfo, err := d.MyInfo(c, accessKey)
		ctx.Convey("Then err should be nil.accountInfo should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(accountInfo, convey.ShouldBeNil)
		})
	})
}

func TestDaoOauth(t *testing.T) {
	var (
		c         = context.TODO()
		uri       = "https://wwww.baidu.com"
		accessKey = "123456"
		from      = "baidu"
	)
	convey.Convey("Oauth", t, func(ctx convey.C) {
		token, err := d.Oauth(c, uri, accessKey, from)
		ctx.Convey("Then err should be nil.token should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(token, convey.ShouldBeNil)
		})
	})
}

func TestDaoLogin(t *testing.T) {
	var (
		c      = context.TODO()
		query  = "123"
		cookie = "123"
	)
	convey.Convey("Login", t, func(ctx convey.C) {
		loginToken, err := d.Login(c, query, cookie)
		ctx.Convey("Then err should be nil.loginToken should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(loginToken, convey.ShouldBeNil)
		})
	})
}

func TestDaoLoginOrigin(t *testing.T) {
	var (
		c      = context.TODO()
		userid = "1"
		rsaPwd = "123456"
	)
	convey.Convey("LoginOrigin", t, func(ctx convey.C) {
		loginToken, err := d.LoginOrigin(c, userid, rsaPwd)
		ctx.Convey("Then err should be nil.loginToken should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(loginToken, convey.ShouldBeNil)
		})
	})
}

func TestDaoRSAKeyOrigin(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("RSAKeyOrigin", t, func(ctx convey.C) {
		key, err := d.RSAKeyOrigin(c)
		ctx.Convey("Then err should be nil.key should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(key, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoRenewToken(t *testing.T) {
	var (
		c    = context.TODO()
		uri  = "https://wwww.baidu.com"
		ak   = "234"
		from = "234"
	)
	convey.Convey("RenewToken", t, func(ctx convey.C) {
		renewToken, err := d.RenewToken(c, uri, ak, from)
		ctx.Convey("Then err should be nil.renewToken should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(renewToken, convey.ShouldBeNil)
		})
	})
}
