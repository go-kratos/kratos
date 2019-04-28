package report

import (
	"flag"
	"path/filepath"
	"time"

	"go-common/app/job/main/tv/conf"
)

var (
	srv *Service
)

func init() {
	dir, _ := filepath.Abs("../../cmd/tv-job-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	srv = New(conf.Conf)
	time.Sleep(time.Second)
}
