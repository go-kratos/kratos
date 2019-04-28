package card

import (
	"context"
	"flag"
	"testing"

	"go-common/app/interface/main/account/conf"
	v1 "go-common/app/service/main/card/api/grpc/v1"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
	c = context.Background()
)

func init() {
	flag.Set("conf", "../../cmd/account-interface-example.toml")
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	s = New(conf.Conf)
}

// go test  -test.v -test.run TestServiceUserCard
func TestServiceUserCard(t *testing.T) {
	Convey("TestServiceUserCard", t, func() {
		res, err := s.UserCard(c, 1)
		t.Logf("%v", res)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestServiceCard
func TestServiceCard(t *testing.T) {
	Convey("TestServiceCard", t, func() {
		res, err := s.Card(c, 1)
		t.Logf("%v", res)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestCardHots
func TestCardHots(t *testing.T) {
	Convey("TestCardHots", t, func() {
		res, err := s.CardHots(c)
		t.Logf("%v", res)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestAllGroup
func TestAllGroup(t *testing.T) {
	Convey("TestAllGroup", t, func() {
		res, err := s.AllGroup(c, 1)
		t.Logf("%v", res)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestCardsByGid
func TestCardsByGid(t *testing.T) {
	Convey("TestCardsByGid", t, func() {
		res, err := s.CardsByGid(c, 1)
		t.Logf("%v", res)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestEquip
func TestEquip(t *testing.T) {
	Convey("TestEquip", t, func() {
		err := s.Equip(c, &v1.ModelArgEquip{
			Mid:    2,
			CardId: 1,
		})
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestDemount
func TestDemount(t *testing.T) {
	Convey("TestDemount", t, func() {
		err := s.Demount(c, 2)
		So(err, ShouldBeNil)
	})
}
