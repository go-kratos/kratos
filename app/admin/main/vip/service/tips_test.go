package service

import (
	"testing"

	"go-common/app/admin/main/vip/model"

	. "github.com/smartystreets/goconvey/convey"
)

// go test  -test.v -test.run TestServiceTipList
func TestServiceTipList(t *testing.T) {
	Convey("TestServiceTipList", t, func() {
		var (
			platform = int8(0)
			state    = int8(0)
			position = int8(2)
		)
		res, err := s.TipList(c, platform, state, position)
		for _, v := range res {
			t.Logf("%+v", v)
		}
		So(len(res) != 0, ShouldBeTrue)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestServiceTipByID
func TestServiceTipByID(t *testing.T) {
	Convey("TestServiceTipByID", t, func() {
		var (
			id int64 = 1
		)
		res, err := s.TipByID(c, id)
		t.Logf("%+v", res)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestServiceAddTip
func TestServiceAddTip(t *testing.T) {
	Convey("TestServiceAddTip", t, func() {
		t := &model.Tips{
			Platform:  2,
			Version:   4000,
			Tip:       "一样",
			Link:      "http://www.baidu.com",
			StartTime: 1528315928,
			EndTime:   1538315928,
			Level:     2,
			JudgeType: 1,
			Operator:  "baihai",
			Position:  2,
		}
		err := s.AddTip(c, t)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestServiceUpdateTip
func TestServiceUpdateTip(t *testing.T) {
	Convey("TestServiceUpdateTip", t, func() {
		t := &model.Tips{
			ID:        1,
			Platform:  2,
			Version:   4000,
			Tip:       "一样2",
			Link:      "http://www.baidu.com",
			StartTime: 1528315928,
			EndTime:   1538315928,
			Level:     2,
			JudgeType: 1,
			Position:  1,
			Operator:  "baihai",
		}
		err := s.TipUpdate(c, t)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestServiceDeleteTip
func TestServiceDeleteTip(t *testing.T) {
	Convey("TestServiceDeleteTip", t, func() {
		var (
			id int64 = 2
		)
		err := s.DeleteTip(c, id, "baihai")
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestServiceExpireTip
func TestServiceExpireTip(t *testing.T) {
	Convey("TestServiceExpireTip", t, func() {
		var (
			id int64 = 3
		)
		err := s.ExpireTip(c, id, "baihai")
		So(err, ShouldBeNil)
	})
}
