package playurl

import (
	"context"
	"flag"
	"go-common/app/job/main/tv/conf"
	"path/filepath"
)

var (
	ctx = context.TODO()
	d   *Dao
)

func WithDao(f func(d *Dao)) func() {
	return func() {
		dir, _ := filepath.Abs("../../cmd/tv-job-test.toml")
		flag.Set("conf", dir)
		conf.Init()
		if d == nil {
			d = New(conf.Conf)
		}
		f(d)
	}
}
