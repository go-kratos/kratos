package monitor

import (
	"context"
	"go-common/app/job/main/aegis/model/monitor"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestMonitorAddToSet(t *testing.T) {
	convey.Convey("AddToSet", t, func(convCtx convey.C) {
		var (
			c    = context.Background()
			keys = []string{"monitor_test"}
			oid  = int64(123)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			_, err := d.AddToSet(c, keys, oid)
			convCtx.Convey("Then err should be nil.logs should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMonitorRemFromSet(t *testing.T) {
	convey.Convey("RemFromSet", t, func(convCtx convey.C) {
		var (
			c    = context.Background()
			keys = []string{"monitor_test"}
			oid  = int64(123)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			_, err := d.RemFromSet(c, keys, oid)
			convCtx.Convey("Then err should be nil.logs should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMonitorClearExpireSet(t *testing.T) {
	convey.Convey("ClearExpireSet", t, func(convCtx convey.C) {
		var (
			c    = context.Background()
			keys = []string{"monitor_test"}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			_, err := d.ClearExpireSet(c, keys)
			convCtx.Convey("Then err should be nil.logs should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMonitorAddToDelArc(t *testing.T) {
	convey.Convey("AddToDelArc", t, func(convCtx convey.C) {
		var (
			c = context.Background()
			a = &monitor.BinlogArchive{
				ID: 123,
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.AddToDelArc(c, a)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMonitorArcDelInfos(t *testing.T) {
	convey.Convey("ArcDelInfos", t, func(convCtx convey.C) {
		var (
			c    = context.Background()
			aids = []int64{123}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			_, err := d.ArcDelInfos(c, aids)
			convCtx.Convey("Then err should be nil.infos should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMonitorMoniRuleStats(t *testing.T) {
	convey.Convey("MoniRuleStats", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			id  = int64(1)
			min = int64(0)
			max = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			_, err := d.MoniRuleStats(c, id, min, max)
			convCtx.Convey("Then err should be nil.stats should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMonitorMoniRuleOids(t *testing.T) {
	convey.Convey("MoniRuleOids", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			id  = int64(1)
			min = int64(0)
			max = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			_, err := d.MoniRuleOids(c, id, min, max)
			convCtx.Convey("Then err should be nil.oidMap should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
