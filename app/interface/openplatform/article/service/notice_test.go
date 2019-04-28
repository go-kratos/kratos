package service

import (
	"testing"

	"go-common/app/interface/openplatform/article/model"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Notice(t *testing.T) {
	a := &model.Notice{ID: 1, Plat: _platAll, Condition: _equal, Build: 20}
	b := &model.Notice{ID: 2, Plat: _platIOS, Condition: _greaterThanOrEqual, Build: 30}
	c := &model.Notice{ID: 3, Plat: _platAndroid, Condition: _lessThanOrEqual, Build: 50}
	s.notices = []*model.Notice{a, b, c}
	Convey("all plat", t, func() {
		So(s.Notice("", 20), ShouldResemble, a)
		So(s.Notice("", 30), ShouldBeNil)
		So(s.Notice("", 10), ShouldBeNil)
	})
	Convey("ios plat", t, func() {
		So(s.Notice("ios", 25), ShouldBeNil)
		So(s.Notice("ios", 30), ShouldResemble, b)
		So(s.Notice("ios", 40), ShouldResemble, b)
	})
	Convey("android plat", t, func() {
		So(s.Notice("android", 25), ShouldResemble, c)
		So(s.Notice("android", 50), ShouldResemble, c)
		So(s.Notice("android", 60), ShouldBeNil)
	})
}
