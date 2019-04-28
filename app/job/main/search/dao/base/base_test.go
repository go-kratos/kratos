package base

import (
	"flag"
	"fmt"
	"path/filepath"
	"testing"

	"go-common/app/job/main/search/conf"

	. "github.com/smartystreets/goconvey/convey"
)

func WithBase(f func(b *Base)) func() {
	return func() {
		dir, _ := filepath.Abs("../dao/cmd/goconvey.toml")
		flag.Set("conf", dir)
		flag.Parse()
		conf.Init()
		d := NewBase(conf.Conf)
		f(d)
	}
}

func Test_NewAppPool(t *testing.T) {
	Convey("newAppPool", t, WithBase(func(b *Base) {
		pool := b.newAppPool(b.D)
		fmt.Println(pool)
	}))
}
