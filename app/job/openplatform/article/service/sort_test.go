package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_UpdateSort(t *testing.T) {
	Convey("get data", t, WithoutProcService(func(s *Service) {
		err := s.UpdateSort(context.TODO())
		So(err, ShouldBeNil)
	}))
}
func Test_loadSortArts(t *testing.T) {
	Convey("get data", t, WithoutProcService(func(s *Service) {
		res, err := s.loadSortArts(context.TODO())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
}

func Test_TrimArts(t *testing.T) {
	arts := [][2]int64{
		{1, 11},
		{0, 10},
		{2, 12},
	}
	Convey("trim arts", t, func() {
		ns := trimArts(arts, 3)
		expt := [][2]int64{
			{2, 12},
			{1, 11},
		}
		So(ns, ShouldResemble, expt)
	})
}
