package redis

import (
	"flag"
	"path/filepath"

	"go-common/app/job/main/videoup-report/conf"
)

var (
	d *Dao
)

func init() {
	dir, _ := filepath.Abs("../../cmd/videoup-report-job.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
}
