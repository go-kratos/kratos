package service

import (
	"context"
	"flag"
	"path/filepath"
	"time"

	"go-common/app/admin/ep/melloi/conf"
)

var (
	s *Service
	c context.Context
)

func init() {
	dir, _ := filepath.Abs("../cmd/convey-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	c = context.TODO()
	time.Sleep(time.Second)
}
