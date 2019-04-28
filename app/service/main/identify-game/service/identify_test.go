package service

import (
	"context"
	"sync"
	"testing"
	"time"

	"go-common/app/service/main/identify-game/conf"
	"go-common/app/service/main/identify-game/model"
	"go-common/library/ecode"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	once sync.Once
	s    *Service
	// 35位
	errKey = "123456789012345678901234567890123451"
	// 32位
	rightKey = "12345678901234567890123456789013"
)

func startService() {
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second * 2)
}

func Test_AccessToken(t *testing.T) {
	once.Do(startService)
	Convey("access token", t, func() {
		obj := &model.AccessInfo{
			Mid:     41862641,
			Token:   "xxxxxxxxxxxxxxxxxxxxx",
			Expires: 60,
		}
		err := s.d.SetAccessCache(context.Background(), obj.Token, obj)
		So(err, ShouldBeNil)

		_, err = s.Oauth(context.Background(), obj.Token, "origin")
		So(err, ShouldBeNil)

		_, err = s.Oauth(context.Background(), errKey, "origin")
		So(err, ShouldEqual, ecode.NoLogin)

		r, err := s.d.AccessCache(context.Background(), errKey)
		So(err, ShouldBeNil)
		So(r, ShouldBeNil)

		_, err = s.Oauth(context.Background(), rightKey, "origin")
		So(err, ShouldEqual, ecode.NoLogin)

		time.Sleep(time.Millisecond * 10)

		r, err = s.d.AccessCache(context.Background(), rightKey)
		So(err, ShouldBeNil)
		So(r, ShouldNotBeNil)
		So(r.Mid, ShouldEqual, _noLogin.Mid)
	})
}

func Test_RenewToken(t *testing.T) {
	once.Do(startService)
	if r, err := s.d.RenewToken(context.Background(), rightKey, "origin"); err != ecode.NoLogin {
		t.Errorf("Test_RenewToken fail. %v, %v", r, err)
		t.FailNow()
	}

	if r, err := s.d.RenewToken(context.Background(), errKey, "origin"); err != ecode.NoLogin {
		t.Errorf("Test_RenewToken fail. %v, %v", r, err)
		t.FailNow()
	}
}

func Test_Target(t *testing.T) {
	once.Do(startService)
	Convey("target parse", t, func() {
		Convey("situation 1", func() {
			res := s.target(rightKey + "_tx")
			So(res, ShouldEqual, "tx")
		})
		Convey("situation 2", func() {
			res := s.target(rightKey + "_")
			So(res, ShouldEqual, "")
		})
		Convey("situation 3", func() {
			res := s.target(rightKey)
			So(res, ShouldEqual, "origin")
		})
	})
}

func Test_DelCache(t *testing.T) {
	once.Do(startService)
	Convey("test delete cache", t, func() {
		err := s.DelCache(context.Background(), "def32243")
		So(err, ShouldBeNil)
	})
}

func Test_GetCookies(t *testing.T) {
	once.Do(startService)
	Convey("test get cookies by token", t, func() {
		cookies, err := s.GetCookieByToken(context.Background(), "def32243111", "")
		So(err, ShouldNotBeNil)
		So(cookies, ShouldBeNil)
	})
}
