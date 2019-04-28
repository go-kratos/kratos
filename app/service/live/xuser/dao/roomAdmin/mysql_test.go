package roomAdmin

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestRoomAdminGetByUserMysql(t *testing.T) {
	convey.Convey("GetByUserMysql", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			uid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			admins, err := d.GetByUserMysql(c, uid)
			ctx.Convey("Then err should be nil.admins should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(admins, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestRoomAdminGetByRoomIdMysql(t *testing.T) {
	convey.Convey("GetByRoomIdMysql", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			roomId = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			admins, err := d.GetByRoomIdMysql(c, roomId)
			ctx.Convey("Then err should be nil.admins should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(admins, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestRoomAdminDelDbAdmin(t *testing.T) {
	convey.Convey("DelDbAdminMysql", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelDbAdminMysql(c, id)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestRoomAdminAddAdminMysql(t *testing.T) {
	convey.Convey("AddAdminMysql", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			uid    = int64(0)
			roomId = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.AddAdminMysql(c, uid, roomId)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestRoomAdminGetByRoomIdUidDb(t *testing.T) {
	convey.Convey("GetByRoomIdUidMysql", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			uid    = int64(0)
			roomId = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			resp, err := d.GetByRoomIdUidMysql(c, uid, roomId)
			ctx.Convey("Then err should be nil.resp should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(resp, convey.ShouldNotBeNil)
			})
		})
	})
}
