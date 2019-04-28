package income

import (
	"flag"
	"path/filepath"
	"time"

	"go-common/app/admin/main/growup/conf"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../../cmd/growup-admin.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		// Reset(func() { CleanCache() })
		f(s)
	}
}
