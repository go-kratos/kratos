package service

import (
	"flag"
	"go-common/app/job/main/figure/conf"
	"go-common/library/log"
	"path/filepath"
	"sync"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	once sync.Once
	s    *Service
)

func initConf() {
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
}

func init() {
	var (
		err error
	)
	dir, _ := filepath.Abs("../cmd/figure-job-test.toml")
	flag.Set("conf", dir)
	if err = conf.Init(); err != nil {
		panic(err)
	}
	if s == nil {
		s = New(conf.Conf)
	}
	time.Sleep(time.Second)
}

func startService() {
	initConf()
	if s == nil {
		s = New(conf.Conf)
	}
	time.Sleep(time.Second * 2)
}

func CleanCache() {
}

func WithService(f func(s *Service)) func() {
	return func() {
		Reset(func() { CleanCache() })
		f(s)
	}
}
