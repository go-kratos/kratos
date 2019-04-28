package service

import (
	"context"
	"testing"

	"go-common/app/interface/main/space/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_Channel(t *testing.T) {
	Convey("test channel", t, WithService(func(s *Service) {
		mid := int64(27515260)
		cid := int64(37)
		res, err := s.Channel(context.Background(), mid, cid)
		So(err, ShouldBeNil)
		Printf("res %+v", res)
		So(res, ShouldHaveSameTypeAs, &model.Channel{})
	}))
}

func TestService_ChannelIndex(t *testing.T) {
	Convey("test channel index", t, WithService(func(s *Service) {
		mid := int64(27515260)
		isGuest := false
		res, err := s.ChannelIndex(context.Background(), mid, isGuest)
		So(err, ShouldBeNil)
		for _, v := range res {
			Printf("res %+v", v)
		}
		var sample []*model.ChannelDetail
		So(res, ShouldHaveSameTypeAs, sample)
	}))
}

func TestService_ChannelList(t *testing.T) {
	Convey("test channel list", t, WithService(func(s *Service) {
		mid := int64(27515260)
		isGuest := false
		res, err := s.ChannelList(context.Background(), mid, isGuest)
		So(err, ShouldBeNil)
		for _, v := range res {
			Printf("res %+v", v)
		}
		var sample []*model.Channel
		So(res, ShouldHaveSameTypeAs, sample)
	}))
}

func TestService_DelChannel(t *testing.T) {
	Convey("del channel", t, WithService(func(s *Service) {
		mid := int64(27515260)
		cid := int64(34)
		err := s.DelChannel(context.Background(), mid, cid)
		So(err, ShouldBeNil)
	}))
}
