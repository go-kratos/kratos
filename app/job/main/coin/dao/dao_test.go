package dao

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/job/main/coin/conf"
)

var (
	d   *Dao
	ctx = context.TODO()
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.account.coin-job")
		flag.Set("conf_appid", "coin-job")
		flag.Set("conf_token", "RiSyxxaGuihmdbKCq9Pba1JwHRTmQjMI")
		flag.Set("tree_id", "2130")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/coin-job-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}
