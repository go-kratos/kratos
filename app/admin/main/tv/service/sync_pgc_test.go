package service

import (
	"go-common/app/admin/main/tv/model"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_EpCheck(t *testing.T) {
	Convey("TestService_EpCheck test", t, WithService(func(s *Service) {
		var cont = model.Content{}
		// already passed to 7
		err := s.EpCheck(&model.Content{
			ID:    106396,
			State: 3,
		})
		So(err, ShouldBeNil)
		s.DB.Where("id=?", 106396).First(&cont)
		So(cont.State == 7, ShouldBeTrue)
		// reject to re_audit
		err = s.EpCheck(&model.Content{
			ID:    106396,
			State: 4,
		})
		So(err, ShouldBeNil)
		s.DB.Where("id=?", 106396).First(&cont)
		So(cont.State == 1, ShouldBeTrue)
	}))
}

func TestService_SeasonCheck(t *testing.T) {
	Convey("TestService_SeasonCheck test", t, WithService(func(s *Service) {
		var season = model.TVEpSeason{}
		// already passed to 7
		err := s.SeasonCheck(&model.TVEpSeason{
			ID:    296,
			Check: 1})
		So(err, ShouldBeNil)
		s.DB.Where("id=?", 296).First(&season)
		So(season.Check == 7, ShouldBeTrue)
		// reject to re_audit
		err = s.SeasonCheck(&model.TVEpSeason{
			ID:    296,
			Check: 0})
		So(err, ShouldBeNil)
		s.DB.Where("id=?", 296).First(&season)
		So(season.Check == 2, ShouldBeTrue)
	}))
}

func TestService_EpDel(t *testing.T) {
	Convey("TestService_EpDel test", t, WithService(func(s *Service) {
		var (
			cont = model.Content{}
			ep   = model.TVEpContent{}
			epid = 5822
		)
		// delete
		err := s.EpDel(int64(epid), 1)
		s.DB.Where("epid=?", epid).First(&cont)
		s.DB.Where("id=?", epid).First(&ep)
		So(err, ShouldBeNil)
		So(cont.IsDeleted == 1, ShouldBeTrue)
		So(ep.IsDeleted == 1, ShouldBeTrue)
		// recover
		err = s.EpDel(30185, 0)
		s.DB.Where("epid=?", 30185).First(&cont)
		s.DB.Where("id=?", 30185).First(&ep)
		So(err, ShouldBeNil)
		So(cont.IsDeleted == 1, ShouldBeTrue)
		So(ep.IsDeleted == 1, ShouldBeTrue)
	}))
}

func TestService_SeasonRemove(t *testing.T) {
	Convey("TestService_SeasonRemove test", t, WithService(func(s *Service) {
		var (
			season   = model.TVEpSeason{}
			seasonID = 20950
		)
		err := s.SeasonRemove(&model.TVEpSeason{
			ID:    int64(seasonID),
			Check: 1,
		})
		s.DB.Where("id=?", seasonID).First(&season)
		So(err, ShouldBeNil)
	}))
}
