package service

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_TreatIntervs(t *testing.T) {
	Convey("Treat Interventions, result's length between 0 and 5", t, WithService(func(s *Service) {
		res, err := s.HomeRecom(ctx)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
		So(len(res), ShouldBeLessThanOrEqualTo, s.ZonesInfo[_homepageID].Top)
	}))
}

func TestService_FollowData(t *testing.T) {
	Convey("Follow Data API is ok", t, WithService(func(s *Service) {
		res := s.FollowData(ctx, "00f04e06e1caf7f8cbe17580d3fa3e62")
		So(res, ShouldNotBeNil)
	}))
}

func TestService_HomeList(t *testing.T) {
	Convey("Follow Data API is ok", t, WithService(func(s *Service) {
		res, latest := s.HomeList()
		So(len(res), ShouldNotEqual, 0)
		So(len(latest), ShouldNotEqual, 0)
	}))
}
