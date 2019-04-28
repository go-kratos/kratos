package service

import (
	"testing"
	"time"

	"go-common/app/service/main/card/model"

	. "github.com/smartystreets/goconvey/convey"
)

// go test  -test.v -test.run TestUserCard
func TestUserCard(t *testing.T) {
	Convey("TestUserCard ", t, func() {
		card, err := s.UserCard(c, 1)
		t.Logf("v(%v)", card)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestUserCards
func TestUserCards(t *testing.T) {
	Convey("TestUserCards ", t, func() {
		card, err := s.UserCards(c, []int64{977771, 977772, 977773})
		time.Sleep(1 * time.Second)
		t.Logf("v(%v)", card)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestEquip
func TestEquip(t *testing.T) {
	Convey("TestEquip ", t, func() {
		err := s.Equip(c, &model.ArgEquip{Mid: 27515232, CardID: 1})
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestDemountEquip
func TestDemountEquip(t *testing.T) {
	Convey("TestDemountEquip ", t, func() {
		err := s.DemountEquip(c, 27515232)
		So(err, ShouldBeNil)
	})
}
