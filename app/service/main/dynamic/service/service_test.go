package service

import (
	"flag"
	"path/filepath"
	"time"

	"go-common/app/service/main/dynamic/conf"
)

var (
	svr *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/dynamic-service-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	svr = New(conf.Conf)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		f(svr)
	}
}
