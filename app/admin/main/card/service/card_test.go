package service

import (
	"context"
	"flag"
	"testing"
	"time"

	"go-common/app/admin/main/card/conf"
	"go-common/app/admin/main/card/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	c = context.TODO()
	s *Service
)

func init() {
	var (
		err error
	)
	flag.Set("conf", "../cmd/test.toml")
	if err = conf.Init(); err != nil {
		panic(err)
	}
	c = context.Background()
	if s == nil {
		s = New(conf.Conf)
	}
	time.Sleep(time.Second)
}

// go test  -test.v -test.run TestCardsByGid
func TestCardsByGid(t *testing.T) {
	Convey("TestCardsByGid ", t, func() {
		card, err := s.CardsByGid(c, 2)
		t.Logf("v(%v)", card)
		So(err, ShouldBeNil)
	})
}

func TestUpdateCardState(t *testing.T) {
	Convey("TestUpdateCardState ", t, func() {
		err := s.UpdateCardState(c, &model.ArgState{ID: 1, State: 1})
		So(err, ShouldBeNil)
	})
}

func TestDeleteCard(t *testing.T) {
	Convey("TestDeleteCard ", t, func() {
		err := s.DeleteCard(c, 1)
		So(err, ShouldBeNil)
	})
}

func TestUpdateGroupState(t *testing.T) {
	Convey("TestUpdateGroupState ", t, func() {
		err := s.UpdateGroupState(c, &model.ArgState{ID: 2, State: 1})
		So(err, ShouldBeNil)
	})
}

func TestDeleteGroup(t *testing.T) {
	Convey("TestDeleteGroup ", t, func() {
		err := s.DeleteGroup(c, 2)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestGroupList
func TestGroupList(t *testing.T) {
	Convey("TestGroupList ", t, func() {
		card, err := s.GroupList(c, &model.ArgQueryGroup{GroupID: 2})
		t.Logf("v(%v)", card)
		So(err, ShouldBeNil)
	})
}

func TestCardOrderChange(t *testing.T) {
	Convey("TestCardOrderChange ", t, func() {
		err := s.CardOrderChange(c, &model.ArgIds{Ids: []int64{2, 3}})
		So(err, ShouldBeNil)
	})
}

func TestGroupOrderChange(t *testing.T) {
	Convey("TestGroupOrderChange ", t, func() {
		err := s.GroupOrderChange(c, &model.ArgIds{Ids: []int64{1}})
		So(err, ShouldBeNil)
	})
}

func TestAddGroup(t *testing.T) {
	Convey("TestAddGroup ", t, func() {
		err := s.AddGroup(c, &model.AddGroup{Name: "test17", State: 0})
		So(err, ShouldBeNil)
	})
}
