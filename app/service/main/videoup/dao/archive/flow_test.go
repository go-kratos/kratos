package archive

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_TxAddFlow(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
	)
	Convey("TxAddFlow", t, func(ctx C) {
		_, err := d.TxAddFlow(tx, 123, 23333, 0, 0, "ut")
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestDao_TxUpFlowState(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
	)
	Convey("TxUpFlowState", t, func(ctx C) {
		_, err := d.TxUpFlowState(tx, 0, 23333)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestDao_CheckFlowGroupID(t *testing.T) {
	var (
		c = context.Background()
	)
	Convey("CheckFlowGroupID", t, func(ctx C) {
		_, err := d.CheckFlowGroupID(c, 0, 23333, 24)
		So(err, ShouldBeNil)
	})
}

func TestDao_FlowPage(t *testing.T) {
	var (
		c = context.Background()
	)
	Convey("FlowPage", t, func(ctx C) {
		_, err := d.FlowPage(c, 0, 23333, 1, 10)
		So(err, ShouldBeNil)
	})
}

func TestDao_TxAddFlowLog(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
	)
	Convey("TxAddFlowLog", t, func(ctx C) {
		_, err := d.TxAddFlowLog(tx, 0, 123, 0, 24, "")
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestArchiveFlows(t *testing.T) {
	var (
		c = context.Background()
	)
	Convey("Flows", t, func(ctx C) {
		_, err := d.Flows(c)
		ctx.Convey("Then err should be nil.fs should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}

func TestFlowsByOID(t *testing.T) {
	var (
		c = context.Background()
	)
	Convey("FlowsByOID", t, func(ctx C) {
		_, err := d.FlowsByOID(c, 2880441)
		ctx.Convey("Then err should be nil.fs should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}

func TestFlowUnique(t *testing.T) {
	var (
		c = context.Background()
	)
	Convey("FlowUnique", t, func(ctx C) {
		_, err := d.FlowUnique(c, 2880441, 11, 0)
		ctx.Convey("Then err should be nil.fs should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}

func TestFlowGroupPools(t *testing.T) {
	var (
		c = context.Background()
	)
	Convey("FlowGroupPools", t, func(ctx C) {
		_, err := d.FlowGroupPools(c, []int64{0})
		ctx.Convey("Then err should be nil.fs should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}

func TestFlowPage(t *testing.T) {
	var (
		c = context.Background()
	)
	Convey("FlowPage", t, func(ctx C) {
		_, err := d.FlowPage(c, 0, 11, 1, 10)
		ctx.Convey("Then err should be nil.fs should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}

func TestFindGroupIDByScope(t *testing.T) {
	var (
		c = context.Background()
	)
	Convey("FindGroupIDByScope", t, func(ctx C) {
		_, err := d.FindGroupIDByScope(c, 0, 11, 1, 1)
		ctx.Convey("Then err should be nil.fs should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}

func TestArchiveWhiteMids(t *testing.T) {
	var (
		c = context.Background()
	)
	Convey("WhiteMids", t, func(ctx C) {
		_, err := d.WhiteMids(c)
		ctx.Convey("Then err should be nil.mids should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}

func TestArchiveCheckActGroupID(t *testing.T) {
	var (
		c       = context.Background()
		groupID = int64(1)
	)
	Convey("CheckActGroupID", t, func(ctx C) {
		_, err := d.CheckActGroupID(c, groupID)
		ctx.Convey("Then err should be nil.state should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}

func TestArchiveCheckFlowMid(t *testing.T) {
	var (
		c    = context.Background()
		pool = int8(1)
		oid  = int64(2)
	)
	Convey("CheckFlowMid", t, func(ctx C) {
		_, err := d.CheckFlowMid(c, pool, oid)
		ctx.Convey("Then err should be nil.flows should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}

func TestArchiveForbidMids(t *testing.T) {
	var (
		c = context.Background()
	)
	Convey("ForbidMids", t, func(ctx C) {
		_, err := d.ForbidMids(c)
		ctx.Convey("Then err should be nil.mids should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}

func TestArchiveAppFeedAids(t *testing.T) {
	var (
		c         = context.Background()
		startTime = time.Now()
		endTime   = time.Now()
	)
	Convey("AppFeedAids", t, func(ctx C) {
		_, err := d.AppFeedAids(c, startTime, endTime)
		ctx.Convey("Then err should be nil.aids should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}
