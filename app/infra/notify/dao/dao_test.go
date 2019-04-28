package dao

import (
	"flag"
	"go-common/app/infra/notify/conf"
	"os"
	"testing"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.common-arch.notify")
		flag.Set("conf_token", "9919040d8742a2124abbed8368626075")
		flag.Set("tree_id", "22513")
		flag.Set("conf_version", "v1.0.0")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}
