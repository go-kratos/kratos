package service

import (
	"context"
	"testing"

	adminmodel "go-common/app/admin/main/mcn/model"
	"go-common/app/interface/main/mcn/model"
	"go-common/app/interface/main/mcn/model/mcnmodel"

	"github.com/smartystreets/goconvey/convey"
)

func TestServicegetMcnWithState(t *testing.T) {
	convey.Convey("getMcnWithState", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mcnmid = int64(0)
			state  = model.MCNSignState(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			mcnSign, err := s.getMcnWithState(c, mcnmid, state)
			ctx.Convey("Then err should be nil.mcnSign should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(mcnSign, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceMcnGetState(t *testing.T) {
	convey.Convey("McnGetState", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &mcnmodel.GetStateReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := s.McnGetState(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceMcnExist(t *testing.T) {
	convey.Convey("McnExist", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &mcnmodel.GetStateReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := s.McnExist(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceMcnApply(t *testing.T) {
	convey.Convey("McnApply", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &mcnmodel.McnApplyReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := s.McnApply(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceMcnBindUpApply(t *testing.T) {
	convey.Convey("McnBindUpApply", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &mcnmodel.McnBindUpApplyReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := s.McnBindUpApply(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceMcnUpConfirm(t *testing.T) {
	convey.Convey("McnUpConfirm", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &mcnmodel.McnUpConfirmReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := s.McnUpConfirm(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceMcnUpGetBind(t *testing.T) {
	convey.Convey("McnUpGetBind", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &mcnmodel.McnUpGetBindReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := s.McnUpGetBind(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceMcnDataSummary(t *testing.T) {
	convey.Convey("McnDataSummary", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &mcnmodel.McnGetDataSummaryReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := s.McnDataSummary(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceMcnDataUpList(t *testing.T) {
	convey.Convey("McnDataUpList", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &mcnmodel.McnGetUpListReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := s.McnDataUpList(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceMcnGetOldInfo(t *testing.T) {
	convey.Convey("McnGetOldInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &mcnmodel.McnGetMcnOldInfoReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := s.McnGetOldInfo(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestCheckPermission(t *testing.T) {
	convey.Convey("checkPermission", t, func(ctx convey.C) {
		var (
			c                   = context.Background()
			mcnMid, upMid int64 = 15555180, 27515410
			permissions         = []adminmodel.AttrBasePermit{adminmodel.AttrBasePermitBit, adminmodel.AttrDataPermitBit}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res := s.checkPermission(c, mcnMid, upMid, permissions...)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(res, convey.ShouldBeTrue)
			})
		})
	})
}
