package service

import (
	"fmt"
	"testing"

	"go-common/app/admin/main/tv/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_SeasonValidation(t *testing.T) {
	Convey("Season Validation Test", t, WithService(func(s *Service) {
		var season model.TVEpSeason
		if err := s.DB.Model(&model.TVEpSeason{}).Where("`check`=?", 1).
			Where("valid=?", 1).Where("is_deleted=?", 0).First(&season).Error; err != nil {
			fmt.Printf("Error:(%v)", err)
			return
		}
		fmt.Printf("Target ID is: %d", season.ID)
		res, sModel := s.snValid(season.ID)
		So(res, ShouldBeTrue)
		So(sModel.ID == season.ID, ShouldBeTrue)
	}))
}

func TestService_Intervs(t *testing.T) {
	Convey("Get Intervention List", t, WithService(func(s *Service) {
		res, err := s.Intervs(&model.IntervListReq{
			Rank:     0,
			Category: 1,
		})
		So(err, ShouldBeNil)
		fmt.Println(res)
	}))
}

func TestService_RemoveInvalids(t *testing.T) {
	Convey("Remove Invalid Test", t, WithService(func(s *Service) {
		var (
			rank     model.Rank
			err      error
			invalids []*model.RankError
		)
		if err = s.DB.Where("is_deleted=?", 0).First(&rank).Error; err != nil {
			fmt.Println(err)
			return
		}
		invalids = append(invalids, &model.RankError{
			ID:       int(rank.ID),
			SeasonID: int(rank.ContID),
		})
		err = s.RemoveInvalids(invalids)
		So(err, ShouldBeNil)
		// recover
		err = s.DB.Model(rank).Where("id=?", rank.ID).Update(map[string]int{"is_deleted": 0}).Error
		So(err, ShouldBeNil)
	}))
}
