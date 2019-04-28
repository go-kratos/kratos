package dao

import (
	"context"
	"go-common/app/service/main/tv/internal/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoCacheUserInfoByMid(t *testing.T) {
	convey.Convey("CacheUserInfoByMid", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(27515308)
		)
		d.AddCacheUserInfoByMid(c, id, &model.UserInfo{})
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.CacheUserInfoByMid(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddCacheUserInfoByMid(t *testing.T) {
	convey.Convey("AddCacheUserInfoByMid", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(27515308)
			val = &model.UserInfo{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.AddCacheUserInfoByMid(c, id, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelCacheUserInfoByMid(t *testing.T) {
	convey.Convey("DelCacheUserInfoByMid", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(27515308)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelCacheUserInfoByMid(c, id)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCachePayParamByToken(t *testing.T) {
	convey.Convey("CachePayParamByToken", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = "TOKEN:34567345678"
		)
		d.AddCachePayParam(context.TODO(), id, &model.PayParam{})
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.CachePayParamByToken(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCachePayParamsByTokens(t *testing.T) {
	convey.Convey("CachePayParamsByTokens", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			ids = []string{"TOKEN:34567345678"}
		)
		d.AddCachePayParam(context.TODO(), ids[0], &model.PayParam{})
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.CachePayParamsByTokens(c, ids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddCachePayParam(t *testing.T) {
	convey.Convey("AddCachePayParam", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = "TOKEN:34567345678"
			val = &model.PayParam{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.AddCachePayParam(c, id, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpdateCachePayParam(t *testing.T) {
	convey.Convey("UpdateCachePayParam", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = "TOKEN:34567345678"
			val = &model.PayParam{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.UpdateCachePayParam(c, id, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
