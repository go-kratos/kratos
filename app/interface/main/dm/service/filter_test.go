package service

import (
	"context"
	"testing"

	"go-common/library/ecode"

	. "github.com/smartystreets/goconvey/convey"
)

// )

func TestAddUserRule(t *testing.T) {
	Convey("test add user rule", t, func() {
		_, err := svr.AddUserRule(context.TODO(), 1, 150781, []string{"aa", "bb"}, "comment")
		So(err, ShouldBeNil)
	})
}

func TestUserRules(t *testing.T) {
	Convey("test get user rule", t, func() {
		rs, err := svr.UserRules(context.TODO(), 27515256)
		So(err, ShouldBeNil)
		So(rs, ShouldNotBeEmpty)
	})
}

func TestDelUserRules(t *testing.T) {
	Convey("test del user rule", t, func() {
		_, err := svr.DelUserRules(context.TODO(), 27515256, []int64{12, 3, 4})
		So(err, ShouldBeNil)
	})
}

func TestAddGlobalRule(t *testing.T) {
	Convey("test add global rule", t, func() {
		_, err := svr.AddGlobalRule(context.TODO(), 1, "aa")
		So(err, ShouldBeNil)
	})
}

func TestGlobalRules(t *testing.T) {
	Convey("test global rule", t, func() {
		rs, err := svr.GlobalRules(context.TODO())
		So(err, ShouldBeNil)
		So(rs, ShouldNotBeEmpty)
	})
}

func TestDelGlobalRules(t *testing.T) {
	Convey("test del global rule", t, func() {
		_, err := svr.DelGlobalRules(context.TODO(), []int64{12, 3, 4})
		So(err, ShouldBeNil)
	})
}

func TestServiceFilterList(t *testing.T) {
	var (
		c              = context.TODO()
		cid, mid int64 = 0, 150781
	)
	Convey("test filter list", t, func() {
		f, err := svr.FilterList(c, mid, cid)
		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)
	})
}

func TestEditFilter(t *testing.T) {
	Convey("test insert regex filter", t, func() {
		err := svr.EditFilter(c, 0, 150781, ".*", 1, 1)
		So(err, ShouldBeNil)
	})
	Convey("test insert wrong regex filter", t, func() {
		err := svr.EditFilter(c, 0, 150781, ".[*", 1, 1)
		So(err, ShouldEqual, ecode.DMFitlerIllegalRegex)
	})
	Convey("test update filter", t, func() {
		err := svr.EditFilter(c, 0, 150781, ".*", 1, 1)
		So(err, ShouldBeNil)
	})
}
