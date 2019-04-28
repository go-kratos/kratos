package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_Attentions(t *testing.T) {
	Convey("test account attentions", t, WithService(func(s *Service) {
		var mid int64 = 27515256
		res, err := s.Attentions(context.Background(), mid)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
		Println(len(res))
	}))
}

func TestService_Card(t *testing.T) {
	Convey("test account card", t, WithService(func(s *Service) {
		var (
			mid      int64 = 27515256
			loginID  int64 = 37515257
			topPhoto       = true
			article        = true
		)
		res, err := s.Card(context.Background(), mid, loginID, topPhoto, article)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}
