package dao

import (
	"flag"
	"path/filepath"

	"go-common/app/job/main/playlist/conf"
)

var d *Dao

func WithDao(f func(d *Dao)) func() {
	return func() {
		dir, _ := filepath.Abs("../cmd/playlist-job-test.toml")
		flag.Set("conf", dir)
		conf.Init()
		if d == nil {
			d = New(conf.Conf)
		}
		f(d)
	}
}
