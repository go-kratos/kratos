package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_HelpList(t *testing.T) {
	Convey("test help HelpList", t, WithService(func(s *Service) {
		pTypeID := "39e79104ca5d433f9d9cdf9df4bf28a0"
		res, err := s.HelpList(context.Background(), pTypeID)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
		Println(len(res))
	}))
}

func TestService_HelpDetail(t *testing.T) {
	Convey("test help HelpDetail", t, WithService(func(s *Service) {
		fID := "39e79104ca5d433f9d9cdf9df4bf28a0"
		qTypeID := "39e79104ca5d433f9d9cdf9df4bf28a0"
		resD, resL, total, err := s.HelpDetail(context.Background(), fID, qTypeID, 0, 1, 15)
		So(err, ShouldBeNil)
		So(total, ShouldBeGreaterThan, 0)
		So(len(resD), ShouldBeGreaterThan, 0)
		Println(len(resD))
		So(len(resL), ShouldBeGreaterThan, 0)
		Println(len(resL))
	}))
}

func TestService_HelpSearch(t *testing.T) {
	Convey("test help HelpSearch", t, WithService(func(s *Service) {
		pTypeID := "39e79104ca5d433f9d9cdf9df4bf28a0"
		keyWords := "aaa"
		res, total, err := s.HelpSearch(context.Background(), pTypeID, keyWords, 0, 1, 15)
		So(err, ShouldBeNil)
		So(total, ShouldBeGreaterThan, 0)
		So(len(res), ShouldBeGreaterThan, 0)
		Println(len(res))
	}))
}
