package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/passport-sns/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDao_snsKey(t *testing.T) {
	convey.Convey("snsKey", t, func(ctx convey.C) {
		var (
			platform = model.PlatformQQStr
			mid      = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res := snsKey(platform, mid)
			ctx.Convey("Then res should not be nil.", func(ctx convey.C) {
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDao_oauth2Key(t *testing.T) {
	convey.Convey("oauth2Key", t, func(ctx convey.C) {
		var (
			platform = model.PlatformQQStr
			openID   = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res := oauth2Key(platform, openID)
			ctx.Convey("Then res should not be nil.", func(ctx convey.C) {
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDao_SetSnsCache(t *testing.T) {
	convey.Convey("SetSnsCache", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(0)
			platform = model.PlatformQQStr
			qq       = &model.SnsProto{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SetSnsCache(c, mid, platform, qq)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDao_SetOauth2Cache(t *testing.T) {
	convey.Convey("SetOauth2Cache", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			openID   = ""
			platform = model.PlatformQQStr
			qq       = &model.Oauth2Proto{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SetOauth2Cache(c, openID, platform, qq)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDao_SnsCache(t *testing.T) {
	convey.Convey("SnsCache", t, func(ctx convey.C) {
		var (
			c           = context.Background()
			platform    = model.PlatformQQStr
			mid         = int64(0)
			midNotExist = int64(-1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			info, err := d.SnsCache(c, mid, platform)
			ctx.Convey("Then err should be nil.info should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(info, convey.ShouldNotBeNil)
			})
			info2, err2 := d.SnsCache(c, midNotExist, platform)
			ctx.Convey("Then err should be nil.info should be nil.", func(ctx convey.C) {
				ctx.So(err2, convey.ShouldBeNil)
				ctx.So(info2, convey.ShouldBeNil)
			})
		})
	})
}

func TestDao_Oauth2Cache(t *testing.T) {
	convey.Convey("Oauth2Cache", t, func(ctx convey.C) {
		var (
			c              = context.Background()
			openID         = ""
			openIDNotExist = "-1"
			platform       = model.PlatformQQStr
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			info, err := d.Oauth2Cache(c, openID, platform)
			ctx.Convey("Then err should be nil.info should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(info, convey.ShouldNotBeNil)
			})
			info2, err2 := d.Oauth2Cache(c, openIDNotExist, platform)
			ctx.Convey("Then err should be nil.info should be nil.", func(ctx convey.C) {
				ctx.So(err2, convey.ShouldBeNil)
				ctx.So(info2, convey.ShouldBeNil)
			})
		})
	})
}

func TestDao_DelSnsCache(t *testing.T) {
	convey.Convey("DelSnsCache", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(0)
			platform = model.PlatformQQStr
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelSnsCache(c, mid, platform)
			ctx.Convey("Then err should be nil.info should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDao_DelOauth2Cache(t *testing.T) {
	convey.Convey("DelOauth2Cache", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			openID   = ""
			platform = model.PlatformQQStr
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelOauth2Cache(c, openID, platform)
			ctx.Convey("Then err should be nil.info should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
