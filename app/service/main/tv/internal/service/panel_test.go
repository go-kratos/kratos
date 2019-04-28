package service

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceGetPpcsMap(t *testing.T) {
	convey.Convey("GetPpcsMap", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := s.GetPpcsMap()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServicePanelPriceConfigsBySuitType(t *testing.T) {
	convey.Convey("PanelPriceConfigsBySuitType", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			st = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ppcs, err := s.PanelPriceConfigsBySuitType(c, st)
			ctx.Convey("Then err should be nil.ppcs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ppcs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServicePanelPriceConfigByProductId(t *testing.T) {
	convey.Convey("PanelPriceConfigByProductId", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			productId = "zc20181206"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ppc, err := s.PanelPriceConfigByProductId(c, productId)
			ctx.Convey("Then err should be nil.ppc should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ppc, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServicePanelPriceConfigByPid(t *testing.T) {
	convey.Convey("PanelPriceConfigByPid", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			pid = int32(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ppc, err := s.PanelPriceConfigByPid(c, pid)
			ctx.Convey("Then err should be nil.ppc should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ppc, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceGuestPanelPriceConfigs(t *testing.T) {
	convey.Convey("GuestPanelPriceConfigs", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ppcs, err := s.GuestPanelPriceConfigs(c)
			ctx.Convey("Then err should be nil.ppcs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ppcs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceMVipPanelPriceConfigs(t *testing.T) {
	convey.Convey("MVipPanelPriceConfigs", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515308)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ppcs, err := s.MVipPanelPriceConfigs(c, mid)
			ctx.Convey("Then err should be nil.ppcs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ppcs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceGuestPanelInfo(t *testing.T) {
	convey.Convey("GuestPanelInfo", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			pi, err := s.GuestPanelInfo(c)
			ctx.Convey("Then err should be nil.pi should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(pi, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServicePanelInfo(t *testing.T) {
	convey.Convey("PanelInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515308)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			pi, err := s.PanelInfo(c, mid)
			ctx.Convey("Then err should be nil.pi should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(pi, convey.ShouldNotBeNil)
			})
		})
	})
}
