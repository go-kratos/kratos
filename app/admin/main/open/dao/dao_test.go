package dao

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"go-common/app/admin/main/open/conf"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "")
		flag.Set("conf_token", "")
		flag.Set("tree_id", "")
		flag.Set("conf_version", "server-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	}
	flag.Parse()
	dir, _ := filepath.Abs("../cmd/open-admin-test.toml")
	if err := flag.Set("conf", dir); err != nil {
		panic(err)
	}

	if err := conf.Init(); err != nil {
		panic(err)
	}

	d = New(conf.Conf)
	m.Run()
	os.Exit(1)
}
