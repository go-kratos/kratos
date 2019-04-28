package block

import (
	"context"
	"testing"
	"time"

	"go-common/app/service/main/member/conf"
	model "go-common/app/service/main/member/model/block"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaouserKey(t *testing.T) {
	convey.Convey("userKey", t, func(ctx convey.C) {
		var (
			mid = int64(46333)
		)
		ctx.Convey("When everything right", func(ctx convey.C) {
			key := userKey(mid)
			ctx.Convey("Then key should equal u_mid.", func(ctx convey.C) {
				ctx.So(key, convey.ShouldEqual, "u_46333")
			})
		})
	})
}

func TestDaoSetUserCache(t *testing.T) {
	convey.Convey("SetUserCache", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			mid       = int64(46333)
			status    = model.BlockStatusForever
			startTime = time.Now().Unix()
			endTime   = time.Now().Add(time.Minute).Unix()
		)
		ctx.Convey("When SetUserCache", func(ctx convey.C) {
			err := d.SetUserCache(c, mid, status, startTime, endTime)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.Convey("When get UsersCache", func(ctx convey.C) {
					res, err := d.UsersCache(c, []int64{mid})
					ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
						ctx.So(err, convey.ShouldBeNil)
						ctx.So(res, convey.ShouldNotBeNil)
					})
				})
				ctx.Convey("When delete UsersCache", func(ctx convey.C) {
					err := d.DeleteUserCache(c, mid)
					ctx.Convey("Then err should be nil.", func(ctx convey.C) {
						ctx.So(err, convey.ShouldBeNil)
					})
				})
			})
		})
	})
}

func TestDaomcUserExpire(t *testing.T) {
	convey.Convey("mcUserExpire", t, func(ctx convey.C) {
		ctx.Convey("When everything right", func(ctx convey.C) {
			sec := d.mcUserExpire("ut_key")
			ctx.Convey("Then sec should >=mcUserExpireBase and <=mcUserExpireBase*conf.Conf.Memcache.Expire.UserMaxRate .", func(ctx convey.C) {
				ctx.So(d.UserTTL, convey.ShouldBeGreaterThan, 0)
				ctx.So(conf.Conf.BlockCacheTTL.UserMaxRate, convey.ShouldBeGreaterThan, 0)
				ctx.So(sec, convey.ShouldBeGreaterThanOrEqualTo, d.UserTTL)
				ctx.So(sec, convey.ShouldBeLessThanOrEqualTo, int32(conf.Conf.BlockCacheTTL.UserMaxRate*float64(d.UserTTL)))
			})
		})
	})
}

func TestDaoSetUserDetailCache(t *testing.T) {
	convey.Convey("SetUserDetailCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(46333)
		)
		ctx.Convey("When SetUserDetailCache", func(ctx convey.C) {
			err := d.SetUserDetailCache(c, mid, &model.MCUserDetail{BlockCount: 12})
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.Convey("When get UserDetailCache", func(ctx convey.C) {
					res, err := d.UserDetailsCache(c, []int64{mid})
					ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
						ctx.So(err, convey.ShouldBeNil)
						ctx.So(res, convey.ShouldNotBeNil)
					})
				})
				ctx.Convey("When delete UsersDetailCache", func(ctx convey.C) {
					err := d.DeleteUserDetailCache(c, mid)
					ctx.Convey("Then err should be nil.", func(ctx convey.C) {
						ctx.So(err, convey.ShouldBeNil)
					})
				})
			})
		})
	})
}
