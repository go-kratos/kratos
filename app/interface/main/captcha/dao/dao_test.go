package dao

import (
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/interface/main/captcha/conf"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	flag.Parse()
	dir, _ := filepath.Abs("../cmd/captcha-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
}
