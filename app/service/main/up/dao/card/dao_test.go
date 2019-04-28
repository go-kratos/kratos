package card

import (
	"flag"
	"os"
	"testing"

	"go-common/app/service/main/up/conf"
	"go-common/app/service/main/up/dao"
	"go-common/app/service/main/up/dao/global"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", dao.AppID)
		flag.Set("conf_token", dao.UatToken)
		flag.Set("tree_id", dao.TreeID)
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/up-service.toml")
	}
	if os.Getenv("UT_LOCAL_TEST") != "" {
		flag.Set("conf", "../../cmd/up-service.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	global.Init(conf.Conf)
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}
