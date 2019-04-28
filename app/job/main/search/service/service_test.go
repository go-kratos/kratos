package service

import (
	"flag"
	"path/filepath"

	"go-common/app/job/main/search/conf"
)

func WithService(f func(s *Service)) func() {
	return func() {
		dir, _ := filepath.Abs("../goconvey.toml")
		flag.Set("conf", dir)
		conf.Init()
		s := New(conf.Conf)
		// s.dao = dao.New(conf.Conf)
		f(s)
	}
}
