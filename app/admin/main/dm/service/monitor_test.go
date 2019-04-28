package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/dm/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMonitorList(t *testing.T) {
	var (
		c                          = context.TODO()
		aid, cid, mid, p, ps int64 = 0, 0, 0, 1, 100
		kw, sort, order            = "", "", ""
		state                int32 = 1
		tp                   int32 = 1
	)
	Convey("test monitor list from search", t, func() {
		res, err := svr.MonitorList(c, tp, aid, cid, mid, state, kw, sort, order, p, ps)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestUpdateMonitor(t *testing.T) {
	var (
		c    = context.TODO()
		tp   = model.SubTypeVideo
		oids = []int64{2, 99}
	)
	Convey("update dm subject monitor", t, func() {
		affect, err := svr.UpdateMonitor(c, tp, oids, 1)
		So(affect, ShouldBeGreaterThan, 0)
		So(err, ShouldBeNil)
	})
}
