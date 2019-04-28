package service

import (
	"context"
	"fmt"
	"go-common/app/admin/main/usersuit/model"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_Medal(t *testing.T) {
	Convey("return sth", t, func() {
		res, err := s.Medal(context.Background())
		for k, re := range res {
			fmt.Printf("%d %+v \n", k, re)
		}
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestService_MedalView(t *testing.T) {
	Convey("return sth", t, func() {
		res, err := s.MedalView(context.Background(), 3)
		fmt.Printf("%+v \n", res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestService_AddMedal(t *testing.T) {
	Convey("return sth", t, func() {
		np := &model.Medal{
			Name:        "出道偶像",
			Description: "播放五十万，准备发新单",
			Condition:   "所有自制视频总播放数>=50万",
			GID:         4,
			Level:       int8(2),
			Sort:        2,
			LevelRank:   "50万",
			IsOnline:    1,
			Image:       "/bfs/face/3f2d64f048b39fb6c26f3db39df47e6080ec0f9c.png",
			ImageSmall:  "/bfs/face/90c35d41d8a19b19474d6bac672394c17b444ce8.png",
		}
		err := s.AddMedal(context.Background(), np)
		fmt.Printf("%+v \n", err)
		So(err, ShouldBeNil)
	})
}

func TestService_UpMedal(t *testing.T) {
	Convey("return sth", t, func() {
		np := &model.Medal{
			Name:        "出道偶像111",
			Description: "播放五十万，准备发新单",
			Condition:   "所有自制视频总播放数>=50万",
			GID:         4,
			Level:       int8(2),
			Sort:        2,
			LevelRank:   "50万",
			IsOnline:    1,
			Image:       "/bfs/face/3f2d64f048b39fb6c26f3db39df47e6080ec0f9c.png",
			ImageSmall:  "/bfs/face/90c35d41d8a19b19474d6bac672394c17b444ce8.png",
		}
		err := s.UpMedal(context.Background(), 1, np)
		fmt.Printf("%+v \n", err)
		So(err, ShouldBeNil)
	})
}

func TestService_MedalGroup(t *testing.T) {
	Convey("return someting", t, func() {
		res, err := s.MedalGroup(context.Background())
		for k, re := range res {
			fmt.Printf("%d %+v \n", k, re)
		}
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestService_MedalGroupInfo(t *testing.T) {
	Convey("return someting", t, func() {
		res, err := s.MedalGroupInfo(context.Background())
		for k, re := range res {
			fmt.Printf("%d %+v \n", k, re)
		}
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestService_MedalGroupParent(t *testing.T) {
	Convey("return someting", t, func() {
		res, err := s.MedalGroupParent(context.Background())
		for k, re := range res {
			fmt.Printf("%d %+v \n", k, re)
		}
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestService_MedalGroupByID(t *testing.T) {
	Convey("return someting", t, func() {
		re, err := s.MedalGroupByGid(context.Background(), 2)
		So(err, ShouldBeNil)
		So(re, ShouldNotBeNil)
		fmt.Printf("%+v \n", re)
	})
}

func TestService_MedalGroupAdd(t *testing.T) {
	Convey("return someting", t, func() {
		pg := &model.MedalGroup{
			Name:     "test",
			PID:      1,
			Rank:     int8(1),
			IsOnline: int8(1),
		}
		err := s.MedalGroupAdd(context.Background(), pg)
		So(err, ShouldBeNil)
		fmt.Printf("%+v \n", err)
	})
}

func TestService_MedalGroupUp(t *testing.T) {
	Convey("return someting", t, func() {
		pg := &model.MedalGroup{
			Name:     "test222",
			PID:      2,
			Rank:     2,
			IsOnline: 0,
		}
		err := s.MedalGroupUp(context.Background(), 37, pg)
		So(err, ShouldBeNil)
		fmt.Printf("%+v \n", err)
	})
}

func TestService_MedalOwner(t *testing.T) {
	Convey("return someting", t, func() {
		res, err := s.MedalOwner(context.Background(), 111)
		for k, re := range res {
			fmt.Printf("%d %+v \n", k, re)
		}
		fmt.Printf("err:%+v \n", err)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestService_MedalOwnerAdd(t *testing.T) {
	Convey("return someting", t, func() {
		err := s.MedalOwnerAdd(context.Background(), 41, 1, "", "", 1)
		fmt.Printf("%+v \n", err)
		So(err, ShouldBeNil)
	})
}

func TestService_MedalAddList(t *testing.T) {
	Convey("return someting", t, func() {
		res, err := s.MedalOwnerAddList(context.Background(), 1)
		for k, re := range res {
			fmt.Printf("%d %+v \n", k, re)
		}
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestService_MedalOwnerUpActivated(t *testing.T) {
	Convey("return someting", t, func() {
		err := s.MedalOwnerUpActivated(context.Background(), 1, 1)
		fmt.Printf("%+v \n", err)
		So(err, ShouldBeNil)
	})
}

func TestService_MedalOwnerDel(t *testing.T) {
	Convey("return someting", t, func() {
		err := s.MedalOwnerDel(context.Background(), 1, 5, 0, "", "")
		fmt.Printf("%+v \n", err)
		So(err, ShouldBeNil)
	})
}
