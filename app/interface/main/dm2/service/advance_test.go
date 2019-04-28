package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAdvanceState(t *testing.T) {
	Convey("test adv State", t, func() {
		_, err := svr.AdvanceState(context.TODO(), 27515330, 10107292, "sp")
		So(err, ShouldBeNil)
	})
}

func TestAdvances(t *testing.T) {
	Convey("test adv", t, func() {
		res, err := svr.Advances(context.TODO(), 27515260)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestPassAdvance(t *testing.T) {
	Convey("test pass adv", t, func() {
		err := svr.PassAdvance(context.TODO(), 7158471, 2)
		So(err, ShouldBeNil)
	})
}

func TestDenyAdvance(t *testing.T) {
	Convey("test deny adv", t, func() {
		err := svr.DenyAdvance(context.TODO(), 27515615, 107)
		So(err, ShouldBeNil)
	})
}

func TestCancelAdvance(t *testing.T) {
	Convey("test cancel adv", t, func() {
		err := svr.CancelAdvance(context.TODO(), 27515615, 122)
		So(err, ShouldBeNil)
	})
}
