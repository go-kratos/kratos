package service

import (
	"context"
	"flag"
	"path/filepath"
	"time"

	"go-common/app/interface/main/spread/conf"
)

var (
	s *Service
	c = context.TODO()
)

func init() {
	dir, _ := filepath.Abs("../cmd/convey-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}
