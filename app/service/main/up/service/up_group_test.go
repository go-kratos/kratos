package service

import (
	"context"
	"testing"

	upgrpc "go-common/app/service/main/up/api/v1"
	"go-common/app/service/main/up/model"
	"go-common/library/net/http/blademaster"

	"github.com/smartystreets/goconvey/convey"
)

func TestServicecheckGroupData(t *testing.T) {
	convey.Convey("checkGroupData", t, func(ctx convey.C) {
		var (
			arg = &model.AddGroupArg{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := s.checkGroupData(arg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestServiceAddGroup(t *testing.T) {
	convey.Convey("AddGroup", t, func(ctx convey.C) {
		var (
			c   = &blademaster.Context{}
			arg = &model.AddGroupArg{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			result, err := s.AddGroup(c, arg)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceUpdateGroup(t *testing.T) {
	convey.Convey("UpdateGroup", t, func(ctx convey.C) {
		var (
			c   = &blademaster.Context{}
			arg = &model.EditGroupArg{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			result, err := s.UpdateGroup(c, arg)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceRemoveGroup(t *testing.T) {
	convey.Convey("RemoveGroup", t, func(ctx convey.C) {
		var (
			c   = &blademaster.Context{}
			arg = &model.RemoveGroupArg{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			result, err := s.RemoveGroup(c, arg)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceGetGroup(t *testing.T) {
	convey.Convey("GetGroup", t, func(ctx convey.C) {
		var (
			c   = &blademaster.Context{}
			arg = &model.GetGroupArg{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := s.GetGroup(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceGetGroupCache(t *testing.T) {
	convey.Convey("GetGroupCache", t, func(ctx convey.C) {
		var (
			groupID = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			group := s.getGroupCache(groupID)
			ctx.Convey("Then group should not be nil.", func(ctx convey.C) {
				ctx.So(group, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceUpGroups(t *testing.T) {
	convey.Convey("UpGroups", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			req = &upgrpc.NoArgReq{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := s.UpGroups(c, req)
			convCtx.Convey("No return values", func(convCtx convey.C) {
				convCtx.So(res, convey.ShouldNotBeNil)
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
