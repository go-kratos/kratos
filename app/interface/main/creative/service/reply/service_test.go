package reply

import (
	"context"
	"flag"
	"go-common/app/interface/main/creative/conf"
	seamdl "go-common/app/interface/main/creative/model/search"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/creative/service"

	"github.com/davecgh/go-spew/spew"
	. "github.com/smartystreets/goconvey/convey"
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

func Test_Replies(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		oid = int64(0)
		res *seamdl.Replies
		err error
		p   *seamdl.ReplyParam
	)
	Convey("Replies", t, WithService(func(s *Service) {
		p = &seamdl.ReplyParam{
			Ak:   "ak",
			Ck:   "ck",
			OMID: mid,
			OID:  oid,
			Pn:   1,
			Ps:   10,
		}
		res, err = s.Replies(c, p)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		spew.Dump(res)
	}))
}

func Test_Archives(t *testing.T) {
	var (
		oids = []int64{1, 2, 3, 4}
		ip   = "127.0.0.1"
		c    = context.TODO()
	)
	Convey("Archives", t, WithService(func(s *Service) {
		res, err := s.arc.Archives(c, oids, ip)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}

func Test_Articles(t *testing.T) {
	var (
		oids = []int64{1, 2, 3, 4}
		ip   = "127.0.0.1"
		c    = context.TODO()
	)
	Convey("Articles", t, WithService(func(s *Service) {
		res, err := s.art.ArticleMetas(c, oids, ip)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}

func Test_Audio(t *testing.T) {
	var (
		oids  = []int64{47276}
		level = 0
		ip    = "127.0.0.1"
		c     = context.TODO()
	)
	Convey("Audio", t, WithService(func(s *Service) {
		res, err := s.mus.Audio(c, oids, level, ip)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		spew.Dump(res)
	}))
}
