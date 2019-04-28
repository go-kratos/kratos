package dao

import (
	"flag"
	"os"
	"testing"

	"go-common/app/job/main/member-cache/conf"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.account.member-cache-job")
		flag.Set("conf_token", "31291cb6cd6b252ff27bf58d7fec72a6")
		flag.Set("tree_id", "56587")
		flag.Set("conf_version", "docker-1")
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
