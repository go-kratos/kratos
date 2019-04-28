package roomAdmin

import (
	"context"
	"go-common/app/service/live/xuser/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestRoomAdminCacheRoomAdminRoom(t *testing.T) {
	convey.Convey("CacheRoomAdminRoom", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.CacheRoomAdminRoom(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestRoomAdminCacheRoomAdminUser(t *testing.T) {
	convey.Convey("CacheRoomAdminUser", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.CacheRoomAdminUser(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestRoomAdminAddCacheKeyAnchorRoom(t *testing.T) {
	convey.Convey("AddCacheKeyAnchorRoom", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(0)
			val = []*model.RoomAdmin{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.AddCacheKeyAnchorRoom(c, id, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestRoomAdminAddCacheRoomAdminUser(t *testing.T) {
	convey.Convey("AddCacheRoomAdminUser", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(0)
			val = []*model.RoomAdmin{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.AddCacheRoomAdminUser(c, id, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestRoomAdminDelCacheKeyAnchorRoom(t *testing.T) {
	convey.Convey("DelCacheKeyAnchorRoom", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelCacheKeyAnchorRoom(c, id)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestRoomAdminDelCacheRoomAdminUser(t *testing.T) {
	convey.Convey("DelCacheRoomAdminUser", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelCacheRoomAdminUser(c, id)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
