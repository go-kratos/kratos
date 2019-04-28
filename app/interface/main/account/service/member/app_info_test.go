package member

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_Account(t *testing.T) {
	Convey("get account info", t, func() {
		Convey("when not timeout", func() {
			time.Sleep(time.Second * 2)
			res, err := s.Account(context.TODO(), 1, "127.0.0.1")
			So(err, ShouldBeNil)
			So(res, ShouldNotBeNil)
		})

	})
}

func TestService_UpdateUname(t *testing.T) {
	Convey("update uname", t, func() {
		Convey("when not timeout", func() {
			time.Sleep(time.Second * 2)
			err := s.UpdateName(context.TODO(), 110001353, "127.0.0.1", "127.0.0.1")
			So(err, ShouldBeNil)
		})

	})
}

func TestService_NickFree(t *testing.T) {
	Convey("update uname", t, func() {
		Convey("when not timeout", func() {
			time.Sleep(time.Second * 2)
			nick, err := s.NickFree(context.TODO(), 110001353)
			So(err, ShouldBeNil)
			So(nick, ShouldNotBeEmpty)
		})

	})
}

func TestService_UpdateSign(t *testing.T) {
	Convey("update sign", t, func() {
		Convey("when not timeout", func() {
			err := s.UpdateSign(context.TODO(), 61, "1989-09-19")
			So(err, ShouldBeNil)
		})
	})
}

func TestService_UpdateSex(t *testing.T) {
	Convey("update sex", t, func() {
		Convey("when not timeout", func() {
			err := s.UpdateSex(context.TODO(), 110001353, 1)
			So(err, ShouldBeNil)
		})
	})
}

func TestService_UpdateBirthday(t *testing.T) {
	Convey("update sex", t, func() {
		Convey("when not timeout", func() {
			err := s.UpdateBirthday(context.TODO(), 61, "1989-09-19")
			So(err, ShouldBeNil)
		})
	})
}
