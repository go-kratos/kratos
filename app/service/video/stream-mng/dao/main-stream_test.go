package dao

import (
	"context"
	"go-common/app/service/video/stream-mng/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoCreateNewStream(t *testing.T) {
	convey.Convey("CreateNewStream", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			stream = &model.MainStream{
				RoomID:     123456,
				StreamName: "tetststest",
				Key:        "testtestetste",
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.CreateNewStream(c, stream)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetMainStreamFromDB(t *testing.T) {
	convey.Convey("GetMainStreamFromDB", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			roomID     = int64(11891462)
			streamName = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.GetMainStreamFromDB(c, roomID, streamName)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoChangeDefaultVendor(t *testing.T) {
	convey.Convey("ChangeDefaultVendor", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			roomID    = int64(11891462)
			newVendor = int64(2)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.ChangeDefaultVendor(c, roomID, newVendor)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoClearMainStreaming(t *testing.T) {
	convey.Convey("ClearMainStreaming", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			roomID     = int64(11891462)
			options    = int64(0)
			newoptions = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.ClearMainStreaming(c, roomID, newoptions, options)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoMainStreamNotify(t *testing.T) {
	convey.Convey("MainStreamNotify", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			roomID     = int64(11891462)
			vendor     = int64(2)
			isOpen     bool
			isOrigin   bool
			options    = int64(0)
			newoptions = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.MainStreamNotify(c, roomID, vendor, isOpen, isOrigin, options, newoptions)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
