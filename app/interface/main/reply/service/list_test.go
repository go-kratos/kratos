package service

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWithinFloor(t *testing.T) {
	var (
		ids = []int64{3, 5, 7, 9, 11, 13, 15}
	)
	// within asc
	Convey("withinFloor asc without page 1", t, WithService(func(s *Service) {
		r := withinFloor(ids, 1, 1, 5, true)
		So(r, ShouldBeFalse)
	}))
	Convey("withinFloor asc within page 1", t, WithService(func(s *Service) {
		r := withinFloor(ids, 4, 1, 5, true)
		So(r, ShouldBeTrue)
	}))
	Convey("withinFloor asc within page 1", t, WithService(func(s *Service) {
		r := withinFloor(ids, 16, 1, 20, true)
		So(r, ShouldBeTrue)
	}))
	Convey("withinFloor asc without page 1", t, WithService(func(s *Service) {
		r := withinFloor(ids, 16, 1, 5, true)
		So(r, ShouldBeFalse)
	}))
	// within desc
	Convey("withinFloor desc within page 1", t, WithService(func(s *Service) {
		r := withinFloor(ids, 16, 1, 5, false)
		So(r, ShouldBeTrue)
	}))
	Convey("withinFloor desc within page 2", t, WithService(func(s *Service) {
		r := withinFloor(ids, 16, 2, 5, false)
		So(r, ShouldBeFalse)
	}))
	Convey("withinFloor desc within page 1", t, WithService(func(s *Service) {
		r := withinFloor(ids, 2, 3, 20, false)
		So(r, ShouldBeTrue)
	}))
	Convey("withinFloor desc without page 1", t, WithService(func(s *Service) {
		r := withinFloor(ids, 1, 1, 5, false)
		So(r, ShouldBeFalse)
	}))
	Convey("withinFloor desc within page 1", t, WithService(func(s *Service) {
		r := withinFloor(ids, 1, 1, 20, true)
		So(r, ShouldBeTrue)
	}))
}
