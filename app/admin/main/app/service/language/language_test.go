package language

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/admin/main/app/conf"
	"go-common/app/admin/main/app/model/language"

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

func TestLanguages(t *testing.T) {
	Convey("get Languages", t, WithService(func(s *Service) {
		res, err := s.Languages(context.TODO())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestLangByID(t *testing.T) {
	Convey("select language by id", t, WithService(func(s *Service) {
		res, err := s.LangByID(context.TODO(), 1)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestInsert(t *testing.T) {
	Convey("get Insert", t, WithService(func(s *Service) {
		a := &language.Param{
			Name:   "2233",
			Remark: "繁体中文",
		}
		err := s.Insert(context.TODO(), a)
		So(err, ShouldBeNil)
	}))
}

func TestUpdate(t *testing.T) {
	Convey("update language", t, WithService(func(s *Service) {
		a := &language.Param{
			ID:     2,
			Name:   "白居易",
			Remark: "简体中文",
		}
		err := s.Update(context.TODO(), a)
		So(err, ShouldBeNil)
	}))
}
