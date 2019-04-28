package service

import (
	"context"
	"testing"

	"go-common/app/job/main/dm/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFlushTrimQueue(t *testing.T) {
	var (
		tp  int32 = 1
		oid int64 = 1
	)
	Convey("", t, func() {
		err := testSvc.flushTrimQueue(context.TODO(), tp, oid)
		So(err, ShouldBeNil)
	})
}

func TestAddTrimQueue(t *testing.T) {
	var (
		tp       int32 = 1
		oid      int64 = 1
		maxlimit int64 = 1
		idx            = &model.DM{ID: 1, Type: tp, Oid: oid, Mid: 1, Progress: 1, State: 0, Pool: 2, Attr: 1}
	)
	Convey("", t, func() {
		err := testSvc.addTrimQueue(context.TODO(), tp, oid, maxlimit, idx)
		So(err, ShouldBeNil)
	})
}

func TestRecoverDM(t *testing.T) {
	var (
		tp       int32 = 1
		oid      int64 = 1
		duration int64 = 10
		maxlimit int64 = 1
		sub            = &model.Subject{ID: 1, Type: tp, Oid: oid, ACount: 2, Count: 2, Maxlimit: maxlimit}
	)
	Convey("", t, func() {
		_, err := testSvc.recoverDM(context.TODO(), sub.Type, sub.Oid, duration)
		So(err, ShouldBeNil)
	})
}

func TestSubject(t *testing.T) {
	var (
		tp  int32 = 1
		oid int64 = 1
	)
	Convey("", t, func() {
		sub, err := testSvc.subject(context.TODO(), tp, oid)
		So(err, ShouldBeNil)
		So(sub, ShouldNotBeNil)
	})
}
