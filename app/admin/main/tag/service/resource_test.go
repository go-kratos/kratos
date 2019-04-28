package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServiceResource(t *testing.T) {
	var (
		oid int64 = 22786099
		pn  int32 = 1
		ps  int32 = 20
		tp  int32 = 3
		op  int32 = 1
	)
	Convey("ResourceByOid", func() {
		testSvc.ResourceByOid(context.TODO(), oid, tp)
	})
	Convey("ResByOperate", func() {
		testSvc.ResByOperate(context.TODO(), op, pn, ps)
	})
	Convey("UpdateResLimitState", func() {
		testSvc.UpdateResLimitState(context.TODO(), oid, tp, op)
	})
}
