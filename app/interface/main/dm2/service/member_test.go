package service

import (
	"context"
	"testing"

	"go-common/app/interface/main/dm2/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDMUpRecent(t *testing.T) {
	var (
		c         = context.TODO()
		mid int64 = 123
		pn  int64 = 1
		ps  int64 = 10
	)
	Convey("dm recent", t, func() {
		res, err := svr.DMUpRecent(c, mid, pn, ps)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestDMUpSearch(t *testing.T) {
	p := &model.SearchDMParams{
		Type:         1,
		Oid:          10131156,
		Mids:         "",
		ProgressFrom: model.CondIntNil,
		ProgressTo:   model.CondIntNil,
		CtimeFrom:    "",
		CtimeTo:      "",
		Mode:         "",
		State:        "0,2,6",
		Pool:         "",
		Pn:           1,
		Ps:           50,
		Order:        "ctime",
		Sort:         "asc",
		Keyword:      "还吃几个",
	}
	Convey("test up dm list", t, func() {
		res, err := svr.DMUpSearch(context.TODO(), 123, p)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		for _, v := range res.Result {
			t.Logf("===========\n%+v", v)
		}
		t.Logf("===========\n%+v", res)
		So(res.Result, ShouldNotBeEmpty)
	})
}
