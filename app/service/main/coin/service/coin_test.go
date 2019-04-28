package service

import (
	"context"
	"testing"

	pb "go-common/app/service/main/coin/api"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAdd(t *testing.T) {
	Convey("add coin", t, func() {
		arg := &pb.AddCoinReq{
			IP:       "",
			Mid:      4780461,
			Upmid:    4052089,
			MaxCoin:  2,
			Aid:      1,
			Business: "archive",
			Number:   1,
			Typeid:   1,
			PubTime:  0,
		}
		_, err := s.AddCoin(context.TODO(), arg)
		So(err, ShouldBeNil)
	})
}

func TestUpdateItemCoins(t *testing.T) {
	Convey("add coin", t, func() {
		err := s.UpdateItemCoins(context.TODO(), 5462972, 1, 22)
		So(err, ShouldBeNil)
	})
}

func TestList(t *testing.T) {
	Convey("list", t, func() {
		arg := &pb.ListReq{
			Mid:      88888929,
			Business: "archive",
			Ts:       0,
		}
		ls, err := s.List(context.TODO(), arg)
		So(err, ShouldBeNil)
		So(ls, ShouldNotBeEmpty)
	})
}
