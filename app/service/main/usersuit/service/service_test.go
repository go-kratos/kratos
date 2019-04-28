package service

import (
	"flag"
	"path/filepath"

	"go-common/app/service/main/usersuit/conf"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/convey-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
}
