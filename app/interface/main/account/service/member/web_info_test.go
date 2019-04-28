package member

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/interface/main/account/model"
	"testing"
	"time"
)

func TestService_SettingsInfo(t *testing.T) {
	Convey("get settingsInfo info", t, func() {
		Convey("when not timeout", func() {
			time.Sleep(time.Second * 2)
			res, err := s.SettingsInfo(context.TODO(), 110001260)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeEmpty)
		})
	})
}

func TestService_LogLogin(t *testing.T) {
	Convey("get settingsInfo info", t, func() {
		Convey("when not timeout", func() {
			time.Sleep(time.Second * 2)
			res, err := s.LogLogin(context.TODO(), 110001260)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeEmpty)
		})
	})
}

func TestService_LogCoin(t *testing.T) {
	Convey("get settingsInfo info", t, func() {
		Convey("when not timeout", func() {
			time.Sleep(time.Second * 2)
			res, err := s.LogCoin(context.TODO(), 1)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeEmpty)
		})
	})
}

func TestService_LogExp(t *testing.T) {
	Convey("get settingsInfo info", t, func() {
		Convey("when not timeout", func() {
			time.Sleep(time.Second * 2)
			res, err := s.LogExp(context.TODO(), 110001260)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeEmpty)
		})
	})
}

func TestService_LogMoral(t *testing.T) {
	Convey("get settingsInfo info", t, func() {
		Convey("when not timeout", func() {
			time.Sleep(time.Second * 2)
			res, err := s.LogMoral(context.TODO(), 110001260)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeEmpty)
		})
	})
}

func TestService_Coin(t *testing.T) {
	Convey("get settingsInfo info", t, func() {
		Convey("when not timeout", func() {
			time.Sleep(time.Second * 2)
			res, err := s.Coin(context.TODO(), 1)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeEmpty)
		})
	})
}

func TestService_UpdateSettings(t *testing.T) {
	Convey("get settingsInfo info", t, func() {
		Convey("when not timeout", func() {
			time.Sleep(time.Second * 2)
			err := s.UpdateSettings(context.TODO(), 110000092, &model.Settings{Uname: "test", Sex: "男"})
			So(err, ShouldBeNil)
		})
	})
}

func TestService_Reward(t *testing.T) {
	Convey("reward info", t, func() {
		Convey("when not timeout", func() {
			time.Sleep(time.Second * 2)
			res, err := s.Reward(context.TODO(), 2)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeEmpty)
		})
	})
}
