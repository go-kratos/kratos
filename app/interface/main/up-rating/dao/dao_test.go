package dao

import (
	"flag"
	"os"
	"testing"

	"go-common/app/interface/main/up-rating/conf"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.up-rating-interface")
		flag.Set("conf_appid", "main.archive.up-rating-interface")
		flag.Set("conf_token", "db506d78d859ee901ac000bac393484c")
		flag.Set("tree_id", "65467")
		flag.Set("conf_version", "0.0.1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/up-rating-interface.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}
