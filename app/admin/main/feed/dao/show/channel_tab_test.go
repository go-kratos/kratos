package show

import (
	"testing"

	"go-common/app/admin/main/feed/model/show"

	"github.com/smartystreets/goconvey/convey"
)

func TestShowChannelTabAdd(t *testing.T) {
	convey.Convey("ChannelTabAdd", t, func(ctx convey.C) {
		var (
			param = &show.ChannelTabAP{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.ChannelTabAdd(param)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestShowChannelTabUpdate(t *testing.T) {
	convey.Convey("ChannelTabUpdate", t, func(ctx convey.C) {
		var (
			param = &show.ChannelTabUP{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.ChannelTabUpdate(param)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestShowChannelTabDelete(t *testing.T) {
	convey.Convey("ChannelTabDelete", t, func(ctx convey.C) {
		var (
			id = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.ChannelTabDelete(id)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestShowChannelTabValid(t *testing.T) {
	convey.Convey("ChannelTabValid", t, func(ctx convey.C) {
		var (
			id       = int64(0)
			tagID    = int64(0)
			sTime    = int64(0)
			eTime    = int64(0)
			priority = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			count, err := d.ChannelTabValid(id, tagID, sTime, eTime, priority)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}
