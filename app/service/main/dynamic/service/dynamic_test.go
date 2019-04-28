package service

import (
	"context"
	"go-common/app/service/main/archive/api"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_RegionTotal(t *testing.T) {
	Convey("RegionTotal BigData", t, WithService(func(s *Service) {
		for _, v := range s.regionTotal {
			So(v, ShouldBeGreaterThan, 0)
		}
	}))
	Convey("RegionTotal Live", t, WithService(func(s *Service) {
		So(s.live, ShouldBeGreaterThan, 0)
	}))
}

func TestService_RegionArcs(t *testing.T) {
	var (
		c                   = context.TODO()
		rid           int32 = 168
		pn, ps, count int   = 1, 10, 0
		err           error
		rs            []*api.Arc
	)

	Convey("Region Arcives", t, WithService(func(s *Service) {
		rs, count, err = s.RegionArcs3(c, rid, pn, ps)
		So(err, ShouldBeNil)
		So(count, ShouldBeGreaterThan, 0)
		Printf("%+v", rs)
	}))

}

func TestService_RegionsArcs(t *testing.T) {
	var (
		c     = context.TODO()
		rids  = []int32{1, 3, 4, 5, 13, 36, 129, 119, 23, 11, 155, 160, 165, 168}
		count = 15
		err   error
		rs    map[int32][]*api.Arc
	)
	Convey("Regions Archives", t, WithService(func(s *Service) {
		rs, err = s.RegionsArcs3(c, rids, count)
		So(err, ShouldBeNil)
		Printf("%+v", rs)
	}))

}

func TestService_RegionTagArcs(t *testing.T) {
	var (
		c                   = context.TODO()
		tid           int64 = 123456
		rid           int32 = 168
		pn, ps, count int   = 1, 5, 0
		err           error
		rs            []*api.Arc
	)
	Convey("Region Tag Archives", t, WithService(func(s *Service) {
		rs, count, err = s.RegionTagArcs3(c, rid, tid, pn, ps)
		So(err, ShouldBeNil)
		Printf("%+v", rs)
		Printf("%d", count)
	}))
}
