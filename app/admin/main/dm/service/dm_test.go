package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/dm/model"

	"github.com/davecgh/go-spew/spew"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDMSearch(t *testing.T) {
	d := &model.SearchDMParams{
		Type:         1,
		Oid:          1,
		Mid:          model.CondIntNil,
		ProgressFrom: model.CondIntNil,
		ProgressTo:   model.CondIntNil,
		CtimeFrom:    model.CondIntNil,
		CtimeTo:      model.CondIntNil,
		State:        "",
		Pool:         "",
		Page:         1,
		Order:        "id",
		Sort:         "asc",
	}
	Convey("test dm list", t, func() {
		res, err := svr.DMSearch(context.TODO(), d)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		spew.Dump(res)
		So(res.Result, ShouldNotBeEmpty)
	})
}

func TestEditDMPool(t *testing.T) {
	Convey("test change pool id", t, func() {
		err := svr.EditDMPool(context.TODO(), 1, 1, 1, []int64{1, 2}, 123)
		So(err, ShouldBeNil)
	})
}

func TestXMLCacheFlush(t *testing.T) {
	Convey("test flush cache", t, func() {
		svr.XMLCacheFlush(context.TODO(), 1, 1221)
	})
}

func TestEditDMState(t *testing.T) {
	dmids := []int64{1, 2}
	Convey("test content status", t, func() {
		res := svr.EditDMState(context.TODO(), 1, 1221, 1, 1, dmids, 10, 123, "admin", "test")
		So(res, ShouldNotBeNil)
	})
}

func TestEditDMAttr(t *testing.T) {
	Convey("test change attr", t, func() {
		err := svr.EditDMAttr(context.TODO(), 1, 1, []int64{1, 2}, model.AttrProtect, 1, 123)
		So(err, ShouldBeNil)
	})
}

func TestDMIndexInfo(t *testing.T) {
	var cid int64 = 9967205
	Convey("test dm index info", t, func() {
		idx, err := svr.DMIndexInfo(context.TODO(), cid)
		So(err, ShouldBeNil)
		So(idx, ShouldNotBeNil)
	})
}
