package manager

import (
	"flag"
	"go-common/app/service/main/resource/conf"
	"os"
	"testing"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "prod" {
		flag.Set("app_id", "main.app-svr.resource-service")
		flag.Set("conf_token", "y79sErNhxggjvULS0O8Czas9PaxHBF5o")
		flag.Set("tree_id", "3232")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")

	} else {
		flag.Set("conf", "../../cmd/resource-service-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}
