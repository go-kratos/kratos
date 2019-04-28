package game

import (
	"context"
	"go-common/app/interface/main/creative/model/game"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

func TestGameList(t *testing.T) {
	var (
		c          = context.TODO()
		keywordStr = "kw"
		ip         = "127.0.0.1"
		err        error
		data       []*game.ListItem
	)
	convey.Convey("List", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("Do", d.gameListURL).Reply(-502)
		data, err = d.List(c, keywordStr, ip)
		ctx.Convey("List", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(data, convey.ShouldBeNil)
		})
	})
}

func TestGameInfo(t *testing.T) {
	var (
		c           = context.TODO()
		gameBaseID  = int64(0)
		platForType = int(0)
		ip          = "127.0.0.1"
		info        *game.Info
		err         error
	)
	convey.Convey("Info", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("Do", d.gameListURL).Reply(-502)
		info, err = d.Info(c, gameBaseID, platForType, ip)
		ctx.Convey("Info", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(info, convey.ShouldBeNil)
		})
	})
}
