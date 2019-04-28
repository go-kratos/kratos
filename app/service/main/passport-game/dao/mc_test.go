package dao

import (
	"context"
	"go-common/app/service/main/passport-game/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyInfoPB(t *testing.T) {
	var (
		mid = int64(12)
	)
	convey.Convey("keyInfoPB", t, func(ctx convey.C) {
		p1 := keyInfoPB(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaokeyTokenPB(t *testing.T) {
	var (
		accessToken = "123"
	)
	convey.Convey("keyTokenPB", t, func(ctx convey.C) {
		p1 := keyTokenPB(accessToken)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaokeyOriginMissMatchFlag(t *testing.T) {
	var (
		identify = "123"
	)
	convey.Convey("keyOriginMissMatchFlag", t, func(ctx convey.C) {
		p1 := keyOriginMissMatchFlag(identify)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaopingMC(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("pingMC", t, func(ctx convey.C) {
		err := d.pingMC(c)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSetInfoCache(t *testing.T) {
	var (
		c    = context.TODO()
		info = &model.Info{}
	)
	convey.Convey("SetInfoCache", t, func(ctx convey.C) {
		err := d.SetInfoCache(c, info)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoInfoCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
	)
	convey.Convey("InfoCache", t, func(ctx convey.C) {
		info, err := d.InfoCache(c, mid)
		ctx.Convey("Then err should be nil.info should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(info, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTokenCache(t *testing.T) {
	var (
		c           = context.TODO()
		accessToken = "123456"
		token       = &model.Perm{
			Mid:         -1,
			AccessToken: accessToken,
		}
	)
	convey.Convey("TestDaoTokenCache", t, func(ctx convey.C) {
		err := d.SetTokenCache(c, token)
		ctx.Convey("SetTokenCache", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})

		token, err := d.TokenCache(c, accessToken)
		ctx.Convey("TokenCache", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(token, convey.ShouldNotBeNil)
		})

		err = d.DelTokenCache(c, accessToken)
		ctx.Convey("DelTokenCache", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})

		token, err = d.TokenCache(c, accessToken)
		ctx.Convey("DeletedTokenCache", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(token, convey.ShouldBeNil)
		})
	})
}

func TestDaoSetOriginMissMatchFlagCache(t *testing.T) {
	var (
		c        = context.TODO()
		identify = "123456"
		flag     = []byte("")
	)
	convey.Convey("SetOriginMissMatchFlagCache", t, func(ctx convey.C) {
		err := d.SetOriginMissMatchFlagCache(c, identify, flag)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoOriginMissMatchFlagCache(t *testing.T) {
	var (
		c        = context.TODO()
		identify = "123456"
	)
	convey.Convey("OriginMissMatchFlagCache", t, func(ctx convey.C) {
		res, err := d.OriginMissMatchFlagCache(c, identify)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelOriginMissMatchFlagCache(t *testing.T) {
	var (
		c        = context.TODO()
		identify = "123456"
	)
	convey.Convey("DelOriginMissMatchFlagCache", t, func(ctx convey.C) {
		err := d.DelOriginMissMatchFlagCache(c, identify)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoOriginTokenCache(t *testing.T) {
	var (
		c           = context.TODO()
		accessToken = "123456"
		token       = &model.Token{
			Mid:         -1,
			AccessToken: accessToken,
		}
	)
	convey.Convey("TestDaoOriginTokenCache", t, func(ctx convey.C) {
		err := d.SetOriginTokenCache(c, token)
		ctx.Convey("SetOriginTokenCache", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})

		token, err := d.OriginTokenCache(c, accessToken)
		ctx.Convey("OriginTokenCache", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(token, convey.ShouldNotBeNil)
		})

		err = d.DelOriginTokenCache(c, accessToken)
		ctx.Convey("DelOriginTokenCache", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})

		token, err = d.OriginTokenCache(c, accessToken)
		ctx.Convey("DeletedOriginTokenCache", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(token, convey.ShouldBeNil)
		})
	})
}
