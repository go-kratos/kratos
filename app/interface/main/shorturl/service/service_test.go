package service

import (
	"flag"
	"path/filepath"
	"time"

	"go-common/app/interface/main/shorturl/conf"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/shorturl-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)

}

func WithService(f func(s *Service)) func() {
	return func() {
		f(s)
	}
}
