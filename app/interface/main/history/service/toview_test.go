package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TestService_AddHistory
func TestService_Toview(t *testing.T) {
	var (
		c          = context.TODO()
		mid  int64 = 27515316
		aids       = []int64{27515316, 27515316}
		ip         = ""
		pn         = 1
		ps         = 10000
	)
	Convey("toview ", t, WithService(func(s *Service) {
		Convey("toview AddMultiToView", func() {
			err := s.AddMultiToView(c, mid, aids, ip)
			So(err, ShouldBeNil)
		})
		Convey("toview RemainingToView", func() {
			_, err := s.RemainingToView(c, mid, "")
			So(err, ShouldBeNil)
		})
		Convey("toview ClearToView", func() {
			err := s.ClearToView(c, mid)
			So(err, ShouldBeNil)
		})
		Convey("toview DelToView", func() {
			err := s.DelToView(c, mid, aids, true, "")
			So(err, ShouldBeNil)
		})
		Convey("toview cache del", func() {
			_, _, err := s.ToView(c, mid, pn, ps, ip)
			So(err, ShouldBeNil)
		})
		Convey("toview Manager", func() {
			_, err := s.ManagerToView(c, mid, ip)
			So(err, ShouldBeNil)
		})
	}))
}
