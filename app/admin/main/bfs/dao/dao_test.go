package dao

import (
	"flag"
	"os"
	"testing"

	"go-common/app/admin/main/bfs/conf"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("conf_appid", "main.common-arch.bfs-admin")
		flag.Set("conf_token", "117c7ca29fdb846f73d96122491bb75d")
		flag.Set("conf_env", "uat")
		flag.Set("tree_id", "73143")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}
