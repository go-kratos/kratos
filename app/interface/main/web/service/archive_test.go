package service

import (
	"context"
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_View(t *testing.T) {
	Convey("test archive view", t, WithService(func(s *Service) {
		var (
			mid int64 = 27515256
			aid int64 = 10110688
			cid int64 = 1
		)
		res, err := s.View(context.Background(), aid, cid, mid, "", "")
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		str, _ := json.Marshal(res)
		Printf("%s", str)
	}))
}

func TestService_ArchiveStat(t *testing.T) {
	Convey("test archive archiveStat", t, WithService(func(s *Service) {
		var aid int64 = 5464686
		res, err := s.ArchiveStat(context.Background(), aid)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}

func TestService_AddShare(t *testing.T) {
	Convey("test archive AddShare", t, WithService(func(s *Service) {
		var (
			mid int64 = 27515256
			aid int64 = 5464686
		)
		res, err := s.AddShare(context.Background(), aid, mid, "", "", "", "", "")
		So(err, ShouldBeNil)
		So(res, ShouldBeGreaterThan, 0)
	}))
}

func TestService_Description(t *testing.T) {
	Convey("test archive Description", t, WithService(func(s *Service) {
		var (
			aid  int64 = 5464686
			page int64 = 1
		)
		res, err := s.Description(context.Background(), aid, page)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
}

func TestService_ArcReport(t *testing.T) {
	Convey("test archive ArcReport", t, WithService(func(s *Service) {
		var (
			mid    int64 = 27515256
			aid    int64 = 5464686
			tp     int64
			reason string
			pics   string
		)
		err := s.ArcReport(context.Background(), mid, aid, tp, reason, pics)
		So(err, ShouldBeNil)
	}))
}

func TestService_AppealTags(t *testing.T) {
	Convey("test archive AppealTags", t, WithService(func(s *Service) {
		res, err := s.AppealTags(context.Background())
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}

func TestService_AuthorRecommend(t *testing.T) {
	Convey("test archive AuthorRecommend", t, WithService(func(s *Service) {
		var aid int64 = 5464686
		res, err := s.AuthorRecommend(context.Background(), aid)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}

func TestService_RelatedArcs(t *testing.T) {
	Convey("test archive RelatedArcs", t, WithService(func(s *Service) {
		var aid int64 = 5464686
		res, err := s.RelatedArcs(context.Background(), aid)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}

func TestService_Detail(t *testing.T) {
	Convey("test archive Detail", t, WithService(func(s *Service) {
		var aid int64 = 5464686
		res, err := s.Detail(context.Background(), aid, 0, "", "")
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}
