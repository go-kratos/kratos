package ad

import (
	"flag"
	"path/filepath"

	"go-common/app/interface/main/web-show/conf"
)

var d *Dao

func WithDao(f func(d *Dao)) func() {
	return func() {
		dir, _ := filepath.Abs("../cmd/web-show-test.toml")
		flag.Set("conf", dir)
		conf.Init()
		if d == nil {
			d = New(conf.Conf)
		}
		f(d)
	}
}
