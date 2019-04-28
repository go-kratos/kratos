package dao

import (
	"flag"
	"go-common/app/admin/main/manager/conf"
	"os"
	"path/filepath"
	"testing"
)

var d *Dao

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.manager.manager-admin")
		flag.Set("conf_token", "b661550b1a2cce530787e00ed1625581")
		flag.Set("tree_id", "11038")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		dir, _ := filepath.Abs("../cmd/manager-admin-test.toml")
		if err := flag.Set("conf", dir); err != nil {
			panic(err)
		}
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}
