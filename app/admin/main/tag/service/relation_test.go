package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServiceRelation(t *testing.T) {
	var (
		tp    int32 = 3
		pn    int32 = 1
		ps    int32 = 10
		oid   int64 = 22786099
		mid   int64 = 35152246
		tid   int64 = 233
		tname       = "unit test"
	)
	Convey("RelationListByTag", func() {
		testSvc.RelationListByTag(context.TODO(), tname, pn, ps)
	})
	Convey("RelationListByOid", func() {
		testSvc.RelationListByOid(context.TODO(), oid, tp, pn, ps)
	})
	Convey("RelationAdd", func() {
		testSvc.RelationAdd(context.TODO(), tname, oid, mid, tp)
	})
	Convey("RelationLock", func() {
		testSvc.RelationLock(context.TODO(), tid, oid, tp)
	})
	Convey("RelationUnLock", func() {
		testSvc.RelationUnLock(context.TODO(), tid, oid, tp)
	})
	Convey("RelationDelete", func() {
		testSvc.RelationDelete(context.TODO(), tid, oid, tp)
	})
}
