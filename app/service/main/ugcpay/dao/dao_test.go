package dao

import (
	"context"
	"flag"
	"go-common/app/service/main/ugcpay/conf"
	"os"
	"testing"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.account.ugcpay-service")
		flag.Set("conf_token", "a3e98fbfb0f63520a8cdc699365460c2")
		flag.Set("tree_id", "62221")
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
	conf.Conf.CacheTTL.OrderTTL = 1
	conf.Conf.CacheTTL.AssetTTL = 1
	conf.Conf.CacheTTL.AssetRelationTTL = 1
	conf.Conf.CacheTTL.AggrIncomeUserTTL = 1
	conf.Conf.CacheTTL.AggrIncomeUserMonthlyTTL = 1
	if err := d.Ping(context.Background()); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
