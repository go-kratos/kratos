package service

import (
	"context"
	"flag"
	"path/filepath"
	"time"

	"go-common/app/service/main/coin/conf"
)

var (
	s   *Service
	ctx = context.TODO()
)

func init() {
	dir, _ := filepath.Abs("../cmd/coin-service-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}
