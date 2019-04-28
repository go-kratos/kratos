package service

import (
	"flag"
	"fmt"
	"path/filepath"

	"go-common/app/interface/main/favorite/conf"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/favorite-test.toml")
	flag.Set("conf", dir)
	if err := conf.Init(); err != nil {
		fmt.Printf("conf.Init() error(%v)", err)
	}
	s = New(conf.Conf)
}
