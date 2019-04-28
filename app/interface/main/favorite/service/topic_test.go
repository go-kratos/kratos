package service

import (
	"context"
	"go-common/library/ecode"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_FavTopics(t *testing.T) {
	Convey("FavTopics", t, func() {
		var (
			mid int64 = 88888894
			pn        = 1
			ps        = 30
		)
		res, err := s.FavTopics(context.TODO(), mid, pn, ps, nil)
		t.Logf("res:%v", res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func Test_IsTopicFavoured(t *testing.T) {
	Convey("IsTopicFavoured", t, func() {
		var (
			mid  int64 = 88888894
			tpID int64 = 3456
		)
		res, err := s.IsTopicFavoured(context.TODO(), mid, tpID)
		t.Logf("res:%v", res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func Test_AddFavTopic(t *testing.T) {
	Convey("AddFavTopic", t, func() {
		var (
			mid    int64 = 88888894
			tpID   int64 = 3456
			ck, ak string
		)
		err := s.AddFavTopic(context.TODO(), mid, tpID, ck, ak)
		t.Logf("err:%v", err)
		So(err, ShouldEqual, ecode.FavTopicExist)
	})
}

func Test_DelFavTopic(t *testing.T) {
	Convey("DelFavTopic", t, func() {
		var (
			mid  int64 = 88888894
			tpID int64 = 3457
		)
		err := s.DelFavTopic(context.TODO(), mid, tpID)
		t.Logf("err:%v", err)
		So(err, ShouldEqual, ecode.FavVideoAlreadyDel)
	})
}
