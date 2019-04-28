package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServiceTag(t *testing.T) {
	var (
		tid   int64 = 2233
		tp    int32 = 3
		state int32 = 1

		tname   = "unit test"
		content = "unit test"
	)
	Convey("TagEdit", func() {
		testSvc.TagEdit(context.TODO(), tid, tp, tname, content)
	})
	Convey("TagInfo", func() {
		testSvc.TagInfo(context.TODO(), tid)
	})
	Convey("TagState", func() {
		testSvc.TagState(context.TODO(), tid, state)
	})
	Convey("TagVerify", func() {
		testSvc.TagVerify(context.TODO(), tid)
	})
}
