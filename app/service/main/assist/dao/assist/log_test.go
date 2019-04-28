package assist

import (
	"context"
	"go-common/library/ecode"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestAssistAddLog(t *testing.T) {
	var (
		c         = context.Background()
		mid       = int64(0)
		assistMid = int64(0)
		tp        = int64(0)
		act       = int64(0)
		subID     = int64(0)
		objIDStr  = ""
		detail    = ""
	)
	convey.Convey("AddLog", t, func(ctx convey.C) {
		id, err := d.AddLog(c, mid, assistMid, tp, act, subID, objIDStr, detail)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(id, convey.ShouldNotBeNil)
		})
	})
}

func TestAssistCancelLog(t *testing.T) {
	var (
		c         = context.Background()
		logID     = int64(0)
		mid       = int64(0)
		assistMid = int64(0)
	)
	convey.Convey("CancelLog", t, func(ctx convey.C) {
		rows, err := d.CancelLog(c, logID, mid, assistMid)
		ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(rows, convey.ShouldNotBeNil)
		})
	})
}

func TestAssistLogInfo(t *testing.T) {
	var (
		c         = context.Background()
		id        = int64(1)
		mid       = int64(27515256)
		assistMid = int64(27515243)
	)
	convey.Convey("LogInfo", t, func(ctx convey.C) {
		a, err := d.LogInfo(c, id, mid, assistMid)
		ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(err, convey.ShouldEqual, ecode.AssistLogNotExist)
			ctx.So(a, convey.ShouldBeNil)
		})
	})
}

func TestAssistLogs(t *testing.T) {
	var (
		c      = context.Background()
		mid    = int64(0)
		start  = int(0)
		offset = int(0)
	)
	convey.Convey("Logs", t, func(ctx convey.C) {
		logs, err := d.Logs(c, mid, start, offset)
		ctx.Convey("Then err should be nil.logs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(logs, convey.ShouldNotBeNil)
		})
	})
}

func TestAssistLogsByAssist(t *testing.T) {
	var (
		c         = context.Background()
		mid       = int64(0)
		assistMid = int64(0)
		start     = int(0)
		offset    = int(0)
	)
	convey.Convey("LogsByAssist", t, func(ctx convey.C) {
		logs, err := d.LogsByAssist(c, mid, assistMid, start, offset)
		ctx.Convey("Then err should be nil.logs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(logs, convey.ShouldNotBeNil)
		})
	})
}

func TestAssistLogsByCtime(t *testing.T) {
	var (
		c      = context.Background()
		mid    = int64(0)
		stime  = time.Now()
		etime  = time.Now()
		start  = int(0)
		offset = int(0)
	)
	convey.Convey("LogsByCtime", t, func(ctx convey.C) {
		logs, err := d.LogsByCtime(c, mid, stime, etime, start, offset)
		ctx.Convey("Then err should be nil.logs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(logs, convey.ShouldNotBeNil)
		})
	})
}

func TestAssistLogsByAssistCtime(t *testing.T) {
	var (
		c         = context.Background()
		mid       = int64(0)
		assistMid = int64(0)
		stime     = time.Now()
		etime     = time.Now()
		start     = int(0)
		offset    = int(0)
	)
	convey.Convey("LogsByAssistCtime", t, func(ctx convey.C) {
		logs, err := d.LogsByAssistCtime(c, mid, assistMid, stime, etime, start, offset)
		ctx.Convey("Then err should be nil.logs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(logs, convey.ShouldNotBeNil)
		})
	})
}

func TestAssistLogCount(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
	)
	convey.Convey("LogCount", t, func(ctx convey.C) {
		totalm, err := d.LogCount(c, mid)
		ctx.Convey("Then err should be nil.totalm should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(totalm, convey.ShouldNotBeNil)
		})
	})
}

func TestAssistLogCnt(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
	)
	convey.Convey("LogCnt", t, func(ctx convey.C) {
		count, err := d.LogCnt(c, mid)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestAssistLogCntAssist(t *testing.T) {
	var (
		c         = context.Background()
		mid       = int64(0)
		assistMid = int64(0)
	)
	convey.Convey("LogCntAssist", t, func(ctx convey.C) {
		count, err := d.LogCntAssist(c, mid, assistMid)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestAssistLogCntCtime(t *testing.T) {
	var (
		c     = context.Background()
		mid   = int64(0)
		stime = time.Now()
		etime = time.Now()
	)
	convey.Convey("LogCntCtime", t, func(ctx convey.C) {
		count, err := d.LogCntCtime(c, mid, stime, etime)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestAssistLogCntAssistCtime(t *testing.T) {
	var (
		c         = context.Background()
		mid       = int64(0)
		assistMid = int64(0)
		stime     = time.Now()
		etime     = time.Now()
	)
	convey.Convey("LogCntAssistCtime", t, func(ctx convey.C) {
		count, err := d.LogCntAssistCtime(c, mid, assistMid, stime, etime)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestAssistAssistsMidsTotal(t *testing.T) {
	var (
		c       = context.Background()
		mid     = int64(27515409)
		assmids = []int64{27515235, 27515233}
	)
	convey.Convey("AssistsMidsTotal", t, func(ctx convey.C) {
		_, err := d.AssistsMidsTotal(c, mid, assmids)
		ctx.Convey("Then err should be nil.totalm should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestAssistLogObj(t *testing.T) {
	var (
		c     = context.Background()
		mid   = int64(0)
		objID = int64(0)
		tp    = int64(0)
		act   = int64(0)
	)
	convey.Convey("LogObj", t, func(ctx convey.C) {
		a, err := d.LogObj(c, mid, objID, tp, act)
		ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(a, convey.ShouldNotBeNil)
		})
	})
}
