package dao

import (
	"context"
	"testing"

	"go-common/app/interface/main/dm/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRptSearchIndex(t *testing.T) {
	Convey("", t, func() {
		s := testDao.rptSearchIndex()
		So(s, ShouldNotBeEmpty)
		t.Log(s)
	})
}

func TestSearchReport(t *testing.T) {
	var (
		c = context.TODO()
		// 		mid, aid, pn, ps int64 = 432230, 9548327, 1, 20 //pre
		mid, aid, pn, ps int64 = 27515615, 10100087, 1, 20
		upOp             int8  = 2
		states                 = []int64{0, 2}
	)
	Convey("", t, func() {
		res, err := testDao.SearchReport(c, mid, aid, pn, ps, upOp, states)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
		t.Logf("%+v", res.Page)
		for _, v := range res.Result {
			t.Logf("%+v", v)
		}
	})
}

func TestSearchReportAid(t *testing.T) {
	var (
		c                 = context.TODO()
		mid, pn, ps int64 = 27515256, 1, 20
		upOp        int8
		state       = []int8{0, 2}
	)
	Convey("", t, func() {
		aids, err := testDao.SearchReportAid(c, mid, upOp, state, pn, ps)
		So(err, ShouldBeNil)
		t.Logf("%+v", aids)
	})
}

func TestUptSearchReport(t *testing.T) {
	var (
		c   = context.TODO()
		upt = &model.UptSearchReport{
			DMid:  1958334770970627,
			Upop:  0,
			Ctime: "2018-07-06 09:20:57",
			Mtime: "2018-07-27 19:20:50",
		}
	)
	Convey("", t, func() {
		err := testDao.UpdateSearchReport(c, []*model.UptSearchReport{upt})
		So(err, ShouldBeNil)
	})
}
