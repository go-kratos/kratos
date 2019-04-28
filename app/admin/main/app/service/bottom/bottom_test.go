package bottom

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/admin/main/app/conf"
	"go-common/app/admin/main/app/model/bottom"

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

func TestBottoms(t *testing.T) {
	Convey("get bottoms", t, WithService(func(s *Service) {
		res, err := s.Bottoms(context.TODO())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestBottomByID(t *testing.T) {
	Convey("select bottom by id ", t, WithService(func(s *Service) {
		res, err := s.BottomByID(context.TODO(), 2)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestInsert(t *testing.T) {
	Convey("insert bottom", t, WithService(func(s *Service) {
		a := &bottom.Param{
			Name:   "wuwu",
			Logo:   "http://i0.hdslb.com/oseA456.jpg",
			Rank:   44,
			Action: 1,
			Param:  "545",
			State:  1,
		}
		err := s.Insert(context.TODO(), a, time.Now())
		So(err, ShouldBeNil)
	}))
}

func TestUpdate(t *testing.T) {
	Convey("update bottom", t, WithService(func(s *Service) {
		a := &bottom.Param{
			ID:     17,
			Name:   "wuwu again",
			Logo:   "http://i0.hdslb.com/oseA456.jpg",
			Rank:   44,
			Action: 1,
			Param:  "545",
			State:  1,
		}
		err := s.Update(context.TODO(), a, time.Now())
		So(err, ShouldBeNil)
	}))
}

func TestPublish(t *testing.T) {
	Convey("update state", t, WithService(func(s *Service) {
		err := s.Publish(context.TODO(), "2,3,5,6", time.Now())
		So(err, ShouldBeNil)
	}))
}

func TestDelete(t *testing.T) {
	Convey("delete", t, WithService(func(s *Service) {
		err := s.Delete(context.TODO(), 11)
		So(err, ShouldBeNil)
	}))
}
