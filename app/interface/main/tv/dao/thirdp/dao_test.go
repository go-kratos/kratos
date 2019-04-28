package thirdp

import (
	"context"
	"flag"
	"path/filepath"

	"go-common/app/interface/main/tv/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d   *Dao
	ctx = context.Background()
)

func init() {
	dir, _ := filepath.Abs("../../cmd/tv-interface.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
}

func WithDao(f func(d *Dao)) func() {
	return func() {
		Reset(func() {})
		f(d)
	}
}
