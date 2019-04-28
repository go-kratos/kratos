package weeklyhonor

import (
	"context"
	model "go-common/app/interface/main/creative/model/weeklyhonor"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestWeeklyhonorstatKey(t *testing.T) {
	convey.Convey("statKey", t, func(ctx convey.C) {
		var (
			mid  = int64(0)
			date = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := statKey(mid, date)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestWeeklyhonorhonorKey(t *testing.T) {
	convey.Convey("honorKey", t, func(ctx convey.C) {
		var (
			mid  = int64(0)
			date = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := honorKey(mid, date)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestWeeklyhonorhonorClickKey(t *testing.T) {
	convey.Convey("honorClickKey", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := honorClickKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestWeeklyhonorStatMC(t *testing.T) {
	convey.Convey("StatMC", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(0)
			date = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			hs, err := d.StatMC(c, mid, date)
			ctx.Convey("Then err should be nil.hs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(hs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestWeeklyhonorHonorMC(t *testing.T) {
	convey.Convey("HonorMC", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(0)
			date = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.HonorMC(c, mid, date)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestWeeklyhonorSetStatMC(t *testing.T) {
	convey.Convey("SetStatMC", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(0)
			date = ""
			hs   = &model.HonorStat{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetStatMC(c, mid, date, hs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestWeeklyhonorSetHonorMC(t *testing.T) {
	convey.Convey("SetHonorMC", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(0)
			date = ""
			hs   = &model.HonorLog{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetHonorMC(c, mid, date, hs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestWeeklyhonorSetClickMC(t *testing.T) {
	convey.Convey("SetClickMC", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetClickMC(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestWeeklyhonorClickMC(t *testing.T) {
	convey.Convey("ClickMC", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.ClickMC(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
