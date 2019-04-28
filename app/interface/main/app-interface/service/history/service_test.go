package history

import (
	"flag"
	"path/filepath"
	"time"

	"go-common/app/interface/main/app-interface/conf"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/app-interface-test.toml")
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
