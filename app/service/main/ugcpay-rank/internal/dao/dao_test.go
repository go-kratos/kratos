package dao

import (
	"flag"
	"os"
	"testing"

	"go-common/app/service/main/ugcpay-rank/internal/conf"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.account.ugcpay-rank-service")
		flag.Set("conf_token", "af0015c6e9072270c243e43df817cebf")
		flag.Set("tree_id", "81367")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New()
	os.Exit(m.Run())
}
