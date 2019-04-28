package vip

import (
	"context"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/live/xuser/model"
	"go-common/library/log"
	"testing"
	"time"
)

func TestDao_GetVipFromCache(t *testing.T) {
	initd()
	Convey("test get vip cache", t, testWithTestUser(func(u *TestUser) {
		var (
			ctx  = context.Background()
			info *model.VipInfo
			err  error
			conn = d.redis.Get(ctx)
			key  = getUserCacheKey(u.Uid)
		)
		log.Info("TestDao_GetVipFromCache uid(%d), key(%s)", u.Uid, key)

		// delete key at begin
		conn.Do("DEL", key)

		// should get nil info and err
		Convey("should get nil info and err", func() {
			info, err = d.GetVipFromCache(ctx, u.Uid)
			So(info, ShouldBeNil)
			So(err, ShouldBeNil)
		})

		// set empty data
		Convey("set empty data", func() {
			conn.Do("HSET", key, _vipFieldName, "")
			info, err = d.GetVipFromCache(ctx, u.Uid)
			So(info, ShouldBeNil)
			So(err, ShouldBeNil)
		})

		// set not a json
		Convey("set not a json", func() {
			conn.Do("HSET", key, _vipFieldName, "test")
			info, err = d.GetVipFromCache(ctx, u.Uid)
			So(info, ShouldBeNil)
			So(err, ShouldBeNil)
		})

		// test vip/svip string format
		Convey("test vip/svip string format", func() {
			vipTime := time.Now().Add(time.Hour * 12).Format(model.TimeNano)
			svipTime := time.Now().AddDate(0, -1, 0).Format(model.TimeNano)
			conn.Do("HSET", key, _vipFieldName, fmt.Sprintf(`{"vip":"1","vip_time":"%s","svip":0,"svip_time":"%s"}`, vipTime, svipTime))
			info, err = d.GetVipFromCache(ctx, u.Uid)
			So(info.Vip, ShouldEqual, 1)
			So(info.Svip, ShouldEqual, 0)
			So(err, ShouldBeNil)
		})

		// set valid data
		Convey("set valid data", func() {
			vipTime := time.Now().Add(time.Hour * 12).Format(model.TimeNano)
			conn.Do("HSET", key, _vipFieldName, fmt.Sprintf(`{"vip":1,"vip_time":"%s","svip":0,"svip_time":"0000-00-00 00:00:00"}`, vipTime))
			info, err = d.GetVipFromCache(ctx, u.Uid)
			So(info, ShouldNotBeNil)
			So(info.Vip, ShouldEqual, 1)
			So(info.VipTime, ShouldEqual, vipTime)
			So(info.Svip, ShouldEqual, 0)
			So(info.SvipTime, ShouldEqual, model.TimeEmpty)
			So(err, ShouldBeNil)
		})

		// expired vip time
		Convey("set valid but expired vip time", func() {
			vip := 1
			vipTime := time.Now().AddDate(0, 0, -1).Format(model.TimeNano)
			svip := 1
			svipTime := time.Now().AddDate(0, -1, -1).Format(model.TimeNano)
			conn.Do("HSET", key, _vipFieldName, fmt.Sprintf(`{"vip":%d,"vip_time":"%s","svip":%d,"svip_time":"%s"}`, vip, vipTime, svip, svipTime))
			info, err = d.GetVipFromCache(ctx, u.Uid)
			log.Info("expired vip time, info(%+v)", info)
			So(info, ShouldNotBeNil)
			So(info.Vip, ShouldEqual, 0)
			So(info.VipTime, ShouldEqual, vipTime)
			So(info.Svip, ShouldEqual, 0)
			So(info.SvipTime, ShouldEqual, svipTime)
			So(err, ShouldBeNil)
		})
	}))
}

func TestDao_SetVipCache(t *testing.T) {
	initd()
	Convey("test set vip cache", t, testWithTestUser(func(u *TestUser) {
		var (
			ctx  = context.Background()
			info *model.VipInfo
			err  error
			conn = d.redis.Get(ctx)
			key  = getUserCacheKey(u.Uid)
		)
		log.Info("TestDao_GetVipFromCache uid(%d), key(%s)", u.Uid, key)
		// delete key at begin
		conn.Do("DEL", key)

		// nil info
		Convey("nil info", func() {
			err = d.SetVipCache(ctx, u.Uid, nil)
			So(err, ShouldBeNil)
			info, err = d.GetVipFromCache(ctx, u.Uid)
			log.Info("TestDao_SetVipCache get info1(%v), err(%v)", info, err)
			So(err, ShouldBeNil)
			So(info.Vip, ShouldEqual, 0)
			So(info.VipTime, ShouldEqual, model.TimeEmpty)
			So(info.Svip, ShouldEqual, 0)
			So(info.SvipTime, ShouldEqual, model.TimeEmpty)
		})

		// set valid info
		Convey("set valid info", func() {
			info = &model.VipInfo{
				Vip:      1,
				VipTime:  time.Now().Add(time.Hour * 12).Format(model.TimeNano),
				Svip:     1,
				SvipTime: time.Now().Add(time.Hour * 6).Format(model.TimeNano),
			}
			err = d.SetVipCache(ctx, u.Uid, info)
			So(err, ShouldBeNil)
			info2, err := d.GetVipFromCache(ctx, u.Uid)
			log.Info("TestDao_SetVipCache get info2(%v), err(%v)", info, err)
			So(err, ShouldBeNil)
			So(info2.Vip, ShouldEqual, info.Vip)
			So(info2.VipTime, ShouldEqual, info.VipTime)
			So(info2.Svip, ShouldEqual, info.Svip)
			So(info2.SvipTime, ShouldEqual, info.SvipTime)
		})
	}))
}

func TestDao_ClearCache(t *testing.T) {
	initd()
	Convey("test clear cache", t, testWithTestUser(func(u *TestUser) {
		var (
			ctx  = context.Background()
			err  error
			conn = d.redis.Get(ctx)
			key  = getUserCacheKey(u.Uid)
		)
		log.Info("TestDao_ClearCache uid(%d), key(%s)", u.Uid, key)
		// delete key at begin
		conn.Do("DEL", key)

		// del already deleted key
		Convey("del already deleted key", func() {
			err = d.ClearCache(ctx, u.Uid)
			So(err, ShouldBeNil)
		})

		// set valid info
		Convey("set valid info", func() {
			err = d.SetVipCache(ctx, u.Uid, nil)
			So(err, ShouldBeNil)
			err = d.ClearCache(ctx, u.Uid)
			So(err, ShouldBeNil)
		})
	}))
}
