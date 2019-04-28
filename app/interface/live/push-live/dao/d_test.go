package dao

import (
	"flag"
	"go-common/app/interface/live/push-live/conf"
	"path/filepath"
)

var d *Dao

func initd() {
	dir, _ := filepath.Abs("../cmd/push-live-test.toml")
	flag.Set("conf", dir)
	flag.Set("conf_env", "uat")
	conf.Init()
	d = New(conf.Conf)
}
