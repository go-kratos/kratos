package dao

import (
	"context"
	"go-common/app/service/video/stream-mng/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetBackupStreamByRoomID(t *testing.T) {
	convey.Convey("GetBackupStreamByRoomID", t, func(ctx convey.C) {
		var (
			ctx2 = context.Background()
			rid  = int64(11891462)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			infos, err := d.GetBackupStreamByRoomID(ctx2, rid)
			ctx.Convey("Then err should be nil.infos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(infos, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCreateBackupStream(t *testing.T) {
	convey.Convey("CreateBackupStream", t, func(ctx convey.C) {
		var (
			ctx2 = context.Background()
			bs   = &model.BackupStream{
				RoomID: 66666,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.CreateBackupStream(ctx2, bs)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetBackupStreamByStreamName(t *testing.T) {
	convey.Convey("GetBackupStreamByStreamName", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			sn = "live_1511284_bs_7317941"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.GetBackupStreamByStreamName(c, sn)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSetBackupStreamStreamingStatus(t *testing.T) {
	convey.Convey("SetBackupStreamStreamingStatus", t, func(ctx convey.C) {
		var (
			c = context.Background()
			p = &model.StreamingNotifyParam{
				SRC:        "bvc",
				Type:       "0",
				StreamName: "live_1511284_bs_7317941",
			}
			bs = &model.BackupStream{
				OriginUpstream: 1,
				StreamName:     "live_1511284_bs_7317941",
			}
			open bool
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.SetBackupStreamStreamingStatus(c, p, bs, open)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
