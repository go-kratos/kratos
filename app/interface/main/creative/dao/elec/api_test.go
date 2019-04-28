package elec

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestElecUserInfo(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		ip  = "127.0.0.1"
	)
	convey.Convey("UserInfo", t, func(ctx convey.C) {
		st, err := d.UserInfo(c, mid, ip)
		ctx.Convey("Then err should be nil.st should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(st, convey.ShouldNotBeNil)
		})
	})
}

func TestElecUserUpdate(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		st  = int8(1)
		ip  = "127.0.0.1"
	)
	convey.Convey("UserUpdate", t, func(ctx convey.C) {
		u, err := d.UserUpdate(c, mid, st, ip)
		ctx.Convey("Then err should be nil.u should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(u, convey.ShouldNotBeNil)
		})
	})
}

func TestElecArcUpdate(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		aid = int64(10106491)
		st  = int8(1)
		ip  = "127.0.0.1"
	)
	convey.Convey("ArcUpdate", t, func(ctx convey.C) {
		err := d.ArcUpdate(c, mid, aid, st, ip)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestElecNotify(t *testing.T) {
	var (
		c  = context.TODO()
		ip = ""
	)
	convey.Convey("Notify", t, func(ctx convey.C) {
		nt, err := d.Notify(c, ip)
		ctx.Convey("Then err should be nil.nt should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(nt, convey.ShouldNotBeNil)
		})
	})
}

func TestElecStatus(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		ip  = "127.0.0.1"
	)
	convey.Convey("Status", t, func(ctx convey.C) {
		st, err := d.Status(c, mid, ip)
		ctx.Convey("Then err should be nil.st should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(st, convey.ShouldNotBeNil)
		})
	})
}

func TestElecUpStatus(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(2089809)
		spday = int(1)
		ip    = "127.0.0.1"
	)
	convey.Convey("UpStatus", t, func(ctx convey.C) {
		err := d.UpStatus(c, mid, spday, ip)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestElecRecentRank(t *testing.T) {
	var (
		c    = context.TODO()
		mid  = int64(2089809)
		size = int64(10)
		ip   = "127.0.0.1"
	)
	convey.Convey("RecentRank", t, func(ctx convey.C) {
		rec, err := d.RecentRank(c, mid, size, ip)
		ctx.Convey("Then err should be nil.rec should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(rec, convey.ShouldNotBeNil)
		})
	})
}

func TestElecCurrentRank(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		ip  = "127.0.0.1"
	)
	convey.Convey("CurrentRank", t, func(ctx convey.C) {
		cur, err := d.CurrentRank(c, mid, ip)
		ctx.Convey("Then err should be nil.cur should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(cur, convey.ShouldNotBeNil)
		})
	})
}

func TestElecTotalRank(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		ip  = "127.0.0.1"
	)
	convey.Convey("TotalRank", t, func(ctx convey.C) {
		tol, err := d.TotalRank(c, mid, ip)
		ctx.Convey("Then err should be nil.tol should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(tol, convey.ShouldNotBeNil)
		})
	})
}

func TestElecDailyBill(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(2089809)
		pn    = int(1)
		ps    = int(10)
		ip    = "127.0.0.1"
		begin = "2018-01-04"
		end   = "2019-01-04"
	)
	convey.Convey("DailyBill", t, func(ctx convey.C) {
		bl, err := d.DailyBill(c, mid, pn, ps, begin, end, ip)
		ctx.Convey("Then err should be nil.bl should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(bl, convey.ShouldNotBeNil)
		})
	})
}

func TestElecBalance(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		ip  = "127.0.0.1"
	)
	convey.Convey("Balance", t, func(ctx convey.C) {
		bal, err := d.Balance(c, mid, ip)
		ctx.Convey("Then err should be nil.bal should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(bal, convey.ShouldNotBeNil)
		})
	})
}

func TestElecRecentElec(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		pn  = int(1)
		ps  = int(10)
		ip  = "127.0.0.1"
	)
	convey.Convey("RecentElec", t, func(ctx convey.C) {
		rec, err := d.RecentElec(c, mid, pn, ps, ip)
		ctx.Convey("Then err should be nil.rec should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(rec, convey.ShouldNotBeNil)
		})
	})
}

func TestElecRemarkList(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(2089809)
		pn    = int(1)
		ps    = int(10)
		begin = ""
		end   = ""
		ip    = "127.0.0.1"
	)
	convey.Convey("RemarkList", t, func(ctx convey.C) {
		rec, err := d.RemarkList(c, mid, pn, ps, begin, end, ip)
		ctx.Convey("Then err should be nil.rec should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(rec, convey.ShouldNotBeNil)
		})
	})
}

func TestElecRemarkDetail(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		id  = int64(0)
		ip  = "127.0.0.1"
	)
	convey.Convey("RemarkDetail", t, func(ctx convey.C) {
		re, err := d.RemarkDetail(c, mid, id, ip)
		ctx.Convey("Then err should be nil.re should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(re, convey.ShouldNotBeNil)
		})
	})
}

func TestElecRemark(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		id  = int64(0)
		msg = ""
		ak  = ""
		ck  = ""
		ip  = "127.0.0.1"
	)
	convey.Convey("Remark", t, func(ctx convey.C) {
		status, err := d.Remark(c, mid, id, msg, ak, ck, ip)
		ctx.Convey("Then err should be nil.status should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(status, convey.ShouldNotBeNil)
		})
	})
}
