package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

func TestDaoLastPlayGame(t *testing.T) {
	convey.Convey("LastPlayGame", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(908085)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("GET", d.lastPlayGameURL).Reply(200).JSON(`{
    "code": 0,
    "games": [
        {
            "website": "https://www.bilibili.com/blackboard/activity-zlzyfk.html",
            "image": "http://i0.hdslb.com/bfs/game/abfe03ca09e2051e5edc2693499f5db4d72e0e79.png",
            "name": ""
        }
    ]
}`)
			data, err := d.LastPlayGame(c, mid)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When code not equal 0", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("GET", d.lastPlayGameURL).Reply(200).JSON(`{"code": -3}`)
			data, err := d.LastPlayGame(c, mid)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.Printf("%+v", data)
			})
		})
	})
}

func TestDaoAppPlayedGame(t *testing.T) {
	convey.Convey("AppPlayedGame", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(89090481)
			platform = _platformAndroid
			pn       = int(1)
			ps       = int(20)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("GET", d.appPlayedGameURL).Reply(200).JSON(`
{
    "code": 0,
    "message": "ok",
    "data": {
        "list": [
            {
                "game_base_id": 97,
                "game_name": "碧蓝航线",
                "game_icon": "https://i0.hdslb.com/bfs/game/3f975ae90395bd323f585a788eb5e852e7af625e.jpg",
                "grade": 8.8,
                "detail_url": "bilibili://game_center/detail?id=97&sourceFrom=666"
            }
        ],
        "total_count": 1
    }
}`)
			data, count, err := d.AppPlayedGame(c, mid, platform, pn, ps)
			ctx.Convey("Then err should be nil.data,count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}
