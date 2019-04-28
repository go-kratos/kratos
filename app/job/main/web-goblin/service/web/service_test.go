package web

import (
	"flag"
	"path/filepath"
	"time"

	"go-common/app/job/main/web-goblin/conf"
)

var (
	srv *Service
)

func init() {
	dir, _ := filepath.Abs("../../cmd/web-goblin-job-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	srv = New(conf.Conf)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		f(srv)
	}
}
