package like

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_fmtStartEnd(t *testing.T) {
	Convey("test fmt start and end", t, WithService(func(s *Service) {
		pn := 1
		ps := 10
		cnt := 11
		typ := "random" //ctime random
		start, end, err := s.fmtStartEnd(pn, ps, cnt, typ)
		So(err, ShouldBeNil)
		Println(start, end)
	}))
}

func TestService_LikeInitialize(t *testing.T) {
	Convey("test LikeInitialize", t, WithService(func(s *Service) {
		lid := int64(13537)
		err := s.LikeInitialize(context.Background(), lid)
		time.Sleep(time.Second)
		So(err, ShouldBeNil)
	}))
}
func TestService_LikeMaxIDInitialize(t *testing.T) {
	Convey("test LikeInitialize", t, WithService(func(s *Service) {
		err := s.LikeMaxIDInitialize(context.Background())
		So(err, ShouldBeNil)
	}))
}

func TestService_LikeUp(t *testing.T) {
	Convey("test LikeUp", t, WithService(func(s *Service) {
		lid := int64(13538)
		err := s.LikeUp(context.Background(), lid)
		time.Sleep(time.Second)
		So(err, ShouldBeNil)
	}))
}

func TestService_AddLikeCtimeCache(t *testing.T) {
	Convey("test LikeUp", t, WithService(func(s *Service) {
		lid := int64(13540)
		err := s.AddLikeCtimeCache(context.Background(), lid)
		So(err, ShouldBeNil)
	}))
}

func TestService_DelLikeCtimeCache(t *testing.T) {
	Convey("test LikeUp", t, WithService(func(s *Service) {
		err := s.DelLikeCtimeCache(context.Background(), 13537, 10296, 5)
		So(err, ShouldBeNil)
	}))
}

func TestService_SetLikeContent(t *testing.T) {
	Convey("test LikeUp", t, WithService(func(s *Service) {
		err := s.SetLikeContent(context.Background(), 13511)
		So(err, ShouldBeNil)
	}))
}

func TestService_AddLikeActCache(t *testing.T) {
	Convey("test LikeUp", t, WithService(func(s *Service) {
		err := s.AddLikeActCache(context.Background(), 10296, 13528, 7)
		So(err, ShouldBeNil)
	}))
}
