package service

import (
	"context"
	"fmt"
	"go-common/app/admin/main/activity/model"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_LikesByLid(t *testing.T) {
	Convey("like get items", t, WithService(func(s *Service) {
		res, err := s.LikesList(context.Background(), &model.LikesParam{Sid: 10256, PageSize: 10, Page: 1, Mid: 155551800, States: []int{0, 1}})
		So(err, ShouldBeEmpty)
		for _, v := range res.Likes {
			fmt.Printf("%+v", v)
		}
	}))
}

func TestService_Likes(t *testing.T) {
	Convey("like get items", t, WithService(func(s *Service) {
		res, err := s.Likes(context.Background(), 10256, []int64{1185, 1256})
		So(err, ShouldBeEmpty)
		for _, v := range res {
			fmt.Printf("%+v", v)
		}
	}))
}

func TestService_UpLike(t *testing.T) {
	Convey("like get items", t, WithService(func(s *Service) {
		res, err := s.AddLike(context.Background(), &model.AddLikes{DealType: "videoAdd", Sid: 10206, Wid: 10210488, Mid: 88895364, Device: 11, Plat: 12, State: 1, Type: 12})
		So(err, ShouldBeEmpty)
		fmt.Printf("%+v", res)
	}))
}

func TestService_UpLikeContents(t *testing.T) {
	Convey("like get items", t, WithService(func(s *Service) {
		res, err := s.UpLike(context.Background(), &model.UpLike{Lid: 13557, Type: 13, State: 1, Message: "ii", Reply: "nono", Image: "jj", Link: "like", Mid: 12345, Wid: 12367, StickTop: 1}, "ly")
		So(err, ShouldBeEmpty)
		fmt.Printf("%+v", res)
	}))
}

func TestService_BatchLikes(t *testing.T) {
	Convey("like get items", t, WithService(func(s *Service) {
		err := s.BatchLikes(context.Background(), &model.BatchLike{Sid: 10299, Type: 13, Mid: 12345, Wid: []int64{1, 4, 5, 6, 7}})
		So(err, ShouldBeEmpty)
	}))
}
