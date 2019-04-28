package service

import (
	"context"
	"testing"

	"go-common/app/interface/main/dm2/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDmNormalIds(t *testing.T) {
	var (
		tp    int32 = 1
		oid   int64 = 19
		cnt   int64 = 1
		n     int64 = 2
		ps    int64
		pe    int64 = 1000
		limit int64 = 100
	)
	Convey("test edit dm state", t, func() {
		dmids, err := svr.dmNormalIds(context.Background(), tp, oid, cnt, n, ps, pe, limit)
		So(err, ShouldBeNil)
		So(dmids, ShouldNotBeNil)
	})
}

func TestSegmentInfo(t *testing.T) {
	var (
		tp  int32 = 1
		aid int64 = 10097265
		oid int64 = 1508
		ps  int64 = 10
		c         = context.TODO()
		seg *model.Segment
		err error
	)
	if seg, err = svr.segmentInfo(c, tp, aid, oid, ps); err != nil {
		t.Fatalf("s.segmentInfo(%d, %d) error(%v)", aid, oid, err)
	}
	t.Logf("aid:%d, oid:%d, segment:%+v", aid, oid, seg)
}

func TestService_JudgeDMList(t *testing.T) {
	var (
		cid, dmid int64 = 9967369, 719232849
	)
	dm, err := svr.JudgeDms(context.TODO(), 1, cid, dmid)
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range dm.List {
		t.Logf("%v", m)
	}
}
