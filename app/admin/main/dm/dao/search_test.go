package dao

import (
	"context"
	"testing"

	"go-common/app/admin/main/dm/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSearchMonitor(t *testing.T) {
	var (
		c                           = context.TODO()
		aid, cid, mid, p, ps int64  = 0, 0, 0, 1, 50
		kw, sort, order      string = "", "", ""
		attr                        = int32(0)
		tp                          = int32(1)
	)
	Convey("get monitor list from search", t, func() {
		res, err := testDao.SearchMonitor(c, tp, aid, cid, mid, attr, kw, sort, order, p, ps)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
		t.Logf("res:%+v", res)
		for _, rpt := range res.Result {
			t.Logf("======%+v", rpt)
		}
	})
}
func TestSearchReports(t *testing.T) {
	var (
		c     = context.TODO()
		tid   = []int64{}
		rpID  = []int64{1, 2}
		upOp  = []int64{0, 1, 2}
		state = []int64{0, 1, 2}
		rt    = &model.Report{
			Aid:    -1,
			UID:    -1,
			RpUID:  -1,
			RpType: -1,
			Cid:    -1,
		}
	)
	Convey("", t, func() {
		res, err := testDao.SearchReport(c, 1, 100, "", "", "rp_time", "desc", "", tid, rpID, state, upOp, rt)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
		t.Logf("======%+v", res.Page)
		for _, rpt := range res.Result {
			t.Logf("======%+v", rpt)
		}
	})
}

func TestSearchReportByID(t *testing.T) {
	var (
		c     = context.TODO()
		dmids = []int64{719218372}
	)
	Convey("", t, func() {
		res, err := testDao.SearchReportByID(c, dmids)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
		t.Logf("===res:%v", res)
		for _, rpt := range res.Result {
			t.Logf("======%+v", rpt)
		}
	})
}

func TestSearchDM(t *testing.T) {
	var (
		c = context.TODO()
		d = &model.SearchDMParams{
			Type:         1,
			Oid:          10131232,
			Mid:          model.CondIntNil,
			ProgressFrom: model.CondIntNil,
			ProgressTo:   model.CondIntNil,
			CtimeFrom:    model.CondIntNil,
			CtimeTo:      model.CondIntNil,
			State:        "",
			Pool:         "",
			Attrs:        "1",
			Page:         1,
			Order:        "id",
			Sort:         "asc",
		}
	)
	Convey("", t, func() {
		res, err := testDao.SearchDM(c, d)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
		t.Logf("===page:%+v", res.Page)
		for _, rpt := range res.Result {
			t.Logf("======%+v", rpt)
		}
	})
}

func TestSendJudgement(t *testing.T) {
	var (
		c  = context.TODO()
		li = make([]*model.ReportJudge, 0)
	)
	i := &model.ReportJudge{
		MID:        150781,
		Operator:   "zhang",
		OperID:     121,
		OContent:   "test dm",
		OTitle:     "test archive",
		OType:      2,
		OURL:       "https://www.bilibili.com",
		ReasonType: 2,
		AID:        11,
		OID:        22,
		RPID:       33,
		Page:       1,
		BTime:      int64(1517824276),
	}
	li = append(li, i)
	Convey("test send judgement", t, func() {
		err := testDao.SendJudgement(c, li)
		So(err, ShouldBeNil)
	})
}

func TestUpSearchDMState(t *testing.T) {
	var (
		tp    int32 = 1
		state int32 = 2
		dmids       = map[int64][]int64{10131232: {1909860427366403, 1909932482887683}}
	)
	Convey("", t, func() {
		err := testDao.UpSearchDMState(context.TODO(), tp, state, dmids)
		So(err, ShouldBeNil)
	})
}

func TestUpSearchDMAttr(t *testing.T) {
	var (
		tp    int32 = 1
		oid   int64 = 10131232
		attr  int32 = 1
		dmids       = []int64{1909860427366403, 1909932482887683}
	)
	Convey("", t, func() {
		err := testDao.UpSearchDMAttr(context.TODO(), tp, oid, attr, dmids)
		So(err, ShouldBeNil)
	})
}

func TestUpSearchDMPool(t *testing.T) {
	var (
		tp    int32 = 1
		oid   int64 = 10131232
		pool  int32 = 1
		dmids       = []int64{1909860427366403, 1909932482887683}
	)
	Convey("", t, func() {
		err := testDao.UpSearchDMPool(context.TODO(), tp, oid, pool, dmids)
		So(err, ShouldBeNil)
	})
}

func TestUptSearchReport(t *testing.T) {
	var (
		uptRpt = &model.UptSearchReport{
			DMid:  719218991,
			Ctime: "2018-04-27 11:06:46",
			Mtime: "2018-06-07 11:06:46",
			State: model.StatSecondIgnore,
		}
		uptRpts = []*model.UptSearchReport{uptRpt}
	)
	Convey("test update search report", t, func() {
		err := testDao.UptSearchReport(context.TODO(), uptRpts)
		So(err, ShouldBeNil)
	})
}

func TestDaoSearchSubjectLog(t *testing.T) {
	Convey("SearchSubjectLog", t, func() {
		data, err := testDao.SearchSubjectLog(context.TODO(), 1, 1221)
		So(err, ShouldBeNil)
		So(data, ShouldNotBeNil)
		for _, v := range data {
			t.Logf("====%+v", v)
		}
	})
}

func TestDaoSearchSubject(t *testing.T) {
	r := &model.SearchSubjectReq{
		//	Oids:  []int64{10131821},
		Mids:  []int64{0},
		Sort:  "desc",
		Order: "mtime",
		Pn:    1,
		Ps:    5,
		State: int64(model.CondIntNil),
	}
	Convey("SearchSubject", t, func() {
		data, page, err := testDao.SearchSubject(context.TODO(), r)
		So(err, ShouldBeNil)
		for _, v := range data {
			t.Logf("====%+v", v)
		}
		t.Logf("====%+v", page)
	})
}

func TestSearchProtectCount(t *testing.T) {
	Convey("get protect dm count", t, func() {
		count, err := testDao.SearchProtectCount(context.TODO(), 1, 1221)
		So(err, ShouldBeNil)
		t.Log(count)
	})
}

func TestUpSearchRecemtDMState(t *testing.T) {
	var (
		tp    int32 = 1
		state int32 = 2
		dmids       = map[int64][]int64{10131232: {1909860427366403, 1909932482887683}}
	)
	Convey("", t, func() {
		err := testDao.UpSearchRecentDMState(context.TODO(), tp, state, dmids)
		So(err, ShouldBeNil)
	})
}

func TestUpSearchRecentDMAttr(t *testing.T) {
	var (
		tp    int32 = 1
		oid   int64 = 10131232
		attr  int32 = 1
		dmids       = []int64{1909860427366403, 1909932482887683}
	)
	Convey("", t, func() {
		err := testDao.UpSearchRecentDMAttr(context.TODO(), tp, oid, attr, dmids)
		So(err, ShouldBeNil)
	})
}

func TestUpSearchRecentDMPool(t *testing.T) {
	var (
		tp    int32 = 1
		oid   int64 = 10131232
		pool  int32 = 1
		dmids       = []int64{1909860427366403, 1909932482887683}
	)
	Convey("", t, func() {
		err := testDao.UpSearchRecentDMPool(context.TODO(), tp, oid, pool, dmids)
		So(err, ShouldBeNil)
	})
}
