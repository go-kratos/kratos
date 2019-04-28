package growup

import (
	"context"
	"go-common/library/ecode"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestGrowupUpStatus(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		ip  = "127.0.0.1"
	)
	convey.Convey("UpStatus", t, func(ctx convey.C) {
		us, err := d.UpStatus(c, mid, ip)
		ctx.Convey("Then err should be nil.us should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldEqual, ecode.CreativeOrderAPIErr)
			ctx.So(us, convey.ShouldNotBeNil)
		})
	})
}

func TestGrowupUpInfo(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		ip  = "127.0.0.1"
	)
	convey.Convey("UpInfo", t, func(ctx convey.C) {
		ui, err := d.UpInfo(c, mid, ip)
		ctx.Convey("Then err should be nil.ui should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(ui, convey.ShouldBeNil)
		})
	})
}

func TestGrowupJoin(t *testing.T) {
	var (
		c      = context.TODO()
		mid    = int64(0)
		accTy  = int(0)
		signTy = int(0)
		ip     = ""
	)
	convey.Convey("Join", t, func(ctx convey.C) {
		err := d.Join(c, mid, accTy, signTy, ip)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldEqual, ecode.CreativeOrderAPIErr)
		})
	})
}

func TestGrowupQuit(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		ip  = "127.0.0.1"
	)
	convey.Convey("Quit1", t, func(ctx convey.C) {
		httpMock("POST", d.quitURL).Reply(200).JSON(`{"code":20008}`)
		err := d.Quit(c, mid, ip)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
	convey.Convey("Quit2", t, func(ctx convey.C) {
		httpMock("POST", d.quitURL).Reply(502)
		err := d.Quit(c, mid, ip)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestGrowupSummary(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		ip  = "127.0.0.1"
	)
	convey.Convey("Summary", t, func(ctx convey.C) {
		sm, err := d.Summary(c, mid, ip)
		ctx.Convey("Then err should be nil.sm should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldEqual, ecode.CreativeOrderAPIErr)
			ctx.So(sm, convey.ShouldBeNil)
		})
	})
}

func TestGrowupStat(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		ip  = "127.0.0.1"
	)
	convey.Convey("Stat", t, func(ctx convey.C) {
		st, err := d.Stat(c, mid, ip)
		ctx.Convey("Then err should be nil.st should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldEqual, ecode.CreativeOrderAPIErr)
			ctx.So(st, convey.ShouldBeNil)
		})
	})
}

func TestGrowupIncomeList(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		ty  = int(0)
		pn  = int(1)
		ps  = int(10)
		ip  = "127.0.0.1"
	)
	convey.Convey("IncomeList", t, func(ctx convey.C) {
		il, err := d.IncomeList(c, mid, ty, pn, ps, ip)
		ctx.Convey("Then err should be nil.il should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldEqual, ecode.CreativeOrderAPIErr)
			ctx.So(il, convey.ShouldBeNil)
		})
	})
}

func TestGrowupBreachList(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		pn  = int(1)
		ps  = int(10)
		ip  = "127.0.0.1"
	)
	convey.Convey("BreachList", t, func(ctx convey.C) {
		bl, err := d.BreachList(c, mid, pn, ps, ip)
		ctx.Convey("Then err should be nil.bl should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldEqual, ecode.CreativeOrderAPIErr)
			ctx.So(bl, convey.ShouldBeNil)
		})
	})
}
