package dao

import (
	"flag"
	"os"
	"testing"

	"go-common/app/service/main/secure/conf"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.account.secure-service")
		flag.Set("conf_appid", "main.account.secure-service")
		flag.Set("conf_token", "02bebea2b9e0b8d5b0433b2a59b8bc68")
		flag.Set("tree_id", "2864")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_env", "10")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/secure-service-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}
