package service

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSubTags(t *testing.T) {
	var (
		ps          = 1
		pn          = 10
		order       = 1
		tid   int64 = 10176
		mid   int64 = 14771787
		vmid  int64 = 14771787
		tids        = []int64{10176, 5578, 11436}
	)
	Convey("SubTags", t, func() {
		testSvc.SubTags(context.Background(), mid, vmid, pn, ps, order)
	})
	Convey("AddSub", t, func() {
		testSvc.AddSub(context.Background(), mid, tids, time.Now())
	})
	Convey("TestCancelSub", t, func() {
		testSvc.CancelSub(context.Background(), tid, mid, time.Now())
	})
	Convey("TestSubArcs", t, func() {
		testSvc.SubArcs(context.Background(), mid, vmid)
	})
}
