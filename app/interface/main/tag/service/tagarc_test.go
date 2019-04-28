package service

import (
	"context"
	"testing"

	"go-common/app/service/main/archive/api"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewArcs(t *testing.T) {
	var (
		tp          int8
		c                 = context.Background()
		ps                = 1
		pn                = 10
		prid        int64 = 17
		mid         int64 = 14771787
		tid         int64 = 4052445
		rid         int32 = 17
		invalidArcs []*api.Arc
	)
	Convey("TestNewArcs", t, func() {
		testSvc.NewArcs(c, tid, ps, pn)
	})
	Convey("TestRegionNewArcs ", t, func() {
		testSvc.RegionNewArcs(c, rid, tid, tp, ps, pn)
	})
	Convey("TestDetailRankArc ", t, func() {
		testSvc.DetailRankArc(c, tid, prid, pn, ps)
	})
	Convey("TestDetail", t, func() {
		testSvc.Detail(c, tid, mid, pn, ps)
	})
	Convey("TestDelInvalidArc ", t, func() {
		testSvc.delInvalidArc(c, tid, rid, invalidArcs)
	})
}
