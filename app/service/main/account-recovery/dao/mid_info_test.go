package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

func TestDaoGetMidInfo(t *testing.T) {
	var (
		c     = context.Background()
		qType = "1"
		qKey  = "silg@yahoo.cn"
	)
	convey.Convey("GetMidInfo", t, func(ctx convey.C) {
		v, err := d.GetMidInfo(c, qType, qKey)
		ctx.Convey("Then err should be nil.v should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(v, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoGetUserInfo(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
	)
	convey.Convey("GetUserInfo", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("GET", d.c.AccRecover.GetUserInfoURL).Reply(200).JSON(`{"code":0,"data":{"mid":21,"email":"raiden131@yahoo.cn","telphone":"","join_time":1245902140}}`)
		v, err := d.GetUserInfo(c, mid)
		ctx.Convey("Then err should be nil.v should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(v, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpdatePwd(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
	)
	convey.Convey("UpdatePwd", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("POST", d.c.AccRecover.UpPwdURL).Reply(200).JSON(`{"code": 0, "data":{"pwd":"d4txsunbb1","userid":"minorin"}}`)
		user, err := d.UpdatePwd(c, mid, "账号找回服务")
		ctx.So(err, convey.ShouldBeNil)
		ctx.So(user, convey.ShouldNotBeNil)
	})
}

func TestDaoCheckSafe(t *testing.T) {
	var (
		c        = context.Background()
		mid      = int64(1)
		question = int8(0)
		answer   = "1"
	)
	convey.Convey("CheckSafe", t, func(ctx convey.C) {
		check, err := d.CheckSafe(c, mid, question, answer)
		ctx.Convey("Then err should be nil.check should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(check, convey.ShouldNotBeNil)
		})
	})
}

//func httpMock(method, url string) *gock.Request {
//	r := gock.New(url)
//	r.Method = strings.ToUpper(method)
//	return r
//}

//func TestDaoGetUserType(t *testing.T) {
//	var (
//		c   = context.Background()
//		mid = int64(2)
//	)
//	convey.Convey("When http request gets code != 0", t, func(ctx convey.C) {
//		defer gock.OffAll()
//		httpMock("GET", d.c.AccRecover.GameURL).Reply(0).JSON(`{"requestId":"0def8d70b7ef11e8a395fa163e01a2e9","ts":"1535440592","code":0,"items":[{"id":14,"name":"SDK测试2","lastLogin":"1500969010"}]}`)
//		games, err := d.GetUserType(c, mid)
//		ctx.So(err, convey.ShouldBeNil)
//		ctx.So(games, convey.ShouldNotBeNil)
//	})
//}

func TestDaoCheckReg(t *testing.T) {
	var (
		c       = context.Background()
		mid     = int64(1)
		regTime = int64(1532441644)
		regType = int8(0)
		regAddr = "中国_上海"
	)
	convey.Convey("CheckReg", t, func(ctx convey.C) {
		v, err := d.CheckReg(c, mid, regTime, regType, regAddr)
		ctx.Convey("Then err should be nil.v should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(v, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpdateBatchPwd(t *testing.T) {
	var (
		c    = context.Background()
		mids = "1,2"
	)
	convey.Convey("UpdateBatchPwd", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("POST", d.c.AccRecover.UpBatchPwdURL).Reply(200).JSON(`{"code":0,"data":{"6":{"pwd":"tgs52r1st9","userid":"腹黑君"},"7":{"pwd":"g20ahzrf7j","userid":"Tzwcard"}}}`)
		userMap, err := d.UpdateBatchPwd(c, mids, "账号找回服务")
		ctx.So(err, convey.ShouldBeNil)
		ctx.So(userMap, convey.ShouldNotBeNil)
	})
}

func TestDaoCheckCard(t *testing.T) {
	var (
		c        = context.Background()
		mid      = int64(1)
		cardType = int8(1)
		cardCode = "123"
	)
	convey.Convey("CheckCard", t, func(ctx convey.C) {
		ok, err := d.CheckCard(c, mid, cardType, cardCode)
		ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ok, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoCheckPwds(t *testing.T) {
	var (
		c    = context.Background()
		mid  = int64(1)
		pwds = "123"
	)
	convey.Convey("CheckPwds", t, func(ctx convey.C) {
		v, err := d.CheckPwds(c, mid, pwds)
		ctx.Convey("Then err should be nil.v should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(v, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoGetLoginIPs(t *testing.T) {
	var (
		c     = context.Background()
		mid   = int64(2)
		limit = int64(10)
	)
	convey.Convey("GetLoginIPs", t, func(ctx convey.C) {
		ipInfo, err := d.GetLoginIPs(c, mid, limit)
		ctx.Convey("Then err should be nil.ipInfo should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ipInfo, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoGetAddrByIP(t *testing.T) {
	var (
		c     = context.Background()
		mid   = int64(111001254)
		limit = int64(10)
	)
	convey.Convey("GetAddrByIP", t, func(ctx convey.C) {
		addrs, err := d.GetAddrByIP(c, mid, limit)
		ctx.Convey("Then err should be nil.addrs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(addrs, convey.ShouldNotBeNil)
		})
	})
}
