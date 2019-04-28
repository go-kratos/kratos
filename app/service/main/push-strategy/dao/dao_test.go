package dao

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"go-common/app/service/main/push-strategy/conf"

	"gopkg.in/h2non/gock.v1"
)

var d *Dao

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.web-svr.push-strategy")
		flag.Set("conf_token", "a626fc404c86f14654f0a74d80fc1da3")
		flag.Set("tree_id", "26092")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		dir, _ := filepath.Abs("../cmd/push-strategy-test.toml")
		flag.Set("conf", dir)
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	d.httpClient.SetTransport(gock.DefaultTransport)
	m.Run()
	os.Exit(0)
}
