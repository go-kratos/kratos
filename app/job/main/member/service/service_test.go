package service

import (
	"context"
	"flag"
	"testing"
	"time"

	"go-common/app/job/main/member/conf"
	"go-common/app/job/main/member/model"
	"go-common/library/log"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func initConf() {
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Xlog)
	defer log.Close()
}

func init() {
	flag.Set("conf", "../cmd/member-job-dev.toml")
	initConf()
	s = New(conf.Conf)
}
func TestAddexp(t *testing.T) {
	Convey("addexp", t, func() {
		err := s.addExp(context.Background(), &model.AddExp{Mid: 1})
		So(err, ShouldBeNil)
	})
}

func TestRecoverMoral(t *testing.T) {
	time.Sleep(time.Second * 2)
	Convey("recoverMoral", t, func() {
		err := s.recoverMoral(context.Background(), 2)
		So(err, ShouldBeNil)
	})
}
