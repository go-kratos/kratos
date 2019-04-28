package feed

import (
	"flag"
	"path/filepath"
	"time"

	"go-common/app/job/main/app/conf"
)

var (
	d *Dao
)

func init() {
	dir, _ := filepath.Abs("../../cmd/app-job-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(5 * time.Second)
}
