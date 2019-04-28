package dao

import (
	"context"
	"testing"

	"go-common/app/interface/main/dm2/model"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSearchDM(t *testing.T) {
	p := &model.SearchDMParams{
		Oid:          10131156,
		Mids:         "",
		Keyword:      "",
		ProgressFrom: model.CondIntNil,
		ProgressTo:   model.CondIntNil,
		CtimeFrom:    "",
		CtimeTo:      "",
		Pn:           1,
		Ps:           100,
		State:        "0",
		Type:         1,
		Mode:         "",
		Sort:         "desc",
		Order:        "ctime",
		Pool:         "",
		Attrs:        "",
	}
	Convey("", t, func() {
		res, err := testDao.SearchDM(context.TODO(), p)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
		t.Logf("%+v", res.Page)
		for i, v := range res.Result {
			t.Logf("%d %v", i, v)
		}
	})
}

func TestSearchDMHisIndex(t *testing.T) {
	Convey("search dm history index", t, func() {
		dates, err := testDao.SearchDMHisIndex(context.TODO(), 1, 10109227, "2018-03")
		So(err, ShouldBeNil)
		So(dates, ShouldNotBeEmpty)
		t.Log(dates)
	})
}

func TestSearchDMHistory(t *testing.T) {
	Convey("search dm history", t, func() {
		ctime, err := time.Parse("2006-01-02 15:04:05", "2016-08-03 23:59:59")
		if err != nil {
			t.Fail()
		}
		dmids, err := testDao.SearchDMHistory(context.TODO(), 1, 1221, ctime.Unix(), 1, 100)
		So(err, ShouldBeNil)
		So(dmids, ShouldNotBeEmpty)
		t.Log(dmids)
	})
}

func TestUptSearchDMPool(t *testing.T) {
	Convey("update search dm pool", t, func() {
		err := testDao.UptSearchDMPool(context.TODO(), []int64{416894555258883, 372530664701955}, 10131156, 1, 1)
		So(err, ShouldBeNil)
	})
}

func TestUptSearchDMState(t *testing.T) {
	Convey("update search dm state", t, func() {
		err := testDao.UptSearchDMState(context.TODO(), []int64{372412118466563}, 10131156, 1, 1)
		So(err, ShouldBeNil)
	})
}

func TestUptSearchDMAttr(t *testing.T) {
	Convey("update search dm attr", t, func() {
		err := testDao.UptSearchDMAttr(context.TODO(), []int64{372412118466563}, 10131156, 0, 1)
		So(err, ShouldBeNil)
	})
}

func TestSearchSubtitles(t *testing.T) {
	Convey("search subtitle", t, func() {

		res, err := testDao.SearchSubtitles(context.TODO(), 1, 10, 0, nil, 0, 0, 0, nil)
		So(err, ShouldBeNil)
		t.Logf("page:%+v", res.Page)
		t.Logf("results:%v", len(res.Results))
		for _, rs := range res.Results {
			t.Logf("rs:%+v", rs)
		}
	})
}

func TestCountSubtitles(t *testing.T) {
	Convey("search subtitle", t, func() {

		res, err := testDao.CountSubtitles(context.TODO(), 0, nil, 0, 0, 0)
		So(err, ShouldBeNil)
		t.Logf("page:%+v", res)
	})
}
