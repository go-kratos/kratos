package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServiceRank(t *testing.T) {
	var (
		prid  int64 = 22
		rid   int64 = 33
		tp    int32 = 3
		tname       = "unit test"
	)
	Convey("RegionRankList", func() {
		testSvc.RegionRankList(context.TODO(), prid, rid, tp)
	})
	Convey("CastRankList", func() {
		testSvc.ArchiveRankList(context.TODO(), prid, rid, tp)
	})
	Convey("OperateHotTag", func() {
		testSvc.OperateHotTag(context.TODO(), tname)
	})
	// Convey("UpdateRank", func() {
	// 	testSvc.UpdateRank(context.TODO(), tname)
	// })
}
