package service

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestInfo(t *testing.T) {
	convey.Convey("Info", t, func() {
		res, err := s.Info(context.TODO(), 1)
		convey.So(res, convey.ShouldNotBeNil)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestCard(t *testing.T) {
	convey.Convey("Card", t, func() {
		res, err := s.Card(context.TODO(), 1)
		convey.So(res, convey.ShouldNotBeNil)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestInfos(t *testing.T) {
	convey.Convey("Infos", t, func() {
		res, err := s.Infos(context.TODO(), []int64{1, 2, 3})
		convey.So(res, convey.ShouldNotBeNil)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestCards(t *testing.T) {
	convey.Convey("Cards", t, func() {
		res, err := s.Cards(context.TODO(), []int64{1, 2, 3})
		convey.So(res, convey.ShouldNotBeNil)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestProfile(t *testing.T) {
	convey.Convey("Profile", t, func() {
		res, err := s.Profile(context.TODO(), 111003471)
		convey.So(res, convey.ShouldNotBeNil)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestInfosByName(t *testing.T) {
	convey.Convey("InfosByName", t, func() {
		res, err := s.InfosByName(context.TODO(), []string{"1", "2"})
		convey.So(res, convey.ShouldNotBeNil)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestProfileWithStat(t *testing.T) {
	convey.Convey("ProfileWithStat", t, func() {
		res, err := s.ProfileWithStat(context.TODO(), 111003471)
		convey.So(res, convey.ShouldNotBeNil)
		convey.So(err, convey.ShouldBeNil)
	})
}
