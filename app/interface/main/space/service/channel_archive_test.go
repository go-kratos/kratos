package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_SortChannelArc(t *testing.T) {
	Convey("test sort channel arc", t, WithService(func(s *Service) {
		var (
			mid    = int64(27515260)
			cid    = int64(37)
			aid    = int64(5462035)
			preAid = 1
		)
		err := s.SortChannelArc(context.Background(), mid, cid, aid, preAid)
		So(err, ShouldBeNil)
	}))
}

func TestService_AddChannelArc(t *testing.T) {
	Convey("add channel arc", t, WithService(func(s *Service) {
		mid := int64(27515260)
		cid := int64(34)
		aid := []int64{5462036}
		_, err := s.AddChannelArc(context.Background(), mid, cid, aid)
		So(err, ShouldBeNil)
	}))
}

func TestService_archives(t *testing.T) {
	Convey("archives", t, WithService(func(s *Service) {
		aids := []int64{4053640, 4053639, 4053638, 4053637, 4053635, 4053634, 4053633, 4053632, 4053631, 4053630, 4053629, 4053628, 4053627, 4053626, 4053625, 4053624, 4053623, 4053622, 4053617, 4053611, 4053608, 4053607, 4053605, 4053604, 4053603, 4053602, 4053600, 4053599, 4053598, 4053596, 4053595, 4053594, 4053593, 4053592, 4053591, 4053590, 4053589, 4053588, 4053587, 4053586, 4053585, 4053584, 4053583, 4053582, 4053581, 4053579, 4053577, 4053576, 4053575, 4053574}
		arcs, err := s.archives(context.Background(), aids)
		So(err, ShouldBeNil)
		for _, v := range arcs {
			Printf("%+v", v)
		}
	}))
}
