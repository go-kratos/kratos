package dao

import (
	"flag"
	"os"
	"testing"

	"go-common/app/service/main/passport-game/conf"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "dev" {
		flag.Set("app_id", "main.account.passport-game-service")
		flag.Set("conf_token", "6461f64e47f084a2c9bbd505182dbe39")
		flag.Set("tree_id", "3791")
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
