package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/broadcast/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyMidServer(t *testing.T) {
	var (
		mid = int64(1)
	)
	convey.Convey("keyMidServer", t, func(ctx convey.C) {
		p1 := keyMidServer(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaokeyKeyServer(t *testing.T) {
	var (
		key = "key"
	)
	convey.Convey("keyKeyServer", t, func(ctx convey.C) {
		p1 := keyKeyServer(key)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaokeyServerOnline(t *testing.T) {
	var (
		key = "key"
	)
	convey.Convey("keyServerOnline", t, func(ctx convey.C) {
		p1 := keyServerOnline(key)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaopingRedis(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("pingRedis", t, func(ctx convey.C) {
		err := d.pingRedis(c)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAddMapping(t *testing.T) {
	var (
		c      = context.Background()
		mid    = int64(1)
		key    = "key"
		server = "server"
	)
	convey.Convey("AddMapping", t, func(ctx convey.C) {
		err := d.AddMapping(c, mid, key, server)
		ctx.So(err, convey.ShouldBeNil)
		has, err := d.ExpireMapping(c, mid, key)
		ctx.So(err, convey.ShouldBeNil)
		ctx.So(has, convey.ShouldBeTrue)
		has, err = d.DelMapping(c, mid, key, server)
		ctx.So(err, convey.ShouldBeNil)
		ctx.So(has, convey.ShouldBeTrue)
		// false
		has, err = d.ExpireMapping(c, mid, key)
		ctx.So(err, convey.ShouldBeNil)
		ctx.So(has, convey.ShouldBeFalse)
		has, err = d.DelMapping(c, mid, key, server)
		ctx.So(err, convey.ShouldBeNil)
		ctx.So(has, convey.ShouldBeFalse)
	})
}

func TestDaoServersByKeys(t *testing.T) {
	var (
		c    = context.Background()
		keys = []string{"key"}
	)
	convey.Convey("ServersByKeys", t, func(ctx convey.C) {
		res, err := d.ServersByKeys(c, keys)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoKeysByMids(t *testing.T) {
	var (
		c    = context.Background()
		mids = []int64{1, 2, 3}
	)
	convey.Convey("KeysByMids", t, func(ctx convey.C) {
		ress, _, err := d.KeysByMids(c, mids)
		ctx.Convey("Then err should be nil.ress should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ress, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddServerOnline(t *testing.T) {
	var (
		c      = context.Background()
		server = "key"
		shard  = 1
		online = &model.Online{RoomCount: map[string]int32{
			"test1": 100,
			"test2": 200,
			"test3": 300,
		}, Updated: 1}
	)
	convey.Convey("AddServerOnline", t, func(ctx convey.C) {
		err := d.AddServerOnline(c, server, int32(shard), online)
		ctx.So(err, convey.ShouldBeNil)
		res, err := d.ServerOnline(c, server, shard)
		ctx.So(err, convey.ShouldBeNil)
		ctx.So(res, convey.ShouldResemble, online)
		err = d.DelServerOnline(c, server, shard)
		ctx.So(err, convey.ShouldBeNil)
	})
}

func TestDaoSetServers(t *testing.T) {
	var (
		c       = context.Background()
		servers = []*model.ServerInfo{}
	)
	convey.Convey("SetServers", t, func(ctx convey.C) {
		ctx.So(d.SetServers(c, servers), convey.ShouldBeNil)
		res, err := d.Servers(c)
		ctx.So(err, convey.ShouldBeNil)
		ctx.So(res, convey.ShouldResemble, servers)
		_, _, err = d.MigrateServers(c)
		ctx.So(err, convey.ShouldBeNil)
		_, err = d.MigrateRooms(c, 0)
		ctx.So(err, convey.ShouldBeNil)
	})
}
