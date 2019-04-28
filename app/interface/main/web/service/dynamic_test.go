package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_DynamicRegion(t *testing.T) {
	Convey("test dynamic DynamicRegion", t, WithService(func(s *Service) {
		var rid int32 = 23
		res, err := s.DynamicRegion(context.Background(), rid, 1, 10)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}

func TestService_DynamicRegionTag(t *testing.T) {
	Convey("test dynamic DynamicRegionTag", t, WithService(func(s *Service) {
		var (
			tagID int64 = 10101
			rid   int32 = 23
		)
		res, err := s.DynamicRegionTag(context.Background(), tagID, rid, 1, 10)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}

func TestService_DynamicRegionTotal(t *testing.T) {
	Convey("test dynamic DynamicRegionTotal", t, WithService(func(s *Service) {
		res, err := s.DynamicRegionTotal(context.Background())
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
		Println(len(res))
	}))
}

func TestService_DynamicRegions(t *testing.T) {
	Convey("test dynamic DynamicRegions", t, WithService(func(s *Service) {
		res, err := s.DynamicRegions(context.Background())
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
		Println(len(res))
	}))
}
