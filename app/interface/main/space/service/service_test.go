package service

import (
	"flag"
	"path/filepath"
	"time"

	"go-common/app/interface/main/space/conf"
	"go-common/library/log"
)

var svf *Service

func WithService(f func(s *Service)) func() {
	return func() {
		dir, _ := filepath.Abs("../cmd/space-test.toml")
		flag.Set("conf", dir)
		conf.Init()
		log.Init(conf.Conf.Log)
		svf = New(conf.Conf)
		time.Sleep(200 * time.Millisecond)
		f(svf)
	}
}
