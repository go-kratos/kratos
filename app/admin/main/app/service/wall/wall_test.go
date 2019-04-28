package wall

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/admin/main/app/conf"
	"go-common/app/admin/main/app/model/wall"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func WithService(f func(s *Service)) func() {
	return func() {
		f(s)
	}
}

func init() {
	dir, _ := filepath.Abs("../../cmd/app-admin-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestWallss(t *testing.T) {
	Convey("get Walls", t, WithService(func(s *Service) {
		res, err := s.Walls(context.TODO())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestWallByID(t *testing.T) {
	Convey("get WallByID", t, WithService(func(s *Service) {
		res, err := s.WallByID(context.TODO(), 1)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestInsert(t *testing.T) {
	Convey("insert wall", t, WithService(func(s *Service) {
		a := &wall.Param{
			Title:    "举杯邀明月",
			Name:     "lssssssss",
			Package:  "对影成三人",
			Logo:     "http://bilibili.com",
			Size:     "25",
			Download: "ssssss",
			Remark:   "sdfsf",
			Rank:     3,
			State:    0,
		}
		err := s.Insert(context.TODO(), a, time.Now())
		So(err, ShouldBeNil)
	}))
}

func TestUpdateWall(t *testing.T) {
	Convey("update wall", t, WithService(func(s *Service) {
		a := &wall.Param{
			ID:       30,
			Name:     "llssxxx",
			Title:    "举杯邀明月",
			Package:  "对影成三人",
			Logo:     "http://bilibili.com",
			Size:     "25",
			Download: "ssssss",
			Remark:   "sdfsf",
			Rank:     3,
			State:    0,
		}
		err := s.UpdateWall(context.TODO(), a, time.Now())
		So(err, ShouldBeNil)
	}))
}

func TestPublish(t *testing.T) {
	Convey("get Publish", t, WithService(func(s *Service) {
		err := s.Publish(context.TODO(), "1,2,3", time.Now())
		So(err, ShouldBeNil)
	}))
}
