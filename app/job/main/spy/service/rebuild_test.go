package service

import (
	"flag"
	"testing"

	"go-common/app/job/main/spy/conf"
	"go-common/library/log"
)

func TestServiceReBuild(t *testing.T) {
	flag.Parse()
	if err := conf.Init(); err != nil {
		t.Errorf("conf.Init() error(%v)", err)
		t.FailNow()
	}
	log.Init(conf.Conf.Xlog)
	defer log.Close()
	if s == nil {
		s = New(conf.Conf)
	}
	testReBuild(t, s)
}

func testReBuild(t *testing.T, s *Service) {
	s.reBuild()
}
