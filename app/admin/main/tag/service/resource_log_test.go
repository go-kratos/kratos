package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServiceResourceLog(t *testing.T) {
	var (
		id     int64 = 123456
		oid    int64 = 22786099
		pn     int32 = 1
		ps     int32 = 20
		tp     int32 = 3
		role   int32 = 2
		action int32 = 1
		state  int32 = 1
	)
	Convey("ResourceLogs", func() {
		testSvc.ResourceLogs(context.TODO(), oid, tp, role, action, pn, ps)
	})
	Convey("UpdateResLogState", func() {
		testSvc.UpdateResLogState(context.TODO(), id, oid, tp, state)
	})
}
