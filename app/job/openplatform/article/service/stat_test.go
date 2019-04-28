package service

import (
	"context"
	"testing"

	artmdl "go-common/app/interface/openplatform/article/model"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Stat(t *testing.T) {

	var (
		err     error
		c       = context.TODO()
		statCnt = int64(5)
		stat    = &artmdl.StatMsg{
			Aid:      88,
			View:     &statCnt,
			Favorite: &statCnt,
			Like:     &statCnt,
			Dislike:  &statCnt,
			Reply:    &statCnt,
			Share:    &statCnt,
		}
	)

	Convey("updateCache", t, WithoutProcService(func(s *Service) {
		err = s.updateCache(c, stat, 0)
		So(err, ShouldBeNil)
	}))

	Convey("updateDB", t, WithoutProcService(func(s *Service) {
		err = s.updateDB(c, stat, 0)
		So(err, ShouldBeNil)
	}))

	Convey("select Stat", t, WithoutProcService(func(s *Service) {
		var stat *artmdl.StatMsg
		stat, err = s.dao.Stat(c, 1)
		So(err, ShouldBeNil)
		So(stat, ShouldNotBeEmpty)
	}))
}
