package article

import (
	"context"
	"flag"
	"go-common/app/interface/main/creative/conf"
	article "go-common/app/interface/openplatform/article/model"
	"path/filepath"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/interface/main/creative/service"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../../cmd/creative.toml")
	flag.Set("conf", dir)
	conf.Init()
	rpcdaos := service.NewRPCDaos(conf.Conf)
	s = New(conf.Conf, rpcdaos)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		Reset(func() {})
		f(s)
	}
}

func Test_Categories(t *testing.T) {
	var (
		c   = context.TODO()
		err error
		res *article.Categories
	)
	Convey("Categories", t, WithService(func(s *Service) {
		res, err = s.Categories(c)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}

func Test_Article(t *testing.T) {
	var (
		c    = context.TODO()
		err  error
		size int
		loc  string
		url  = "https://i0.hdslb.com/bfs/article/8fce42a26ce128140d2ee7dec599a46cd1bccbb6.jpg"
	)
	Convey("Capture", t, WithService(func(s *Service) {
		loc, size, err = s.ArticleCapture(c, url)
		So(err, ShouldBeNil)
		spew.Dump(loc, size)
	}))
}
