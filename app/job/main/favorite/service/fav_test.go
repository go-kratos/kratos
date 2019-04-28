package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_setRelationCache(t *testing.T) {
	Convey("setRelationCache", t, func() {
		var (
			typ int8  = 2
			mid int64 = 88888894
			fid int64 = 289
		)
		err := s.setRelationCache(context.TODO(), typ, mid, fid)
		t.Logf("err:%v", err)
		So(err, ShouldBeNil)
	})
}

func Test_folder(t *testing.T) {
	Convey("folder", t, func() {
		var (
			typ int8  = 1
			mid int64 = 88888894
			fid int64 = 1
		)
		res, err := s.folder(context.TODO(), typ, mid, fid)
		t.Logf("res:%v", res)
		t.Logf("err:%v", err)
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}

func Test_addCoin(t *testing.T) {
	Convey("addMoney", t, func() {
		var (
			isAdd       = true
			count       = 200
			typ   int8  = 1
			oid   int64 = 123
		)
		err := s.addCoin(context.TODO(), isAdd, count, typ, oid)
		t.Logf("err:%v", err)
		So(err, ShouldBeNil)
	})
}
