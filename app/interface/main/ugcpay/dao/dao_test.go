package dao

import (
	"context"
	"flag"
	"go-common/app/interface/main/ugcpay/conf"
	"os"
	"testing"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.account.ugcpay-interface")
		flag.Set("conf_token", "e6b959f71772ec8207eb80b8a86137cc")
		flag.Set("tree_id", "63917")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	d.Ping(context.Background())
	i := m.Run()
	d.Close()
	os.Exit(i)
}
