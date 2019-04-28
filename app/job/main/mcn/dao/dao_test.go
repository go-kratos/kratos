package dao

import (
	"flag"
	"os"
	"testing"

	"go-common/app/job/main/mcn/conf"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.mcn-job")
		flag.Set("conf_token", "6214e16a2d849e375d899cbe37df31d1")
		flag.Set("tree_id", "58888")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/mcn-job-test.toml")
	}
	if os.Getenv("UT_LOCAL_TEST") != "" {
		flag.Set("conf", "../cmd/mcn-job-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}
