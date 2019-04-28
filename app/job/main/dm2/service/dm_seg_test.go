package service

import (
	"context"
	"testing"

	"go-common/app/job/main/dm2/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDmsByid(t *testing.T) {
	var (
		tp     int32 = 1
		oid    int64 = 1221
		missed       = []int64{719150141, 719150142}
	)
	Convey("", t, func() {
		dms, err := svr.dmsByid(context.TODO(), tp, oid, missed)
		So(err, ShouldBeNil)
		So(dms, ShouldNotBeEmpty)
		for _, dm := range dms {
			t.Log(dm)
		}
	})
}

func TestDMSeg(t *testing.T) {
	var (
		tp        int32 = 1
		oid       int64 = 1221
		childpool int32 = 1
		limit     int64 = 10
		p               = &model.Page{Num: 1, Size: model.DefaultVideoEnd, Total: 1}
	)
	Convey("", t, func() {
		res, err := svr.dmSeg(context.TODO(), tp, oid, limit, childpool, p)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		t.Logf("%v,length:%d", res, len(res.Elems))
	})
}

func TestPageInfo(t *testing.T) {
	Convey("", t, func() {
		dm := &model.DM{ID: 719182141, Type: 1, Oid: 1221, Progress: 0, Pool: 2}
		p, err := svr.pageinfo(context.TODO(), 12345, dm)
		So(err, ShouldBeNil)
		t.Log(p)
	})
}
