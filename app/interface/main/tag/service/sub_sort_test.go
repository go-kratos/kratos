package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCustomSubTags(t *testing.T) {
	var (
		ps          = 1
		pn          = 10
		order       = 1
		tp          = 0
		mid   int64 = 14771787
		tids        = []int64{10176, 5578, 11436}
	)
	Convey("CustomSubTags", t, func() {
		testSvc.CustomSubTags(context.Background(), mid, order, tp, ps, pn)
	})
	Convey("UpCustomSubTags", t, func() {
		testSvc.UpCustomSubTags(context.Background(), mid, tids, tp)
	})
}
