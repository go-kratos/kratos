package audit

import (
	"flag"
	"path/filepath"
	"time"

	"go-common/app/interface/main/tv/conf"
)

var (
	srv *Service
)

func init() {
	dir, _ := filepath.Abs("../../cmd/tv-interface.toml")
	flag.Set("conf", dir)
	conf.Init()
	srv = New(conf.Conf)
	time.Sleep(time.Second)
}
