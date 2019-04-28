package dao

import (
	"context"
	"fmt"
	"net"
	"testing"

	"go-common/app/admin/main/activity/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoimgAddKey(t *testing.T) {
	convey.Convey("imgAddKey", t, func(ctx convey.C) {
		var (
			uri = "http://i2.hdslb.com/bfs/face/5d2c92beb774a4bb30762538bb102d23670ae9c0.gif"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			url := imgAddKey(uri)
			ctx.Convey("Then url should not be nil.", func(ctx convey.C) {
				ctx.So(url, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetLikeContent(t *testing.T) {
	convey.Convey("GetLikeContent", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{1, 2}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			outRes, err := d.GetLikeContent(c, ids)
			ctx.Convey("Then err should be nil.outRes should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%+v", outRes)
			})
		})
	})
}

func TestDaoActSubject(t *testing.T) {
	convey.Convey("ActSubject", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sid = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rp, err := d.ActSubject(c, sid)
			ctx.Convey("Then err should be nil.rp should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%+v", rp)
			})
		})
	})
}

func TestDaoMusics(t *testing.T) {
	convey.Convey("Musics", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			aids = []int64{1, 2}
			ip   = "10.248.25.36"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			music, err := d.Musics(c, aids, ip)
			ctx.Convey("Then err should be nil.music should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%+v", music)
			})
		})
	})
}

func TestDaoBatchLike(t *testing.T) {
	convey.Convey("BatchLike", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			item = &model.Like{Sid: 10256, Mid: 55}
			wids = []int64{}
			ipv6 = net.ParseIP("192.168.3.32")
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.BatchLike(c, item, wids, ipv6)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
