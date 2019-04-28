package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAddUserFilters(t *testing.T) {
	Convey("test add user rule", t, func() {
		fltMap := map[string]string{"11233": "this is comment"}
		res, err := svr.AddUserFilters(context.TODO(), 150781, 1, fltMap)
		So(err, ShouldBeNil)
		for _, v := range res {
			t.Logf("%+v", v)
		}
	})
}

func TestUserFilters(t *testing.T) {
	Convey("test get user rule", t, func() {
		rs, err := svr.UserFilters(context.TODO(), 27515256)
		So(err, ShouldBeNil)
		So(rs, ShouldNotBeEmpty)
	})
}

func TestDelUserFilters(t *testing.T) {
	Convey("test del user rule", t, func() {
		_, err := svr.DelUserFilters(context.TODO(), 27515256, []int64{12, 3, 4})
		So(err, ShouldBeNil)
	})
}

func TestAddUpFilters(t *testing.T) {
	Convey("test add user rule", t, func() {
		fltMap := map[string]string{"\\q": "this is comment"}
		err := svr.AddUpFilters(context.TODO(), 27515256, 1, fltMap)
		So(err, ShouldBeNil)
	})
}

func TestUpFilters(t *testing.T) {
	Convey("test update user rule", t, func() {
		rs, err := svr.UpFilters(context.TODO(), 10097377)
		So(err, ShouldBeNil)
		So(rs, ShouldNotBeEmpty)
	})
}

func TestEditUpFilters(t *testing.T) {
	Convey("test edit user rule", t, func() {
		_, err := svr.EditUpFilters(context.TODO(), 27515256, 1, 0, []string{"\\q", "bb"})
		So(err, ShouldBeNil)
	})
}

func TestAddGlobalFilter(t *testing.T) {
	Convey("test add global rule", t, func() {
		_, err := svr.AddGlobalFilter(context.TODO(), 1, "test")
		So(err, ShouldBeNil)
	})
}

func TestGlobalFilters(t *testing.T) {
	Convey("test global rule", t, func() {
		rs, err := svr.GlobalFilters(context.TODO())
		So(err, ShouldBeNil)
		So(rs, ShouldNotBeEmpty)
	})
}

func TestDelGlobalFilters(t *testing.T) {
	Convey("test del global rule", t, func() {
		_, err := svr.DelGlobalFilters(context.TODO(), []int64{12, 3, 4})
		So(err, ShouldBeNil)
	})
}
