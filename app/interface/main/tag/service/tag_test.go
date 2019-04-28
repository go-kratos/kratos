package service

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInfoByID(t *testing.T) {
	var (
		name    string
		tids    []int64
		names   []string
		rid     int64 = 1
		hotType int64 = 1
		mid     int64 = 14771787
		tid     int64 = 10176
	)
	Convey("TestInfoByID service", t, func() {
		testSvc.InfoByID(context.Background(), mid, tid)
	})
	Convey("MinfoByIDs service", t, func() {
		testSvc.MinfoByIDs(context.Background(), mid, tids)
	})
	Convey("TestInfoByName service", t, func() {
		testSvc.InfoByName(context.Background(), mid, name)
	})
	Convey("TestMinfoByNames service", t, func() {
		testSvc.MinfoByNames(context.Background(), mid, names)
	})
	Convey("TestHotTags service", t, func() {
		testSvc.HotTags(context.Background(), mid, rid, hotType)
	})
	Convey("TestSimilarTags", t, func() {
		testSvc.SimilarTags(context.Background(), rid, tid)
	})
	Convey("TestChangeSim", t, func() {
		testSvc.ChangeSim(context.Background(), tid)
	})
	Convey("TestAddActivityTag", t, func() {
		testSvc.AddActivityTag(context.Background(), "", time.Now())
	})
	Convey("TestRecommandTag", t, func() {
		testSvc.RecommandTag(context.Background())
	})
}
