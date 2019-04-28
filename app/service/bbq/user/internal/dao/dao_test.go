package dao

import (
	"flag"
	"go-common/app/service/bbq/user/internal/conf"
	"go-common/library/conf/paladin"
	"os"
	"testing"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "bbq.bbq.user")
		flag.Set("conf_token", "3523c3d66dbc91081454d8ef7af9a680")
		flag.Set("tree_id", "75578")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/test.toml")
		flag.Set("deploy.env", "uat")
	}
	flag.Parse()
	if err := paladin.Init(); err != nil {
		panic(err)
	}
	if err := paladin.Watch("test.toml", conf.Conf); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}
