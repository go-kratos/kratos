package service

import (
	"context"
	"testing"

	"go-common/app/interface/main/dm/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	c       = context.TODO()
	cid     = int64(10106598)
	dmid    = int64(719918177)
	uid     = int64(1234567)
	reason  = int8(1)
	content = "aaaaaa"
)

func TestAddReport(t *testing.T) {
	Convey("test add  report", t, func() {
		id, err := svr.AddReport(c, cid, dmid, uid, reason, content)
		So(err, ShouldBeNil)
		So(id, ShouldBeGreaterThan, 0)
	})
}

func TestReportList(t *testing.T) {
	var (
		mid, aid, page, size int64 = 27515615, 0, 1, 100
		upOp                 int8
		state                = []int64{0, 2}
	)
	Convey("test report list", t, func() {
		list, err := svr.ReportList(c, mid, aid, page, size, upOp, state)
		So(err, ShouldBeNil)
		So(list, ShouldNotBeNil)
	})
}

func TestReportArchives(t *testing.T) {
	var (
		mid, pn, ps int64 = 27515256, 1, 20
		upOp        int8
		states      = []int8{0, 2}
	)
	Convey("test report archive list", t, func() {
		res, err := svr.ReportArchives(c, mid, upOp, states, pn, ps)
		So(err, ShouldBeNil)
		if res != nil {
			for _, v := range res.Result {
				t.Logf("%+v", v)
			}
		}
	})
}

func TestEditReport(t *testing.T) {
	var (
		cid  int64 = 10114205
		dmid int64 = 719218893
		mid  int64 = 27515615
		upOp       = int8(model.StateDelete)
	)
	Convey("test edit report", t, func() {
		_, err := svr.EditReport(c, 1, cid, mid, dmid, upOp)
		So(err, ShouldBeNil)
	})
}
