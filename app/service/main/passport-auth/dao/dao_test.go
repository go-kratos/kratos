package dao

import (
	"flag"
	"go-common/app/service/main/passport-auth/conf"
	"os"
	"testing"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.passport.passport-auth-service")
		flag.Set("conf_token", "3f34a05a9bf3bc0add9c135c1f83ff2f")
		flag.Set("tree_id", "35706")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "/Users/bilibili/gopath/src/go-common/app/service/main/passport-auth/cmd/passport-auth-service.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}
