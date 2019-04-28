package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_UserFolders(t *testing.T) {
	Convey("UserFolders", t, func() {
		var (
			typ  int8  = 1
			mid  int64 = 88888894
			vmid int64 = 88888894
			oid  int64 = 123123
		)
		res, err := s.UserFolders(context.TODO(), typ, mid, vmid, oid, typ)
		t.Logf("res:%+v", res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func Test_recentOids(t *testing.T) {
	Convey("recentOids", t, func() {
		var (
			typ  int8  = 1
			mid  int64 = 88888894
			fids       = []int64{1}
		)
		res, err := s.recentOids(context.TODO(), typ, mid, fids)
		t.Logf("res:%+v", res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}
func Test_CntUserFolders(t *testing.T) {
	Convey("CntUserFolders", t, func() {
		var (
			typ  int8  = 1
			mid  int64 = 88888894
			vmid int64 = 88888894
		)
		res, err := s.CntUserFolders(context.TODO(), typ, mid, vmid)
		t.Logf("res:%+v", res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}
func Test_InDefaultFolder(t *testing.T) {
	Convey("InDefaultFolder", t, func() {
		var (
			typ int8  = 1
			mid int64 = 88888894
			oid int64 = 987
		)
		res, err := s.InDefaultFolder(context.TODO(), typ, mid, oid)
		t.Logf("res:%+v", res)
		So(err, ShouldBeNil)
		So(res, ShouldBeFalse)
	})
}
func Test_OidCount(t *testing.T) {
	Convey("OidCount", t, func() {
		var (
			typ int8  = 1
			oid int64 = 987
		)
		res, err := s.OidCount(context.TODO(), typ, oid)
		t.Logf("res:%+v", res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func Test_OidsCount(t *testing.T) {
	Convey("OidsCount", t, func() {
		var (
			typ  int8 = 1
			oids      = []int64{1, 2, 3}
		)
		res, err := s.OidsCount(context.TODO(), typ, oids)
		t.Logf("res:%+v", res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func Test_BatchFavs(t *testing.T) {
	Convey("BatchFavs", t, func() {
		var (
			typ   int8  = 1
			mid   int64 = 8888894
			limit       = 1000
		)
		res, err := s.BatchFavs(context.TODO(), typ, mid, limit)
		t.Logf("res:%+v", res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}
