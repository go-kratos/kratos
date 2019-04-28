package feed

import (
	"flag"
	"path/filepath"
	"time"

	"go-common/app/interface/main/app-feed/conf"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../../cmd/app-feed-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}
